package adapter

import (
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
)

// ---------------------------------------------------- PROFILE ----------------------------------------------------- //

type ProfileType struct {
	IPXETemplate      string
	AdditionalContent []Content
}

// ---------------------------------------------------- CONTENT ----------------------------------------------------- //

type Content struct {
	Name string
	ID   uuid.UUID

	PostTransformers []TransformerConfig
	ResolverKind     ContentResolverKind

	Inline        string
	ObjectRef     *ObjectRef
	WebhookConfig *WebhookConfig
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
