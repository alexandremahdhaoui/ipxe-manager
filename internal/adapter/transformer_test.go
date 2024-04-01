//go:build unit

package adapter_test

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestButaneTransformer(t *testing.T) {
	var ()

	setup := func(t *testing.T) func() {
		t.Helper()

		return func() {
			t.Helper()
		}
	}

	t.Run("Transform", func(t *testing.T) {
		defer setup(t)()
	})
}

func WebhookTransformer(t *testing.T) {
	var (
		ctx      context.Context
		expected []byte

		inputConfig     types.TransformerConfig
		inputContent    []byte
		inputAttributes types.IpxeSelectors

		objectRefResolver *mockadapters.MockObjectRefResolver
		transformer       adapter.Transformer
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		expected = []byte("this has been templated: arm64, TODO(UUID)")

		inputConfig = types.TransformerConfig{
			Kind: types.WebhookTransformerKind,
			Webhook: &types.WebhookConfig{
				URL:                "test.example.com",
				MTLSObjectRef:      nil,
				BasicAuthObjectRef: nil,
			},
		}

		inputContent = []byte("this should be templated: buildarch, uuid")

		objectRefResolver = mockadapters.NewMockObjectRefResolver(t)
		adapter.NewWebhookTransformer(objectRefResolver)

		return func() {
			t.Helper()

			objectRefResolver.AssertExpectations(t)
		}
	}

	t.Run("Transform", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			objectRefResolver.EXPECT().ResolvePaths(mock.Anything, mock.Anything, mock.Anything)
			objectRefResolver.EXPECT().ResolvePaths(mock.Anything, mock.Anything, mock.Anything)

			//TODO: support input attributes for dynamic templating
			actual, err := transformer.Transform(ctx, inputConfig, inputContent, inputAttributes)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	})
}
