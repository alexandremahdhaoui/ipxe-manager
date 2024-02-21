package adapter

import (
	"github.com/google/uuid"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ----------------------------------------------------- TYPES ------------------------------------------------------ //

type ContentResolverKind int

const (
	InlineResolverKind ContentResolverKind = iota
	ObjectRefResolverKind
	WebhookResolverKind
)

// --------------------------------------------------- INTERFACE ---------------------------------------------------- //

type Resolver interface {
	Resolve(cfg Content) []byte
}

// -------------------------------------------------- CONSTRUCTORS -------------------------------------------------- //

func NewInlineContent(
	id uuid.UUID,
	name, inline string,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		ID:               id,
		Name:             name,
		ResolverKind:     InlineResolverKind,
		PostTransformers: postTransformers,
		Inline:           inline,
	}
}

func NewObjectRefContent(
	id uuid.UUID,
	name string,
	objectRef ObjectRef,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		ID:               id,
		Name:             name,
		ResolverKind:     ObjectRefResolverKind,
		PostTransformers: postTransformers,
		ObjectRef:        &objectRef,
	}
}

func NewWebhookContent(
	id uuid.UUID,
	name string,
	cfg WebhookConfig,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		ID:               id,
		Name:             name,
		ResolverKind:     WebhookResolverKind,
		PostTransformers: postTransformers,
		WebhookConfig:    &cfg,
	}
}

// ------------------------------------------------- INLINE RESOLVER ------------------------------------------------ //

type inlineResolver struct{}

func (r *inlineResolver) Resolve(cfg Content) []byte {
	return []byte(cfg.Inline)
}

func NewInlineResolver() Resolver {
	return &inlineResolver{}
}

// ---------------------------------------------- OBJECT REF RESOLVER ----------------------------------------------- //

type objectRefResolver struct {
	client client.Client
}

func (r *objectRefResolver) Resolve(cfg Content) []byte {
	//TODO implement me
	panic("implement me")
}

func NewObjectRefResolver(c client.Client) Resolver {
	return &objectRefResolver{client: c}
}

// ------------------------------------------------ WEBHOOK RESOLVER ------------------------------------------------ //

type webhookResolver struct {
	client http.Client
}

func (r *webhookResolver) Resolve(cfg Content) []byte {
	//TODO implement me
	panic("implement me")
}

func NewWebhookResolver(c http.Client) Resolver {
	return &webhookResolver{client: c}
}
