//go:build unit

package adapter_test

import (
	"context"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	resolverserver2 "github.com/alexandremahdhaoui/ipxer/pkg/generated/resolverserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
	"strings"
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

		mock *resolverserver2.Mock

		cl       *fake.FakeDynamicClient
		resolver adapter.Resolver
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()

		// -------------------------------------------------- Content ----------------------------------------------- //

		content = testutil.NewTypesContentWebhook()
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.UsernameJSONPath.Parse(`{.data.username}`))
		require.NoError(t, content.WebhookConfig.BasicAuthObjectRef.PasswordJSONPath.Parse(`{.data.password}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientKeyJSONPath.Parse(`{.data.client\.key}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.ClientCertJSONPath.Parse(`{.data.client\.crt}`))
		require.NoError(t, content.WebhookConfig.MTLSObjectRef.CaBundleJSONPath.Parse(`{.data.ca\.crt}`))

		// -------------------------------------------------- Webhook Server  --------------------------------------- //

		addr := strings.SplitN(content.WebhookConfig.URL, "/", 2)[0]
		mock = resolverserver2.NewMock(t, addr)

		clientKey, clientCert, err := mock.CA.NewCertifiedKeyPEM(addr)
		require.NoError(t, err)
		caCert := mock.CA.Cert()

		// -------------------------------------------------- Basic Auth -------------------------------------------- //

		username, password := "qwe123", "321ewq"
		mock.Echo.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Validator: func(u string, p string, _ echo.Context) (bool, error) {
				return u == username && p == password, nil
			},
		}))

		basicAuthObject = &unstructured.Unstructured{}
		basicAuthObject.SetUnstructuredContent(map[string]any{"data": map[string]any{
			"username": username,
			"password": password,
		}})

		basicAuthObject.SetName(content.WebhookConfig.BasicAuthObjectRef.Name)
		basicAuthObject.SetNamespace(content.WebhookConfig.BasicAuthObjectRef.Namespace)
		basicAuthObject.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "yoursecret.alexandre.mahdhaoui.com",
			Version: "v1beta2",
			Kind:    "YourSecret",
		})

		// -------------------------------------------------- mTLS  ------------------------------------------------- //

		mtlsObject = &unstructured.Unstructured{}
		mtlsObject.SetUnstructuredContent(map[string]any{
			"data": map[string]any{
				"client.key": string(clientKey),
				"client.crt": string(clientCert),
				"ca.crt":     string(caCert),
			},
		})

		mtlsObject.SetName(content.WebhookConfig.MTLSObjectRef.Name)
		mtlsObject.SetNamespace(content.WebhookConfig.MTLSObjectRef.Namespace)
		mtlsObject.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "core",
			Version: "v1",
			Kind:    "Secret",
		})

		// -------------------------------------------------- Client and Adapter ------------------------------------ //

		cl = fake.NewSimpleDynamicClient(runtime.NewScheme(), basicAuthObject, mtlsObject)

		objectRefResolver := adapter.NewObjectRefResolver(cl)
		resolver = adapter.NewWebhookResolver(objectRefResolver)

		return func() {
			t.Helper()

			mock.AssertExpectationsAndShutdown()
		}
	}

	t.Run("Resolve", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			mock.AppendExpectation(func(e echo.Context, params resolverserver2.ResolveParams) error {
				expected = []byte(fmt.Sprintf("hello world: %s + %s", params.Buildarch, params.Uuid.String()))
				e.Response().Write(expected)
				return nil
			})

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

		t.Run("Fail", func(t *testing.T) {
			defer setup(t)()

			cl.PrependReactor("get", "YourSecret", func(_ k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				basicAuthObject.SetUnstructuredContent(map[string]any{"data": map[string]any{
					"username": "not a username",
					"password": "not a password",
				}})

				return true, basicAuthObject, nil
			})

			cl.PrependReactor("get", "Secret", func(_ k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, mtlsObject, nil
			})

			expected = []byte(`{"message":"Unauthorized"}`)
			expected = append(expected, byte(0x0a))

			actual, err := resolver.Resolve(ctx, content)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}
