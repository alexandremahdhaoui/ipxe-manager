package types

import (
	"github.com/google/uuid"
	"k8s.io/client-go/util/jsonpath"
)

// ---------------------------------------------------- PROFILE ----------------------------------------------------- //

type Profile struct {
	IPXETemplate string

	AdditionalContent        map[string]*Content
	AdditionalExposedContent map[uuid.UUID]*Content
}

// ---------------------------------------------------- CONTENT ----------------------------------------------------- //

type Content struct {
	PostTransformers []TransformerConfig
	ResolverKind     ResolverKind

	Inline        string
	ObjectRef     *ObjectRef
	WebhookConfig *WebhookConfig
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
