//go:build unit

package adapter_test

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
	"testing"
)

func TestInlineResolver(t *testing.T) {
	var (
		resolver adapter.Resolver
	)

	setup := func() {
		resolver = adapter.NewInlineResolver()
	}

	t.Run("Resolve", func(t *testing.T) {
		setup()

		expected := []byte("test")
		input := types.Content{
			Inline: string(expected),
		}

		out, err := resolver.Resolve(nil, input)
		assert.NoError(t, err)
		assert.Equal(t, expected, out)
	})
}

func TestObjectRefResolver(t *testing.T) {
	var (
		ctx context.Context

		expected []byte
		content  types.Content
		object   *unstructured.Unstructured

		cl       *fake.FakeDynamicClient
		resolver adapter.Resolver
	)

	setup := func(t *testing.T) {
		t.Helper()

		ctx = context.Background()

		expected = []byte("qwe")

		content = testutil.NewTypesContentObjectRef()
		require.NoError(t, content.ObjectRef.JSONPath.Parse("{.data.test}"))

		object = &unstructured.Unstructured{}
		object.SetName(content.ObjectRef.Name)
		object.SetNamespace(content.ObjectRef.Namespace)
		object.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   content.ObjectRef.Group,
			Version: content.ObjectRef.Version,
			Kind:    content.ObjectRef.Resource,
		})

		cl = fake.NewSimpleDynamicClient(runtime.NewScheme(), object)
		resolver = adapter.NewObjectRefResolver(cl)
	}

	t.Run("Resolve", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			setup(t)

			object.SetUnstructuredContent(map[string]any{"data": map[string]any{"test": string(expected)}})
			cl.PrependReactor("get", "ConfigMap", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, object, nil
			})

			actual, err := resolver.Resolve(ctx, content)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}

func TestWebhookResolver(t *testing.T) {
	var (
		ctx context.Context

		expected []byte

		basicAuthObject *unstructured.Unstructured
		mtlsObject      *unstructured.Unstructured
		content         types.Content

		cl       *fake.FakeDynamicClient
		resolver adapter.Resolver
	)

	setup := func(t *testing.T) {
		t.Helper()

		expected = []byte("expected additional content")

		basicAuthObject = &unstructured.Unstructured{}
		mtlsObject = &unstructured.Unstructured{}

		content = testutil.NewTypesContentWebhookConfig()
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.UsernameJSONPath.Parse(`{.data.username}`))
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.PasswordJSONPath.Parse(`{.data.password}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientKeyJSONPath.Parse(`{.data."client.key"}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientCertJSONPath.Parse(`{.data."client.cert"}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.CaBundleJSONPath.Parse(`{.data."ca.bundle"}`))

		cl = fake.NewSimpleDynamicClient(runtime.NewScheme(), basicAuthObject, mtlsObject)
		objectRefResolver := adapter.NewObjectRefResolver(cl)
		resolver = adapter.NewWebhookResolver(objectRefResolver)
	}

	//TODO: finish this test

	t.Run("Resolve", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			setup(t)

			actual, err := resolver.Resolve(ctx, content)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}
