//go:build unit

package adapter_test

import (
	"context"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/alexandremahdhaoui/ipxer/pkg/resolverserver"
	"github.com/labstack/echo/v4"
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

		webhookResolverServer resolverserver.ServerInterface

		cl       *fake.FakeDynamicClient
		resolver adapter.Resolver
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		expected = []byte("expected additional content")

		// -------------------------------------------------- Content ----------------------------------------------- //
		content = testutil.NewTypesContentWebhookConfig()
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.UsernameJSONPath.Parse(`{.data.username}`))
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.PasswordJSONPath.Parse(`{.data.password}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientKeyJSONPath.Parse(`{.data.client\.key}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientCertJSONPath.Parse(`{.data.client\.crt}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.CaBundleJSONPath.Parse(`{.data.ca\.crt}`))

		// -------------------------------------------------- Basic Auth -------------------------------------------- //

		basicAuthObject = &unstructured.Unstructured{}
		basicAuthObject.SetUnstructuredContent(map[string]any{"data": map[string]any{
			"username": "qwe123",
			"password": "321ewq",
		}})
		basicAuthObject.SetName(content.WebhookConfig.BasicAuthObjectRef.Name)
		basicAuthObject.SetNamespace(content.WebhookConfig.BasicAuthObjectRef.Namespace)
		basicAuthObject.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "yoursecret.alexandre.mahdhaoui.com",
			Version: "v1beta2",
			Kind:    "YourSecret",
		})

		// -------------------------------------------------- mTLS  ------------------------------------------------- //

		//TODO: create func or use existing library for creating self signed client/server cert.

		mtlsObject = &unstructured.Unstructured{}
		mtlsObject.SetUnstructuredContent(map[string]any{"data": map[string]any{
			"client.key": "TODO", //TODO
			"client.crt": "TODO", //TODO
			"ca.crt":     "TODO", //TODO
		}})

		mtlsObject.SetName(content.WebhookConfig.MTLSObjectRef.Name)
		mtlsObject.SetNamespace(content.WebhookConfig.MTLSObjectRef.Namespace)
		mtlsObject.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "core",
			Version: "v1",
			Kind:    "Secret",
		})

		// -------------------------------------------------- Webhook Server  --------------------------------------- //

		//TODO: create a server to test:
		// - validate mTLS certs
		// - validate basic auth

		server := echo.New()
		webhookResolverServer = nil //TODO: create the dummy server interface.
		resolverserver.RegisterHandlersWithBaseURL(server, webhookResolverServer, content.WebhookConfig.URL)

		err := server.Start(fmt.Sprintf(":%d", testutil.WebhookServerPort))
		require.NoError(t, err)

		// -------------------------------------------------- Client and Adapter ------------------------------------ //

		cl = fake.NewSimpleDynamicClient(runtime.NewScheme(), basicAuthObject, mtlsObject)

		objectRefResolver := adapter.NewObjectRefResolver(cl)
		resolver = adapter.NewWebhookResolver(objectRefResolver)

		return func() {
			t.Helper()

			err := server.Shutdown(context.Background())
			require.NoError(t, err)
		}
	}

	t.Run("Resolve", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			setup(t)

			cl.PrependReactor("get", "YourSecret", func(_ k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, basicAuthObject, nil
			})

			cl.PrependReactor("get", "Secret", func(_ k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, mtlsObject, nil
			})

			actual, err := resolver.Resolve(ctx, content)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}
