//go:build unit

package adapters_test

import (
	"context"
	"testing"

	"github.com/alexandremahdhaoui/ipxer/internal/adapters"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockclient"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAssignment(t *testing.T) {
	var (
		ctx       context.Context
		namespace string

		inputBuildarch string

		expectedProfile     string
		expectedListOptions []interface{}

		cl         *mockclient.MockClient
		assignment adapters.Assignment
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		namespace = "test-assignment"

		inputBuildarch = "arm64"

		cl = mockclient.NewMockClient(t)
		assignment = adapters.NewAssignment(cl, namespace)

		return func() {
			t.Helper()

			cl.AssertExpectations(t)
		}
	}

	list := func(t *testing.T) {
		t.Helper()

		cl.EXPECT().List(ctx, mock.Anything, expectedListOptions...).
			RunAndReturn(func(_ context.Context, objList client.ObjectList, options ...client.ListOption) error {
				l := objList.(*v1alpha1.AssignmentList)
				l.Items = []v1alpha1.Assignment{{Spec: v1alpha1.AssignmentSpec{
					ProfileName: expectedProfile,
				}}}

				return nil
			})
	}

	t.Run("FindDefaultProfile", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expectedProfile = uuid.New().String()
			expectedListOptions = []interface{}{
				client.MatchingLabels{v1alpha1.BuildarchAssignmentLabel: inputBuildarch},
				client.HasLabels{v1alpha1.DefaultAssignmentLabel},
			}

			list(t)

			actual, err := assignment.FindDefaultProfile(ctx, inputBuildarch)
			assert.NoError(t, err)
			assert.Equal(t, expectedProfile, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("ListError", func(t *testing.T) {
				defer setup(t)()

				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

				actual, err := assignment.FindDefaultProfile(ctx, inputBuildarch)
				assert.ErrorIs(t, err, assert.AnError)
				assert.Empty(t, actual)
			})

			t.Run("NotFound", func(t *testing.T) {
				defer setup(t)()

				// No assignment found.
				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				actual, err := assignment.FindDefaultProfile(ctx, inputBuildarch)
				assert.ErrorIs(t, err, adapters.ErrAssignmentNotFound)
				assert.Empty(t, actual)
			})
		})
	})

	t.Run("FindProfileBySelectors", func(t *testing.T) {
	})
}
