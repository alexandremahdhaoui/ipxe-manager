package adapter

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxe-api/pkg/v1alpha1"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Profile interface {
	FindBySelectors(ctx context.Context, selectors IpxeSelectors) (ProfileType, error)

	// NB: CRUD operations are done via the controller-runtime client.Client; only FindBySelectorsAndRender is
	// required.
}

func NewProfile(c client.Client, namespace string) Profile {
	return profile{
		client:    c,
		namespace: namespace,
	}
}

// --------------------------------------------- CONCRETE IMPLEMENTATION -------------------------------------------- //

type profile struct {
	client    client.Client
	namespace string
}

func (p profile) FindBySelectors(ctx context.Context, selectors IpxeSelectors) (ProfileType, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := p.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return ProfileType{}, err //TODO: wrap err.
	}

	if list == nil || len(list.Items) == 0 {
		// Get the list of default matching the buildarch
		if err := p.client.List(ctx, list, toDefaultListOptions(selectors.Buildarch)...); err != nil {
			return ProfileType{}, err //TODO: wrap err.
		}

		if list == nil || len(list.Items) == 0 {
			return ProfileType{}, errors.New("TODO") //TODO: err
		}
	}

	profileName := list.Items[0].Spec.ProfileName
	prof := new(v1alpha1.Profile)

	if err := p.client.Get(ctx, types.NamespacedName{Namespace: p.namespace, Name: profileName}, prof); err != nil {
		return ProfileType{}, err //TODO: wrap err
	}

	profileType, err := convertProfileType(prof)
	if err != nil {
		return ProfileType{}, err //TODO: wrap err
	}

	return profileType, nil
}

func toListOptions(selectors IpxeSelectors) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: selectors.Buildarch,
		},
		client.HasLabels{
			v1alpha1.UUIDLabelSelector(selectors.Uuid),
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

// ------------------------------------------------ TO PROFILE MODEL ------------------------------------------------ //

func convertProfileType(input *v1alpha1.Profile) (ProfileType, error) {
	out := ProfileType{}
	out.IPXETemplate = input.Spec.IPXETemplate
	out.AdditionalContent = make([]Content, 0)

	for _, cntSpec := range input.Spec.AdditionalContent {
		id, err := convertProfileID(cntSpec.Name, input.Status)
		if err != nil {
			return ProfileType{}, err //TODO: wrap err
		}

		transformers := convertTransformers(cntSpec.PostTransformations)

		var content Content
		switch {
		case cntSpec.Inline != nil:
			content = NewInlineContent(id, cntSpec.Name, *cntSpec.Inline, transformers...)
		case cntSpec.ObjectRef != nil:
			content = NewObjectRefContent(id, cntSpec.Name, convertObjectRef(cntSpec.ObjectRef), transformers...)
		case cntSpec.Webhook != nil:
			content = NewWebhookContent(id, cntSpec.Name, convertWebhookConfig(cntSpec.Webhook), transformers...)
		}

		out.AdditionalContent = append(out.AdditionalContent, content)
	}

	return out, nil
}

func convertProfileID(name string, status v1alpha1.ProfileStatus) (uuid.UUID, error) {
	id, ok := status.ExposedAdditionalContent[name]
	if !ok {
		return uuid.Nil, errors.New("TODO") //TODO: err
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err //TODO: wrap err
	}

	return uid, nil
}

func convertObjectRef(objectRef *v1alpha1.ObjectRef) ObjectRef {
	return ObjectRef{
		Ref:  objectRef.TypedObjectReference,
		Path: objectRef.Path,
	}
}

func convertTransformers(input []v1alpha1.Transformer) []TransformerConfig {
	out := make([]TransformerConfig, 0)

	for _, t := range input {
		var cfg TransformerConfig

		switch {
		case t.ButaneToIgnition:
			cfg.Kind = ButaneTransformerKind
		case t.Webhook != nil:
			cfg.Kind = WebhookTransformerKind
			cfg.Webhook = ptr(convertWebhookConfig(t.Webhook))
		}

		out = append(out, cfg)
	}

	return out
}

func convertWebhookConfig(input *v1alpha1.WebhookConfig) WebhookConfig {
	out := WebhookConfig{
		URL: input.URL,
	}

	if input.MTLSObjectRef != nil {
		out.MTLSObjectRef = &MTLSObjectRef{
			Ref:            input.MTLSObjectRef.TypedObjectReference,
			ClientKeyPath:  input.MTLSObjectRef.ClientKeyPath,
			ClientCertPath: input.MTLSObjectRef.ClientCertPath,
			CaBundlePath:   input.MTLSObjectRef.CaBundlePath,
		}
	}

	if input.BasicAuthObjectRef != nil {
		out.BasicAuthObjectRef = &BasicAuthObjectRef{
			Ref:          input.BasicAuthObjectRef.TypedObjectReference,
			UsernamePath: input.BasicAuthObjectRef.UsernamePath,
			PasswordPath: input.BasicAuthObjectRef.PasswordPath,
		}
	}

	return out
}
