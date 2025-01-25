//go:build unit

package controller_test

import (
	"context"
	"fmt"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapter"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestResolveTransformerMux(t *testing.T) {
	var (
		ctx            context.Context
		inputSelectors types.IpxeSelectors
		inputBatch     []types.Content

		inlineResolver    *mockadapter.MockResolver
		objectRefResolver *mockadapter.MockResolver
		webhookResolver   *mockadapter.MockResolver

		butaneTransformer  *mockadapter.MockTransformer
		webhookTransformer *mockadapter.MockTransformer

		resolvers    map[types.ResolverKind]adapter.Resolver
		transformers map[types.TransformerKind]adapter.Transformer

		mux controller.ResolveTransformerMux
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		inputSelectors = types.IpxeSelectors{
			UUID:      uuid.New(),
			Buildarch: "arm64",
		}

		inputBatch = make([]types.Content, 0)

		inlineResolver = mockadapter.NewMockResolver(t)
		objectRefResolver = mockadapter.NewMockResolver(t)
		webhookResolver = mockadapter.NewMockResolver(t)

		butaneTransformer = mockadapter.NewMockTransformer(t)
		webhookTransformer = mockadapter.NewMockTransformer(t)

		resolvers = map[types.ResolverKind]adapter.Resolver{
			types.InlineResolverKind:    inlineResolver,
			types.ObjectRefResolverKind: objectRefResolver,
			types.WebhookResolverKind:   webhookResolver,
		}

		transformers = map[types.TransformerKind]adapter.Transformer{
			types.ButaneTransformerKind:  butaneTransformer,
			types.WebhookTransformerKind: webhookTransformer,
		}

		mux = controller.NewResolveTransformerMux(resolvers, transformers)

		return func() {
			t.Helper()

			inlineResolver.AssertExpectations(t)
			objectRefResolver.AssertExpectations(t)
			webhookResolver.AssertExpectations(t)

			butaneTransformer.AssertExpectations(t)
			webhookTransformer.AssertExpectations(t)
		}
	}

	setup(t)()

	t.Run("ResolveAndTransformBatch", func(t *testing.T) {
		for kind := range resolvers {

			t.Run(resolverKindString(t, kind), func(t *testing.T) {
				t.Run("Success", func(t *testing.T) {
					defer setup(t)()

					expected := make(map[string][]byte)

					// generate n content for the batch
					for i := 0; i < 3; i++ {
						inputContent := types.Content{
							Name: fmt.Sprintf("%s-%d", t.Name(), i),
							PostTransformers: []types.TransformerConfig{{
								Kind: types.ButaneTransformerKind,
							}, {
								Kind:    types.WebhookTransformerKind,
								Webhook: types.Ptr(testutil.NewTypesWebhookConfig()),
							}},
							ResolverKind:  kind,
							Inline:        "this is an inline content",
							ObjectRef:     types.Ptr(testutil.NewTypesObjectRef()),
							WebhookConfig: types.Ptr(testutil.NewTypesWebhookConfig()),
						}

						inputBatch = append(inputBatch, inputContent)

						expectedResolverResult := []byte("expectedResolverResult")
						expectedTransformationResult0 := []byte("expectedTransformationResult0")
						expectedTransformationResult1 := []byte("expectedTransformationResult1")
						expected[inputContent.Name] = expectedTransformationResult1

						resolvers[kind].(*mockadapter.MockResolver).EXPECT().
							Resolve(ctx, inputContent).
							Return(expectedResolverResult, nil).
							Once()

						butaneTransformer.EXPECT().
							Transform(ctx, inputContent.PostTransformers[0], expectedResolverResult, inputSelectors).
							Return(expectedTransformationResult0, nil).
							Once()

						webhookTransformer.EXPECT().
							Transform(ctx, inputContent.PostTransformers[1], expectedTransformationResult0, inputSelectors).
							Return(expectedTransformationResult1, nil).
							Once()
					}

					actual, err := mux.ResolveAndTransformBatch(ctx, inputBatch, inputSelectors)
					assert.NoError(t, err)
					assert.Equal(t, expected, actual)
				})
			})
		}

		t.Run("Failure", func(t *testing.T) {
			t.Run("unknown resolver", func(t *testing.T) {
				defer setup(t)()

				inputBatch = append(inputBatch, types.Content{
					Name:         t.Name(),
					ResolverKind: -1,
				})

				_, err := mux.ResolveAndTransformBatch(ctx, inputBatch, inputSelectors)
				assert.ErrorIs(t, err, controller.ErrResolverUnknown)
			})

			t.Run("unknown transformer", func(t *testing.T) {
				defer setup(t)()

				inputBatch = append(inputBatch, types.Content{
					Name: t.Name(),
					PostTransformers: []types.TransformerConfig{{
						Kind: -1,
					}},
					ResolverKind: types.InlineResolverKind,
				})

				resolvers[types.InlineResolverKind].(*mockadapter.MockResolver).EXPECT().
					Resolve(mock.Anything, mock.Anything).
					Return([]byte("whatever"), nil).
					Once()

				_, err := mux.ResolveAndTransformBatch(ctx, inputBatch, inputSelectors)
				assert.ErrorIs(t, err, controller.ErrTransformerUnknown)
			})

			t.Run("resolver error", func(t *testing.T) {
				defer setup(t)()

				inputBatch = append(inputBatch, types.Content{
					Name:         t.Name(),
					ResolverKind: types.InlineResolverKind,
				})

				resolvers[types.InlineResolverKind].(*mockadapter.MockResolver).EXPECT().
					Resolve(mock.Anything, mock.Anything).
					Return(nil, assert.AnError).
					Once()

				_, err := mux.ResolveAndTransformBatch(ctx, inputBatch, inputSelectors)
				assert.ErrorIs(t, err, assert.AnError)
			})

			t.Run("transformer error", func(t *testing.T) {
				defer setup(t)()

				inputBatch = append(inputBatch, types.Content{
					Name:         t.Name(),
					ResolverKind: types.InlineResolverKind,
					PostTransformers: []types.TransformerConfig{{
						Kind: types.ButaneTransformerKind,
					}},
				})

				resolvers[inputBatch[0].ResolverKind].(*mockadapter.MockResolver).EXPECT().
					Resolve(mock.Anything, mock.Anything).
					Return([]byte("whatever"), nil).
					Once()

				transformers[inputBatch[0].PostTransformers[0].Kind].(*mockadapter.MockTransformer).EXPECT().
					Transform(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, assert.AnError).
					Once()

				_, err := mux.ResolveAndTransformBatch(ctx, inputBatch, inputSelectors)
				assert.ErrorIs(t, err, assert.AnError)
			})
		})
	})
}

func resolverKindString(t *testing.T, kind types.ResolverKind) string {
	t.Helper()

	switch kind {
	case types.InlineResolverKind:
		return "InlineResolver"
	case types.ObjectRefResolverKind:
		return "ObjectRefResolver"
	case types.WebhookResolverKind:
		return "WebhookResolver"
	default:
		panic("abort")
	}
}
