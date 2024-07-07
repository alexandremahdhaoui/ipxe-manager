//go:build unit

package adapter_test

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"testing"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
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

		inputBuildarch                 string
		expectedBuildarchLabelSelector string

		expectedAssignment  types.Assignment
		expectedListOptions []interface{}

		cl         *mockclient.MockClient
		assignment adapter.Assignment
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		namespace = "test-assignment"

		inputBuildarch = string(v1alpha1.Arm64)
		expectedBuildarchLabelSelector = v1alpha1.Arm64BuildarchLabelSelector

		cl = mockclient.NewMockClient(t)
		assignment = adapter.NewAssignment(cl, namespace)

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
					ProfileName: expectedAssignment.ProfileName,
				}}}

				return nil
			})
	}

	t.Run("FindDefaultByBuildarch", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expectedAssignment = types.Assignment{
				Name:        "",
				ProfileName: uuid.New().String(),
			}

			expectedListOptions = []interface{}{
				client.HasLabels{expectedBuildarchLabelSelector},
				client.HasLabels{v1alpha1.DefaultAssignmentLabel},
			}

			list(t)

			actual, err := assignment.FindDefaultByBuildarch(ctx, inputBuildarch)
			assert.NoError(t, err)
			assert.Equal(t, expectedAssignment, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("ListError", func(t *testing.T) {
				defer setup(t)()

				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

				actual, err := assignment.FindDefaultByBuildarch(ctx, inputBuildarch)
				assert.ErrorIs(t, err, assert.AnError)
				assert.Empty(t, actual)
			})

			t.Run("NotFound", func(t *testing.T) {
				defer setup(t)()

				// No assignment found.
				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				actual, err := assignment.FindDefaultByBuildarch(ctx, inputBuildarch)
				assert.ErrorIs(t, err, adapter.ErrAssignmentNotFound)
				assert.Empty(t, actual)
			})
		})
	})

	t.Run("FindBySelectors", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expectedAssignment = types.Assignment{
				Name:        "",
				ProfileName: uuid.New().String(),
			}

			id := uuid.New()
			selectors := types.IpxeSelectors{
				UUID:      id,
				Buildarch: inputBuildarch,
			}

			expectedListOptions = []any{
				client.HasLabels{expectedBuildarchLabelSelector},
				client.HasLabels{v1alpha1.NewUUIDLabelSelector(id)},
			}

			list(t)

			actual, err := assignment.FindBySelectors(ctx, selectors)
			assert.NoError(t, err)
			assert.Equal(t, expectedAssignment, actual)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("ListError", func(t *testing.T) {
				defer setup(t)()

				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError)

				actual, err := assignment.FindBySelectors(ctx, types.IpxeSelectors{})
				assert.ErrorIs(t, err, assert.AnError)
				assert.Empty(t, actual)
			})

			t.Run("NotFound", func(t *testing.T) {
				defer setup(t)()

				// No assignment found.
				cl.EXPECT().List(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				actual, err := assignment.FindBySelectors(ctx, types.IpxeSelectors{})
				assert.ErrorIs(t, err, adapter.ErrAssignmentNotFound)
				assert.Empty(t, actual)
			})
		})
	})
}
