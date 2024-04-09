package types

import (
	"github.com/google/uuid"
	"k8s.io/client-go/util/jsonpath"
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
	Group     string
	Version   string
	Resource  string
	Namespace string
	Name      string

	// JSONPath is optional for types that extends this struct.
	JSONPath *jsonpath.JSONPath
}

type WebhookConfig struct {
	URL string

	MTLSObjectRef      *MTLSObjectRef
	BasicAuthObjectRef *BasicAuthObjectRef
}

type BasicAuthObjectRef struct {
	ObjectRef

	UsernameJSONPath *jsonpath.JSONPath
	PasswordJSONPath *jsonpath.JSONPath
}

type MTLSObjectRef struct {
	ObjectRef

	ClientKeyJSONPath  *jsonpath.JSONPath
	ClientCertJSONPath *jsonpath.JSONPath
	CaBundleJSONPath   *jsonpath.JSONPath

	TLSInsecureSkipVerify bool
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
