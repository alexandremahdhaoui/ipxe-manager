package adapter

import (
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//TODO:
// When recursively templating:
//   - First check if we have DAG.
//   - If any cycles is spotted we should not allow such an operation to be performed.
// Additionally, we should ensure that no v1alpha1.Profile Custom Resource can be created if there are cycles:
//   - We need to create a Validating webhook that checks for cycles.
// Finally, we also need a form of runtime invalidation mechanism for dynamic non-DAGs:
//   - Indeed, to ensure that the content of "v1alpha1.ArbitraryResources"- or "v1alpha1.WebhookContent"-
//     v1alpha1.AdditionalContent does not contain cycles.
//   - We can create a DAG upon requesting all those information however, we should use BFS in order to avoid infinite
//     cycles. A max depth when running BFS might be a good solution.

// Template
//
// There is a body containing references to other configs that can themselves contain references to other configs.
//
// Templating happens in 3 phases:
//   1. find references
//   2. resolve references
//   3. render the template
//
// After rendering the template, we can recursively search for any references to a template.

type ContentResolverType int

const (
	InlineResolverType ContentResolverType = iota
	ObjectRefResolverType
	WebhookResolverType
)

func NewInlineContent(s string, postTransformers ...TransformerConfig) Content {
	return Content{
		ResolverType:     InlineResolverType,
		PostTransformers: postTransformers,
		Inline:           s,
	}
}

func NewObjectRefContent(objectRef ObjectRef, postTransformers ...TransformerConfig) Content {
	return Content{
		ResolverType:     ObjectRefResolverType,
		PostTransformers: postTransformers,
		ObjectRef:        &objectRef,
	}
}

func NewWebhookContent(cfg WebhookConfig, postTransformers ...TransformerConfig) Content {
	return Content{
		ResolverType:     WebhookResolverType,
		PostTransformers: postTransformers,
		WebhookConfig:    &cfg,
	}
}

type (
	Content struct {
		ResolverType     ContentResolverType
		PostTransformers []TransformerConfig

		Inline        string
		ObjectRef     *ObjectRef
		WebhookConfig *WebhookConfig
	}

	ObjectRef struct {
		corev1.TypedObjectReference `json:",inline"`

		// Path to the desired content in the resource. E.g. `.data."private.key"`
		Path string `json:"path"`
	}

	WebhookConfig struct {
		URL           string                  `json:"url"`
		MtlsObjectRef *corev1.ObjectReference `json:"mtlsSecretRef,omitempty"`
	}

	MtlsObjectRef struct {
		corev1.TypedObjectReference `json:",inline"`

		ClientKeyPath  string  `json:"clientKeyPath"`
		ClientCertPath string  `json:"clientCertPath"`
		CaBundlePath   *string `json:"caBundlePath,omitempty"`
	}
)

type Resolver interface {
	Resolve(cfg Content) []byte
}

type inlineResolver struct{}

func (r *inlineResolver) Resolve(cfg Content) []byte {
	return []byte(cfg.Inline)
}

func NewInlineResolver() Resolver {
	return &inlineResolver{}
}

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
