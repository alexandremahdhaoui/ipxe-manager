package webhook

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.CustomValidator = &Profile{}
var _ webhook.CustomDefaulter = &Profile{}

func NewProfile() *Profile {
	return &Profile{}
}

type Profile struct{}

func (p *Profile) Default(ctx context.Context, obj runtime.Object) error {
	_, ok := obj.(*v1alpha1.Profile)
	if !ok {
		return NewUnsupportedResource(obj) //TODO: wrap err
	}

	if err := p.validateProfileStatic(ctx, obj); err != nil {
		return err //TODO: wrap err
	}

	//TODO implement me
	panic("implement me")
}

func (p *Profile) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	if err := p.validateProfileStatic(ctx, obj); err != nil {
		return nil, err //TODO: wrap err
	}

	if err := p.validateProfileDynamic(ctx, obj); err != nil {
		return nil, err //TODO: wrap err
	}

	return nil, nil
}

func (p *Profile) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	if err := p.validateProfileStatic(ctx, newObj); err != nil {
		return nil, err //TODO: wrap err
	}

	if err := p.validateProfileDynamic(ctx, newObj); err != nil {
		return nil, err //TODO: wrap err
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
			return err //TODO: wrap err
		}
	}

	return nil
}

func (p *Profile) validateProfileDynamic(ctx context.Context, obj runtime.Object) error {
	for _, f := range []validatingFunc{
		//TODO
	} {
		if err := f(ctx, obj); err != nil {
			return err //TODO: wrap err
		}
	}

	return nil
}

func validateIPXETemplate(ctx context.Context, obj runtime.Object) error {}

func validateAdditionalContent(ctx context.Context, obj runtime.Object) error {
	profile := obj.(*v1alpha1.Profile)

	for _, content := range profile.Spec.AdditionalContent {
		//TODO implement me
		panic("implement me")
	}

	return nil
}

func validateTransformer(ctx context.Context, obj runtime.Object) error        {}
func validateObjectRef(ctx context.Context, obj runtime.Object) error          {}
func validateWebhookConfig(ctx context.Context, obj runtime.Object) error      {}
func validateBasicAuthObjectRef(ctx context.Context, obj runtime.Object) error {}
func validateMTLSObjectRef(ctx context.Context, obj runtime.Object) error      {}
func validateResourceRef(ctx context.Context, obj runtime.Object) error        {}
func validateWebhookJSONPath(ctx context.Context, obj runtime.Object) error    {}
