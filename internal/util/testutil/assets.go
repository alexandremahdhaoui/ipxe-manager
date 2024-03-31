package testutil

import (
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"k8s.io/client-go/util/jsonpath"
)

const (
	inlineName    = "test-inline"
	inlineContent = "test inline content"

	objectRefName = "test-object-ref"
	webhookName   = "test-webhook"

	ipxeTemplate = "abc123"
)

func NewV1alpha1Profile() v1alpha1.Profile {
	return v1alpha1.Profile{
		Spec: v1alpha1.ProfileSpec{
			IPXETemplate: ipxeTemplate,
			AdditionalContent: []v1alpha1.AdditionalContent{
				NewV1alpha1AdditionalContentInline(),
				NewV1alpha1AdditionalContentObjectRef(),
				NewV1alpha1AdditionalContentWebhook(),
			},
		},
	}
}

func NewV1alpha1AdditionalContentInline() v1alpha1.AdditionalContent {
	return v1alpha1.AdditionalContent{
		Name:                inlineName,
		Exposed:             false,
		PostTransformations: nil,
		Inline:              types.Ptr(inlineContent),
	}
}

func NewV1alpha1AdditionalContentObjectRef() v1alpha1.AdditionalContent {
	return v1alpha1.AdditionalContent{
		Name:                objectRefName,
		Exposed:             false,
		PostTransformations: nil,
		ObjectRef: &v1alpha1.ObjectRef{
			ResourceRef: v1alpha1.ResourceRef{
				Group:     "core",
				Version:   "v1",
				Resource:  "ConfigMap",
				Namespace: "test-namespace",
				Name:      "test-cm",
			},
			JSONPath: ".data.jsonPath",
		},
	}
}

func NewV1alpha1AdditionalContentWebhook() v1alpha1.AdditionalContent {
	return v1alpha1.AdditionalContent{
		Name:                webhookName,
		Exposed:             false,
		PostTransformations: nil,
		Webhook: &v1alpha1.WebhookConfig{
			URL: "alexandre.mahdhaoui.com/s3-test",
			MTLSObjectRef: &v1alpha1.MTLSObjectRef{
				ResourceRef: v1alpha1.ResourceRef{
					Group:     "core",
					Version:   "v1",
					Resource:  "Secret",
					Namespace: "test-namespace",
					Name:      "test-mtls",
				},
				ClientKeyJSONPath:  ".data.\"client.key\"",
				ClientCertJSONPath: ".data.\"client.crt\"",
				CaBundleJSONPath:   ".data.\"ca.crt\"",
			},
			BasicAuthObjectRef: &v1alpha1.BasicAuthObjectRef{
				ResourceRef: v1alpha1.ResourceRef{
					Group:     "yoursecret.alexandre.mahdhaoui.com",
					Version:   "v1beta2",
					Resource:  "YourSecret",
					Namespace: "test-namespace",
					Name:      "test-custom-secret",
				},
				UsernameJSONPath: ".data.username",
				PasswordJSONPath: ".data.password",
			},
		},
	}
}

func NewTypesProfile() types.Profile {
	return types.Profile{
		IPXETemplate: ipxeTemplate,
		AdditionalContent: []types.Content{
			NewTypesContentInline(),
			NewTypesContentObjectRef(),
			NewTypesContentWebhookConfig(),
		},
	}
}

func NewTypesContentInline() types.Content {
	return types.Content{
		Name:             inlineName,
		PostTransformers: []types.TransformerConfig{},
		ResolverKind:     types.InlineResolverKind,
		Inline:           inlineContent,
	}
}

func NewTypesContentObjectRef() types.Content {
	return types.Content{
		Name:             objectRefName,
		PostTransformers: []types.TransformerConfig{},
		ResolverKind:     types.ObjectRefResolverKind,
		ObjectRef:        types.Ptr(NewTypesObjectRef()),
	}
}

func NewTypesObjectRef() types.ObjectRef {
	return types.ObjectRef{
		Group:     "core",
		Version:   "v1",
		Resource:  "ConfigMap",
		Namespace: "test-namespace",
		Name:      "test-cm",
		JSONPath:  &jsonpath.JSONPath{}, // to annoying
	}
}

func NewTypesContentWebhookConfig() types.Content {
	return types.Content{
		Name:             webhookName,
		PostTransformers: []types.TransformerConfig{},
		ResolverKind:     types.WebhookResolverKind,
		WebhookConfig:    types.Ptr(NewTypesWebhookConfig()),
	}
}

func NewTypesWebhookConfig() types.WebhookConfig {
	return types.WebhookConfig{
		URL: "alexandre.mahdhaoui.com/s3-test",
		MTLSObjectRef: &types.MTLSObjectRef{
			ObjectRef: types.ObjectRef{
				Group:     "core",
				Version:   "v1",
				Resource:  "Secret",
				Namespace: "test-namespace",
				Name:      "test-mtls",
				JSONPath:  nil,
			},
			ClientKeyJSONPath:  &jsonpath.JSONPath{}, // to annoying
			ClientCertJSONPath: &jsonpath.JSONPath{}, // to annoying
			CaBundleJSONPath:   &jsonpath.JSONPath{}, // to annoying
		},
		BasicAuthObjectRef: &types.BasicAuthObjectRef{
			ObjectRef: types.ObjectRef{
				Group:     "yoursecret.alexandre.mahdhaoui.com",
				Version:   "v1beta2",
				Resource:  "YourSecret",
				Namespace: "test-namespace",
				Name:      "test-custom-secret",
				JSONPath:  nil,
			},
			UsernameJSONPath: &jsonpath.JSONPath{}, // to annoying
			PasswordJSONPath: &jsonpath.JSONPath{}, // to annoying
		},
	}
}

func MakeContentComparable(content types.Content) types.Content {
	if content.ObjectRef != nil {
		content.ObjectRef.JSONPath = &jsonpath.JSONPath{}
	}

	if content.WebhookConfig != nil {
		if content.WebhookConfig.BasicAuthObjectRef != nil {
			content.WebhookConfig.BasicAuthObjectRef.UsernameJSONPath = &jsonpath.JSONPath{}
			content.WebhookConfig.BasicAuthObjectRef.PasswordJSONPath = &jsonpath.JSONPath{}
		}

		if content.WebhookConfig.MTLSObjectRef != nil {
			content.WebhookConfig.MTLSObjectRef.CaBundleJSONPath = &jsonpath.JSONPath{}
			content.WebhookConfig.MTLSObjectRef.ClientCertJSONPath = &jsonpath.JSONPath{}
			content.WebhookConfig.MTLSObjectRef.ClientKeyJSONPath = &jsonpath.JSONPath{}
		}
	}

	return content
}

func MakeProfileComparable(profile types.Profile) types.Profile {
	for i := range profile.AdditionalContent {
		profile.AdditionalContent[i] = MakeContentComparable(profile.AdditionalContent[i])
	}

	return profile
}
