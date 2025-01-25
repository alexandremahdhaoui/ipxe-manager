//go:build unit

package controller_test

import (
	"context"
	"testing"

	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapter"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockcontroller"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	mustBeFilteredOut = "mustBeFilteredOut"
	mustBeReturned    = "mustBeReturned"
)

func TestContent(t *testing.T) {
	var (
		ctx context.Context

		inputConfigID uuid.UUID
		ipxeSelectors types.IPXESelectors

		expectedProfileResult []types.Profile
		expectedProfileErr    error

		expectedMuxResult []byte
		expectedMuxErr    error

		profile *mockadapter.MockProfile
		mux     *mockcontroller.MockResolveTransformerMux
		content controller.Content
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()

		inputConfigID = uuid.New()
		ipxeSelectors = types.IPXESelectors{}

		profile = mockadapter.NewMockProfile(t)
		mux = mockcontroller.NewMockResolveTransformerMux(t)
		content = controller.NewContent(profile, mux)

		expectedProfileResult = nil
		expectedProfileErr = nil

		expectedMuxResult = nil
		expectedMuxErr = nil

		return func() {
			t.Helper()

			profile.AssertExpectations(t)
			mux.AssertExpectations(t)
		}
	}

	expectProfile := func() {
		profile.EXPECT().
			ListByContentID(ctx, inputConfigID).
			Return(expectedProfileResult, expectedProfileErr).
			Once()
	}

	expectMux := func() {
		mux.EXPECT().
			ResolveAndTransform(mock.Anything, mock.Anything, mock.Anything).
			Return(expectedMuxResult, expectedMuxErr).
			Once()
	}

	t.Run("GetByID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expected := []byte("qwe")
			expectedProfileResult = []types.Profile{
				{
					IPXETemplate: "ipxe qwerty",
					AdditionalContent: map[string]types.Content{
						mustBeFilteredOut: {
							Name: mustBeFilteredOut,
						},
						mustBeReturned: {
							Name:        mustBeReturned,
							ExposedUUID: inputConfigID,
						},
					},
				},
			}

			expectedMuxResult = expected

			expectProfile()
			expectMux()

			actual, err := content.GetByID(ctx, inputConfigID, types.IPXESelectors{})
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("Content not found", func(t *testing.T) {
				defer setup(t)()

				expectedProfileResult = nil // no results
				expectProfile()

				_, err := content.GetByID(ctx, inputConfigID, ipxeSelectors)
				assert.ErrorIs(t, err, controller.ErrContentNotFound)
			})

			t.Run("Profile Err", func(t *testing.T) {
				defer setup(t)()

				expectedProfileErr = assert.AnError
				expectProfile()

				_, err := content.GetByID(ctx, inputConfigID, ipxeSelectors)
				assert.ErrorIs(t, err, expectedProfileErr)
			})

			t.Run("Mux Err", func(t *testing.T) {
				defer setup(t)()

				expectedProfileResult = []types.Profile{
					{
						IPXETemplate: "ipxe qwerty",
						AdditionalContent: map[string]types.Content{
							mustBeFilteredOut: {
								Name: mustBeFilteredOut,
							},
							mustBeReturned: {
								Name:        mustBeReturned,
								ExposedUUID: inputConfigID,
							},
						},
					},
				}

				expectedMuxErr = assert.AnError

				expectProfile()
				expectMux()

				_, err := content.GetByID(ctx, inputConfigID, ipxeSelectors)
				assert.ErrorIs(t, err, expectedMuxErr)
			})
		})
	})
}
