//go:build unit

package adapter_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/fakes/transformerserverfake"
	"github.com/alexandremahdhaoui/ipxer/internal/util/httputil"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapter"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/alexandremahdhaoui/ipxer/pkg/generated/transformerserver"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestButaneTransformer(t *testing.T) {
	var transformer adapter.Transformer

	setup := func(t *testing.T) {
		t.Helper()

		transformer = adapter.NewButaneTransformer()
	}

	t.Run("Transform", func(t *testing.T) {
		setup(t)

		inputCfg := types.TransformerConfig{Kind: types.ButaneTransformerKind}
		inputContent := []byte(`
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
`)

		inputSelectors := types.IPXESelectors{
			UUID:      uuid.New(),
			Buildarch: "arm64",
		}

		expected := []byte(`{"ignition":{"version":"3.4.0"},"passwd":{"users":[{"name":"core"}]}}`)

		ctx := context.Background()
		actual, err := transformer.Transform(ctx, inputCfg, inputContent, inputSelectors)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestWebhookTransformer(t *testing.T) {
	var (
		ctx      context.Context
		expected string

		inputConfig     types.TransformerConfig
		inputContent    []byte
		inputAttributes types.IPXESelectors

		objectRefResolver *mockadapter.MockObjectRefResolver
		transformer       adapter.Transformer
		serverMock        *transformerserverfake.Fake
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		id := uuid.New() //nolint:varnamelen
		buildarch := "arm64"
		expected = fmt.Sprintf("this has been templated: %s, %s", id.String(), buildarch)

		// -------------------------------------------------- Inputs ------------------------------------------------ //

		inputConfig = testutil.NewTypesTransformerConfigWebhook()
		inputContent = []byte("this should be templated: {{ .uuid }}, {{ .buildarch }}")
		inputAttributes = types.IPXESelectors{
			UUID:      id,
			Buildarch: buildarch,
		}

		// -------------------------------------------------- Client and Adapter ------------------------------------ //

		objectRefResolver = mockadapter.NewMockObjectRefResolver(t)
		transformer = adapter.NewWebhookTransformer(objectRefResolver)

		// -------------------------------------------------- Webhook Server Fake ----------------------------------- //

		addr := strings.SplitN(inputConfig.Webhook.URL, "/", 2)[0]
		serverMock = transformerserverfake.New(t, addr)

		clientKey, clientCert, err := serverMock.CA.NewCertifiedKeyPEM(addr)
		require.NoError(t, err)
		caCert := serverMock.CA.Cert()

		// -------------------------------------------------- mTLS  ------------------------------------------------- //

		objectRefResolver.EXPECT().
			ResolvePaths(mock.Anything, mock.Anything, mock.Anything).
			Return([][]byte{clientKey, clientCert, caCert}, nil).
			Once()

		// -------------------------------------------------- Basic Auth -------------------------------------------- //

		username, password := "qwe123", "321ewq"

		currentHandler := serverMock.Server.Handler
		httputil.BasicAuth(currentHandler, func(u, p string, _ *http.Request) (bool, error) {
			return u == username && p == password, nil
		})

		objectRefResolver.EXPECT().
			ResolvePaths(mock.Anything, mock.Anything, mock.Anything).
			Return([][]byte{[]byte(username), []byte(password)}, nil).
			Once()

		// -------------------------------------------------- Teardown  --------------------------------------------- //

		return func() { //nolint:contextcheck
			t.Helper()

			objectRefResolver.AssertExpectations(t)
			serverMock.AssertExpectationsAndShutdown()
		}
	}

	t.Run("Transform", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expected = fmt.Sprintf("{\"data\":\"%s + %s\"}\n", inputAttributes.Buildarch, inputAttributes.UUID.String())

			serverMock.AppendExpectation(func(_ context.Context, request transformerserver.TransformRequestObject) (transformerserver.TransformResponseObject, error) { //nolint:lll
				t.Helper()

				return transformerserver.Transform200JSONResponse{
					TransformRespJSONResponse: transformerserver.TransformRespJSONResponse{
						Data: ptr.To(fmt.Sprintf("%s + %s", request.Body.Attributes.Buildarch, request.Body.Attributes.Uuid.String())), //nolint:lll
					},
				}, nil
			})

			actual, err := transformer.Transform(ctx, inputConfig, inputContent, inputAttributes)
			require.NoError(t, err)
			assert.Equal(t, expected, string(actual))
		})

		t.Run("Failure", func(t *testing.T) {
			defer setup(t)()

			expected = "{\"code\":400,\"message\":\"error\"}\n"

			serverMock.AppendExpectation(func(_ context.Context, _ transformerserver.TransformRequestObject) (transformerserver.TransformResponseObject, error) { //nolint:lll
				t.Helper()

				return transformerserver.Transform400JSONResponse{
					N400JSONResponse: transformerserver.N400JSONResponse{
						Code:    400,
						Message: "error",
					},
				}, nil
			})

			actual, err := transformer.Transform(ctx, inputConfig, inputContent, inputAttributes)
			require.NoError(t, err)
			assert.Equal(t, expected, string(actual))
		})
	})
}
