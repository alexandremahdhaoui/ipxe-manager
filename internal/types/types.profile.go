package types

import (
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
)

// ---------------------------------------------------- PROFILE ----------------------------------------------------- //

type Profile struct {
	IPXETemplate      string
	AdditionalContent []Content
}

// ---------------------------------------------------- CONTENT ----------------------------------------------------- //

type Content struct {
	Name string

	PostTransformers []TransformerConfig
	ResolverKind     ResolverKind

	Inline          string
	ObjectRef       *ObjectRef
	WebhookConfig   *WebhookConfig
	ExposedConfigID uuid.UUID
}

type ObjectRef struct {
	Ref corev1.TypedObjectReference

	// Path to the desired content in the resource. E.g. `.data."private.key"`
	Path string
}

type WebhookConfig struct {
	URL string

	MTLSObjectRef      *MTLSObjectRef
	BasicAuthObjectRef *BasicAuthObjectRef
}

type BasicAuthObjectRef struct {
	Ref corev1.TypedObjectReference

	UsernamePath string
	PasswordPath string
}

type MTLSObjectRef struct {
	Ref corev1.TypedObjectReference

	ClientKeyPath  string
	ClientCertPath string
	CaBundlePath   *string
}

// --------------------------------------------------- RESOLVER ----------------------------------------------------- //

type ResolverKind int

const (
	InlineResolverKind ResolverKind = iota
	ObjectRefResolverKind
	WebhookResolverKind
)

// -------------------------------------------------- TRANSFORMER --------------------------------------------------- //

type TransformerKind int

const (
	ButaneTransformerKind TransformerKind = iota
	WebhookTransformerKind
)

type TransformerConfig struct {
	Kind TransformerKind

	Webhook *WebhookConfig
}

// ---------------------------------------------- CONTENT CONSTRUCTORS ---------------------------------------------- //

func NewInlineContent(
	name, inline string,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		Name:             name,
		ResolverKind:     InlineResolverKind,
		PostTransformers: postTransformers,
		Inline:           inline,
	}
}

func NewObjectRefContent(
	name string,
	objectRef ObjectRef,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		Name:             name,
		ResolverKind:     ObjectRefResolverKind,
		PostTransformers: postTransformers,
		ObjectRef:        &objectRef,
	}
}

func NewWebhookContent(
	name string,
	cfg WebhookConfig,
	postTransformers ...TransformerConfig,
) Content {
	return Content{
		Name:             name,
		ResolverKind:     WebhookResolverKind,
		PostTransformers: postTransformers,
		WebhookConfig:    &cfg,
	}
}

func NewExposedContent(id uuid.UUID, name string) Content {
	return Content{
		Name:            name,
		ExposedConfigID: id,
	}
}
