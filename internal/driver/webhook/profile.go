package webhook

import (
	"context"
	"errors"
	"regexp"

	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/jsonpath"

	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	_ webhook.CustomDefaulter = &Profile{}
	_ webhook.CustomValidator = &Profile{}

	// Regexes

	contentNameRegex = regexp.MustCompile("")
)

func NewProfile() *Profile {
	return &Profile{}
}

type Profile struct{}

func (p *Profile) Default(ctx context.Context, obj runtime.Object) error {
	profile, ok := obj.(*v1alpha1.Profile)
	if !ok {
		return NewUnsupportedResource(obj) // TODO: wrap err
	}

	if err := p.validateProfileStatic(ctx, obj); err != nil {
		return err // TODO: wrap err
	}

	// 1. get config UUIDs
	reverseIDMap := make(map[string]string)
	for k, value := range profile.Labels {
		if v1alpha1.IsUUIDLabelSelector(k) {
			reverseIDMap[value] = k // "content name" -> "label holding uuid"
		}
	}

	// 2. Remove all "internal" labels.
	for k := range profile.Labels {
		if !v1alpha1.IsInternalLabel(k) {
			delete(profile.Labels, k)
		}
	}

	// 3. Set labels preserving old UUIDs. (this is a bit overengineered, but may prevent a few race conditions).
	for name, content := range profile.Spec.AdditionalContent {
		if content.Exposed {
			if id, ok := reverseIDMap[name]; ok {
				profile.Labels[id] = name
			} else {
				profile.Labels[uuid.New().String()] = name
			}
		}
	}

	return nil
}

func (p *Profile) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	if err := p.validateProfileStatic(ctx, obj); err != nil {
		return nil, err // TODO: wrap err
	}

	if err := p.validateProfileDynamic(ctx, obj); err != nil {
		return nil, err // TODO: wrap err
	}

	return nil, nil
}

func (p *Profile) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	if err := p.validateProfileStatic(ctx, newObj); err != nil {
		return nil, err // TODO: wrap err
	}

	if err := p.validateProfileDynamic(ctx, newObj); err != nil {
		return nil, err // TODO: wrap err
	}

	return nil, nil
}

func (p *Profile) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (p *Profile) validateProfileStatic(ctx context.Context, obj runtime.Object) error {
	for _, f := range []validatingFunc{
		validateIPXETemplate,
		validateAdditionalContent,
	} {
		if err := f(ctx, obj); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func (p *Profile) validateProfileDynamic(ctx context.Context, obj runtime.Object) error {
	for _, f := range []validatingFunc{
		// TODO
	} {
		if err := f(ctx, obj); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func validateIPXETemplate(_ context.Context, _ runtime.Object) error {
	return nil
}

func validateAdditionalContent(ctx context.Context, obj runtime.Object) error {
	profile := obj.(*v1alpha1.Profile)

	for name, content := range profile.Spec.AdditionalContent {
		if !contentNameRegex.MatchString(name) { // TODO: create the regex
			return errors.New("TODO") // TODO: err + wrap err
		}

		for _, transformer := range content.PostTransformations {
			if err := validateTransformer(transformer); err != nil {
				return err // TODO: wrap err
			}
		}

		var i uint
		for _, ptr := range []any{
			content.Inline,
			content.ObjectRef,
			content.Webhook,
		} {
			if ptr != nil {
				i += 1
			}
		}

		switch {
		case i == 0:
			return errors.New("TODO") // TODO: err + wrap err
		case i > 1:
			return errors.New("TODO") // TODO: err + wrap err
		case content.Inline != nil:
			return nil
		case content.ObjectRef != nil:
			if err := validateObjectRef(content.ObjectRef); err != nil {
				return err // TODO: wrap err
			}
		case content.Webhook != nil:
			if err := validateWebhookConfig(content.Webhook); err != nil {
				return err // TODO: wrap err
			}
		}
	}

	panic("open an issue on github") // this branch does not exist, open an issue if you manage to pass the above s/c.
}

func validateObjectRef(ref *v1alpha1.ObjectRef) error {
	if err := validateResourceRef(ref.ResourceRef); err != nil {
		return err // TODO: wrap err
	}

	if err := validateJSONPath(ref.JSONPath); err != nil {
		return err // TODO: wrap err
	}

	return nil
}

func validateWebhookConfig(cfg *v1alpha1.WebhookConfig) error {
	if cfg.BasicAuthObjectRef != nil {
		if err := validateBasicAuthObjectRef(cfg.BasicAuthObjectRef); err != nil {
			return err // TODO: wrap err
		}
	}

	if cfg.MTLSObjectRef != nil {
		if err := validateMTLSObjectRef(cfg.MTLSObjectRef); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func validateTransformer(transformer v1alpha1.Transformer) error {
	if transformer.Webhook != nil {
		if transformer.ButaneToIgnition == true {
			return errors.New("a transformer must either enable butaneToIgnition or specify a webhook") // TODO: wrap err
		}

		if err := validateWebhookConfig(transformer.Webhook); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func validateBasicAuthObjectRef(ref *v1alpha1.BasicAuthObjectRef) error {
	if err := validateResourceRef(ref.ResourceRef); err != nil {
		return err // TODO: wrap err
	}

	for _, s := range []string{
		ref.UsernameJSONPath,
		ref.PasswordJSONPath,
	} {
		if err := validateJSONPath(s); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func validateMTLSObjectRef(ref *v1alpha1.MTLSObjectRef) error {
	if err := validateResourceRef(ref.ResourceRef); err != nil {
		return err // TODO: wrap err
	}

	for _, s := range []string{
		ref.CaBundleJSONPath,
		ref.ClientCertJSONPath,
		ref.ClientKeyJSONPath,
	} {
		if err := validateJSONPath(s); err != nil {
			return err // TODO: wrap err
		}
	}

	return nil
}

func validateResourceRef(ref v1alpha1.ResourceRef) error {
	if ref.Name == "" || len(ref.Name) > 63 {
		return errors.New("invalid name")
	}

	return nil
}

func validateJSONPath(s string) error {
	if _, err := jsonpath.Parse("", s); err != nil {
		return err // TODO: wrap err
	}

	return nil
}
