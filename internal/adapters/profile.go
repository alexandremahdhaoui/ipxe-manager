package adapters

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/google/uuid"
	apierrrors "k8s.io/apimachinery/pkg/api/errors"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// Useful errors

	ErrProfileNotFound = errors.New("profile cannot be found")

	// FindBySelectors

	ErrFindBySelectors                     = errors.New("error finding profile by selectors")
	errListAssignment                      = errors.New("listing assignment")
	errFallbackToDefaultAssignment         = errors.New("fallback to default assignment")
	errSelectingAssignment                 = errors.New("selecting assignment")
	fmtCannotSelectAssignmentWithSelectors = "cannot select profile assignment with selectors: uuid=%q & buildarch=%q"

	// FindByID

	ErrProfileFindByID = errors.New("error finding profile by id")

	// Conversions

	errConvertingProfile                     = errors.New("converting profile")
	errToProfileID                           = errors.New("converting to profile uuid")
	errExposedAdditionalContentCannotBeFound = errors.New("profile cannot be found in exposed additional content")
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Profile interface {
	FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Profile, error)
	FindByID(ctx context.Context, id uuid.UUID) (types.Profile, error)

	// NB: CRUD operations are done via the reconciler-runtime client.Client; only FindBySelectorsAndRender is
	// required.
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewProfile(c client.Client, namespace string) Profile {
	return &profile{
		client:    c,
		namespace: namespace,
	}
}

// --------------------------------------------- CONCRETE IMPLEMENTATION -------------------------------------------- //

type profile struct {
	client    client.Client
	namespace string
}

// --------------------------------------------- FindBySelectors ---------------------------------------------------- //

func (p *profile) FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Profile, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := p.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return types.Profile{}, errors.Join(err, errListAssignment, ErrFindBySelectors)
	}

	// Handle
	if list == nil || len(list.Items) == 0 {
		selectErr := fmt.Errorf(fmtCannotSelectAssignmentWithSelectors, selectors.UUID.String(), selectors.Buildarch)
		errChain := errors.Join(selectErr, errFallbackToDefaultAssignment, errSelectingAssignment, ErrFindBySelectors)
		// Get the list of default matching the buildarch
		if err := p.client.List(ctx, list, toDefaultListOptions(selectors.Buildarch)...); err != nil {
			return types.Profile{}, errors.Join(err, errListAssignment, errChain)
		}

		if list == nil || len(list.Items) == 0 {
			return types.Profile{}, errors.Join(ErrProfileNotFound, errChain)
		}
	}

	profileName := list.Items[0].Spec.ProfileName
	obj := new(v1alpha1.Profile)
	key := k8stypes.NamespacedName{Name: profileName, Namespace: p.namespace}

	if err := p.client.Get(ctx, key, obj); err != nil {
		return types.Profile{}, errors.Join(err, ErrFindBySelectors)
	}

	out, err := fromV1alpha1.toProfile(obj)
	if err != nil {
		return types.Profile{}, errors.Join(err, ErrFindBySelectors)
	}

	return out, nil
}

func toListOptions(selectors types.IpxeSelectors) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: selectors.Buildarch,
		},
		client.HasLabels{
			v1alpha1.UUIDLabelSelector(selectors.UUID),
		},
	}
}

func toDefaultListOptions(buildarch string) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: buildarch,
		},
		client.HasLabels{
			v1alpha1.DefaultAssignmentLabel,
		},
	}
}

// --------------------------------------------- FindByID ----------------------------------------------------------- //

//TODO(question): change this func, we are trying to find profile by the machine id, not by the profile's ID or name.
//TODO(question): what if there are multiple profile with a label targeting a machine ID? (use webhook validations to disallow creating 2 profile for the same machine ID?)

func (p *profile) FindByID(ctx context.Context, id uuid.UUID) (types.Profile, error) {
	obj := new(v1alpha1.Profile)
	key := k8stypes.NamespacedName{Name: id.String(), Namespace: p.namespace}
	if err := p.client.Get(ctx, key, obj); apierrrors.IsNotFound(err) {
		return types.Profile{}, errors.Join(err, ErrProfileNotFound, ErrProfileFindByID)
	} else if err != nil {
		return types.Profile{}, errors.Join(err, ErrProfileFindByID)
	}

	out, err := fromV1alpha1.toProfile(obj)
	if err != nil {
		return types.Profile{}, errors.Join(err, ErrProfileFindByID)
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

		transformers := fromV1alpha1.toTransformerConfig(ac.PostTransformations)

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
			content = types.NewObjectRefContent(ac.Name, fromV1alpha1.toObjectRef(ac.ObjectRef),
				transformers...)
		case ac.Webhook != nil:
			content = types.NewWebhookContent(ac.Name, fromV1alpha1.toWebhookConfig(ac.Webhook),
				transformers...)
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

func (ipxev1a1) toObjectRef(objectRef *v1alpha1.ObjectRef) types.ObjectRef {
	return types.ObjectRef{
		Ref:  objectRef.TypedObjectReference,
		Path: objectRef.Path,
	}
}

func (ipxev1a1) toTransformerConfig(input []v1alpha1.Transformer) []types.TransformerConfig {
	out := make([]types.TransformerConfig, 0)

	for _, t := range input {
		var cfg types.TransformerConfig

		switch {
		case t.ButaneToIgnition:
			cfg.Kind = types.ButaneTransformerKind
		case t.Webhook != nil:
			cfg.Kind = types.WebhookTransformerKind
			cfg.Webhook = types.Ptr(fromV1alpha1.toWebhookConfig(t.Webhook))
		}

		out = append(out, cfg)
	}

	return out
}

func (ipxev1a1) toWebhookConfig(input *v1alpha1.WebhookConfig) types.WebhookConfig {
	out := types.WebhookConfig{
		URL: input.URL,
	}

	if input.MTLSObjectRef != nil {
		out.MTLSObjectRef = &types.MTLSObjectRef{
			Ref:            input.MTLSObjectRef.TypedObjectReference,
			ClientKeyPath:  input.MTLSObjectRef.ClientKeyPath,
			ClientCertPath: input.MTLSObjectRef.ClientCertPath,
			CaBundlePath:   input.MTLSObjectRef.CaBundlePath,
		}
	}

	if input.BasicAuthObjectRef != nil {
		out.BasicAuthObjectRef = &types.BasicAuthObjectRef{
			Ref:          input.BasicAuthObjectRef.TypedObjectReference,
			UsernamePath: input.BasicAuthObjectRef.UsernamePath,
			PasswordPath: input.BasicAuthObjectRef.PasswordPath,
		}
	}

	return out
}
