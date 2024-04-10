package adapter

import (
	"context"
	"errors"

	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/google/uuid"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
	errProfileGet      = errors.New("error getting profile")

	errProfileListByConfigID = errors.New("listing profile by config id")

	// Conversions

	errConvertingProfile                     = errors.New("converting profile")
	errToProfileID                           = errors.New("converting to profile uuid")
	errExposedAdditionalContentCannotBeFound = errors.New("profile cannot be found in exposed additional content")
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Profile interface {
	Get(ctx context.Context, name string) (types.Profile, error)
	ListByConfigID(ctx context.Context, configID uuid.UUID) ([]types.Profile, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewProfile(c client.Client, namespace string) Profile {
	return &v1a1Profile{
		client:    c,
		namespace: namespace,
	}
}

// --------------------------------------------- CONCRETE IMPLEMENTATION -------------------------------------------- //

type v1a1Profile struct {
	client    client.Client
	namespace string
}

// --------------------------------------------- Get ----------------------------------------------------------- //

func (p *v1a1Profile) Get(ctx context.Context, name string) (types.Profile, error) {
	obj := new(v1alpha1.Profile)

	if err := p.client.Get(ctx, k8stypes.NamespacedName{
		Name:      name,
		Namespace: p.namespace,
	}, obj); apierrors.IsNotFound(err) {
		return types.Profile{}, errors.Join(err, ErrProfileNotFound, errProfileGet)
	} else if err != nil {
		return types.Profile{}, errors.Join(err, errProfileGet)
	}

	out, err := fromV1alpha1.toProfile(obj)
	if err != nil {
		return types.Profile{}, errors.Join(err, errProfileGet)
	}

	return out, nil
}

// --------------------------------------------- ListByConfigID ------------------------------------------------------ //

// ListByConfigID retrieve at most one Profile by a config ID. The nature of UUIDs and the defaulting webhook driver
// ensures the list contains at most 1 ID.
func (p *v1a1Profile) ListByConfigID(ctx context.Context, configID uuid.UUID) ([]types.Profile, error) {
	// list profiles
	obj := new(v1alpha1.ProfileList)
	if err := p.client.List(ctx, obj,
		uuidLabelSelector(configID),
	); apierrors.IsNotFound(err) || len(obj.Items) == 0 {
		return nil, errors.Join(err, ErrProfileNotFound, errProfileListByConfigID)
	} else if err != nil {
		return nil, errors.Join(err, errProfileListByConfigID)
	}

	out := make([]types.Profile, 0, len(obj.Items))
	for i := range obj.Items {
		profile, err := fromV1alpha1.toProfile(&obj.Items[i])
		if err != nil {
			return nil, errors.Join(err, errProfileListByConfigID)
		}

		out = append(out, profile)
	}

	return out, nil
}

// --------------------------------------------------- CONVERSION --------------------------------------------------- //

var fromV1alpha1 ipxev1a1

type ipxev1a1 struct{}

func (ipxev1a1) toProfile(input *v1alpha1.Profile) (types.Profile, error) {
	idNameMap, rev, err := v1alpha1.UUIDLabelSelectors(input.Labels)
	if err != nil {
		return types.Profile{}, err //TODO: wrap err
	}

	out := types.Profile{
		IPXETemplate:       input.Spec.IPXETemplate,
		AdditionalContent:  make(map[string]types.Content),
		ContentIDToNameMap: idNameMap,
	}

	for name, c := range input.Spec.AdditionalContent {
		content := types.Content{}

		// exposed
		if c.Exposed {
			content.Exposed = true

			id, ok := rev[name]
			if !ok {
				return types.Profile{}, errors.New("additional content is exposed but doesn't have a UUID") //TODO: err + wrap err
			}

			content.ExposedUUID = id
		}

		// post transformers
		transformers, err := fromV1alpha1.toTransformerConfig(c.PostTransformations)
		if err != nil {
			return types.Profile{}, err // TODO: wrap err
		}

		content.PostTransformers = transformers

		switch {
		case c.Inline != nil:
			content.ResolverKind = types.InlineResolverKind
			content.Inline = *c.Inline
		case c.ObjectRef != nil:
			ref, err := fromV1alpha1.toObjectRef(c.ObjectRef)
			if err != nil {
				return types.Profile{}, err // TODO: wrap err
			}

			content.ResolverKind = types.ObjectRefResolverKind
			content.ObjectRef = &ref
		case c.Webhook != nil:
			cfg, err := fromV1alpha1.toWebhookConfig(c.Webhook)
			if err != nil {
				return types.Profile{}, err // TODO: wrap err
			}

			content.ResolverKind = types.WebhookResolverKind
			content.WebhookConfig = &cfg
		}

		out.AdditionalContent[name] = content
	}

	return out, nil
}

func (ipxev1a1) toObjectRef(objectRef *v1alpha1.ObjectRef) (types.ObjectRef, error) {
	jp, err := toJSONPath(objectRef.JSONPath)
	if err != nil {
		return types.ObjectRef{}, err // TODO: wrap err
	}

	return types.ObjectRef{
		Group:     objectRef.Group,
		Version:   objectRef.Version,
		Resource:  objectRef.Resource,
		Namespace: objectRef.Namespace,
		Name:      objectRef.Name,
		JSONPath:  jp,
	}, nil
}

func (ipxev1a1) toTransformerConfig(input []v1alpha1.Transformer) ([]types.TransformerConfig, error) {
	out := make([]types.TransformerConfig, 0)

	for _, t := range input {
		var cfg types.TransformerConfig

		switch {
		case t.ButaneToIgnition:
			cfg.Kind = types.ButaneTransformerKind
		case t.Webhook != nil:
			typesCfg, err := fromV1alpha1.toWebhookConfig(t.Webhook)
			if err != nil {
				return nil, err // TODO: wrap err
			}

			cfg.Kind = types.WebhookTransformerKind
			cfg.Webhook = types.Ptr(typesCfg)
		}

		out = append(out, cfg)
	}

	return out, nil
}

func (ipxev1a1) toWebhookConfig(input *v1alpha1.WebhookConfig) (types.WebhookConfig, error) {
	out := types.WebhookConfig{}
	out.URL = input.URL

	if input.MTLSObjectRef != nil {
		ref, err := fromV1alpha1.toMTLSObjectRef(input.MTLSObjectRef)
		if err != nil {
			return types.WebhookConfig{}, err // TODO: wrap err
		}

		out.MTLSObjectRef = ref
	}

	if input.BasicAuthObjectRef != nil {
		ref, err := fromV1alpha1.toBasicAuthObjectRef(input.BasicAuthObjectRef)
		if err != nil {
			return types.WebhookConfig{}, err // TODO: wrap err
		}

		out.BasicAuthObjectRef = ref
	}

	return out, nil
}

func (ipxev1a1) toMTLSObjectRef(ref *v1alpha1.MTLSObjectRef) (*types.MTLSObjectRef, error) {
	ckjp, err := toJSONPath(ref.ClientKeyJSONPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	ccjp, err := toJSONPath(ref.ClientCertJSONPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	cbjp, err := toJSONPath(ref.CaBundleJSONPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	return &types.MTLSObjectRef{
		ObjectRef: types.ObjectRef{
			Group:     ref.Group,
			Version:   ref.Version,
			Resource:  ref.Resource,
			Namespace: ref.Namespace,
			Name:      ref.Name,
		},
		ClientKeyJSONPath:     ckjp,
		ClientCertJSONPath:    ccjp,
		CaBundleJSONPath:      cbjp,
		TLSInsecureSkipVerify: ref.TLSInsecureSkipVerify,
	}, nil
}

func (ipxev1a1) toBasicAuthObjectRef(ref *v1alpha1.BasicAuthObjectRef) (*types.BasicAuthObjectRef, error) {
	ujp, err := toJSONPath(ref.UsernameJSONPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	pjp, err := toJSONPath(ref.PasswordJSONPath)
	if err != nil {
		return nil, err // TODO: wrap err
	}

	return &types.BasicAuthObjectRef{
		ObjectRef: types.ObjectRef{
			Group:     ref.Group,
			Version:   ref.Version,
			Resource:  ref.Resource,
			Namespace: ref.Namespace,
			Name:      ref.Name,
		},
		UsernameJSONPath: ujp,
		PasswordJSONPath: pjp,
	}, nil
}

func toJSONPath(s string) (*jsonpath.JSONPath, error) {
	jp := jsonpath.New("")
	if err := jp.Parse(s); err != nil {
		return nil, err // TODO: wrap err
	}

	return jp, nil
}
