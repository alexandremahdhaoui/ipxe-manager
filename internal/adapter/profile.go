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
	ErrProfileNotFound = errors.New("profile cannot be found")

	errProfileGet = errors.New("error getting profile")

	// Conversions

	errConvertingProfile                     = errors.New("converting profile")
	errToProfileID                           = errors.New("converting to profile uuid")
	errExposedAdditionalContentCannotBeFound = errors.New("profile cannot be found in exposed additional content")
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Profile interface {
	Get(ctx context.Context, name string) (types.Profile, error)
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

// --------------------------------------------------- CONVERSION --------------------------------------------------- //

var fromV1alpha1 ipxev1a1

type ipxev1a1 struct{}

func (ipxev1a1) toProfile(input *v1alpha1.Profile) (types.Profile, error) {
	out := types.Profile{}
	out.IPXETemplate = input.Spec.IPXETemplate
	out.AdditionalContent = make([]types.Content, 0)

	for _, ac := range input.Spec.AdditionalContent {

		transformers, err := fromV1alpha1.toTransformerConfig(ac.PostTransformations)
		if err != nil {
			return types.Profile{}, err // TODO: wrap err
		}

		var content types.Content
		switch {
		case ac.Exposed:
			id, err := fromV1alpha1.toProfileID(ac.Name, input.Status)
			if err != nil {
				return types.Profile{}, errors.Join(err, errConvertingProfile)
			}

			content = types.NewExposedContent(id, ac.Name)
		case ac.Inline != nil:
			content = types.NewInlineContent(ac.Name, *ac.Inline, transformers...)
		case ac.ObjectRef != nil:
			ref, err := fromV1alpha1.toObjectRef(ac.ObjectRef)
			if err != nil {
				return types.Profile{}, err // TODO: wrap err
			}

			content = types.NewObjectRefContent(ac.Name, ref, transformers...)
		case ac.Webhook != nil:
			config, err := fromV1alpha1.toWebhookConfig(ac.Webhook)
			if err != nil {
				return types.Profile{}, err // TODO: wrap err
			}

			content = types.NewWebhookContent(ac.Name, config, transformers...)
		}

		out.AdditionalContent = append(out.AdditionalContent, content)
	}

	return out, nil
}

func (ipxev1a1) toProfileID(name string, status v1alpha1.ProfileStatus) (uuid.UUID, error) {
	id, ok := status.ExposedAdditionalContent[name]
	if !ok {
		return uuid.Nil, errors.Join(errExposedAdditionalContentCannotBeFound, errToProfileID)
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.Join(err, errToProfileID)
	}

	return uid, nil
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
		ClientKeyJSONPath:  ckjp,
		ClientCertJSONPath: ccjp,
		CaBundleJSONPath:   cbjp,
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
