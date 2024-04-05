//go:build unit

package adapter_test

import (
	"context"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapters"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/alexandremahdhaoui/ipxer/pkg/transformerserver"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestButaneTransformer(t *testing.T) {
	var (
		transformer adapter.Transformer
	)

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

		inputSelectors := types.IpxeSelectors{
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
		expected []byte

		inputConfig     types.TransformerConfig
		inputContent    []byte
		inputAttributes types.IpxeSelectors

		objectRefResolver *mockadapters.MockObjectRefResolver
		transformer       adapter.Transformer
		serverMock        *transformerserver.Mock
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		id := uuid.New()
		buildarch := "arm64"
		expected = []byte(fmt.Sprintf("this has been templated: %s, %s", id.String(), buildarch))

		// -------------------------------------------------- Inputs ------------------------------------------------ //

		inputConfig = testutil.NewTypesTransformerConfigWebhook()
		inputContent = []byte("this should be templated: {{ .uuid }}, {{ .buildarch }}")
		inputAttributes = types.IpxeSelectors{
			UUID:      id,
			Buildarch: buildarch,
		}

		// -------------------------------------------------- Client and Adapter ------------------------------------ //

		objectRefResolver = mockadapters.NewMockObjectRefResolver(t)
		transformer = adapter.NewWebhookTransformer(objectRefResolver)

		// -------------------------------------------------- Webhook Server Mock ----------------------------------- //

		addr := strings.SplitN(inputConfig.Webhook.URL, "/", 2)[0]
		serverMock = transformerserver.NewMock(t, addr)

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
		serverMock.Echo.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Validator: func(u string, p string, _ echo.Context) (bool, error) {
				return u == username && p == password, nil
			},
		}))

		objectRefResolver.EXPECT().
			ResolvePaths(mock.Anything, mock.Anything, mock.Anything).
			Return([][]byte{[]byte(username), []byte(password)}, nil).
			Once()

		// -------------------------------------------------- Teardown  --------------------------------------------- //

		return func() {
			t.Helper()

			objectRefResolver.AssertExpectations(t)
			serverMock.AssertExpectationsAndShutdown()
		}
	}

	t.Run("Transform", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			serverMock.AppendExpectation(func(e echo.Context) error {
				_, err := e.Response().Write(expected)
				require.NoError(t, err)
				return nil
			})

			actual, err := transformer.Transform(ctx, inputConfig, inputContent, inputAttributes)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			defer setup(t)()

			expected = []byte(`{"code":400,"message":"error"}`)

			serverMock.AppendExpectation(func(e echo.Context) error {
				e.Response().Status = 400
				e.Response().Write(expected)
				return nil
			})

			actual, err := transformer.Transform(ctx, inputConfig, inputContent, inputAttributes)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}
