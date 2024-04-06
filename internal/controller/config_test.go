//go:build unit

package controller_test

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/controller"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockadapter"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockcontroller"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		ctx              context.Context
		inputProfileName string
		inputConfigID    uuid.UUID

		expectedProfileResult types.Profile
		expectedProfileErr    error

		expectedMuxResult map[string][]byte
		expectedMuxErr    error

		profile *mockadapter.MockProfile
		mux     *mockcontroller.MockResolveTransformerMux
		config  controller.Config
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		inputProfileName = "profile-name"
		inputConfigID = uuid.New()

		profile = mockadapter.NewMockProfile(t)
		mux = mockcontroller.NewMockResolveTransformerMux(t)
		config = controller.NewConfig(profile, mux)

		return func() {
			t.Helper()

			mux.AssertExpectations(t)
			expectedProfileErr = nil
			expectedMuxErr = nil
		}
	}

	expectProfile := func() {
		profile.EXPECT().
			Get(ctx, inputProfileName).
			Return(expectedProfileResult, expectedProfileErr).
			Once()
	}

	expectMux := func() {
		mux.EXPECT().
			ResolveAndTransformBatch(mock.Anything, mock.Anything, mock.Anything).
			Return(expectedMuxResult, expectedMuxErr).
			Once()

	}

	t.Run("GetByID", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expected := []byte("qwe")
			expectedProfileResult = types.Profile{
				IPXETemplate: "ipxe qwerty",
				AdditionalContent: []types.Content{{
					Name: "this will shouldn't be returned",
				}, {
					Name:            "this should be returned",
					ExposedConfigID: inputConfigID,
				}},
			}

			expectedMuxResult = map[string][]byte{
				expectedProfileResult.AdditionalContent[0].Name: []byte("asd"),
				expectedProfileResult.AdditionalContent[1].Name: expected,
			}

			expectProfile()
			expectMux()

			actual, err := config.GetByID(ctx, inputProfileName, inputConfigID)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("Config not found", func(t *testing.T) {
				defer setup(t)()

				expectedProfileResult = types.Profile{
					IPXETemplate: "ipxe qwerty",
					AdditionalContent: []types.Content{{
						Name: "this will shouldn't be returned",
					}},
				}

				expectProfile()

				_, err := config.GetByID(ctx, inputProfileName, inputConfigID)
				assert.ErrorIs(t, err, controller.ErrConfigNotFound)
			})

			t.Run("Profile Err", func(t *testing.T) {
				defer setup(t)()

				expectedProfileErr = assert.AnError
				expectProfile()

				_, err := config.GetByID(ctx, inputProfileName, inputConfigID)
				assert.ErrorIs(t, err, expectedProfileErr)
			})

			t.Run("Mux Err", func(t *testing.T) {
				defer setup(t)()

				expectedProfileResult = types.Profile{
					IPXETemplate: "ipxe qwerty",
					AdditionalContent: []types.Content{{
						Name: "this will shouldn't be returned",
					}, {
						Name:            "this should be returned",
						ExposedConfigID: inputConfigID,
					}},
				}

				expectedMuxErr = assert.AnError

				expectProfile()
				expectMux()

				_, err := config.GetByID(ctx, inputProfileName, inputConfigID)
				assert.ErrorIs(t, err, expectedMuxErr)
			})
		})
	})
}
