//go:build unit

package adapters_test

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/adapters"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockclient"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestProfile(t *testing.T) {
	namespace := "test-profile"
	buildarch := "arm64"

	var (
		ctx context.Context

		id                 uuid.UUID
		assignments        []v1alpha1.Assignment
		defaultAssignments []v1alpha1.Assignment

		expectedErrList        error
		expectedErrListDefault error
		expectedErrGet         error

		c *mockclient.MockClient
		p adapters.Profile
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()

		id = uuid.New()
		assignments = nil
		defaultAssignments = nil

		expectedErrList = nil
		expectedErrListDefault = nil
		expectedErrGet = nil

		c = mockclient.NewMockClient(t)
		p = adapters.NewProfile(c, namespace)

		return func() {
			t.Helper()

			c.AssertExpectations(t)
		}
	}

	list := func(t *testing.T) {
		t.Helper()

		c.EXPECT().
			List(mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, list client.ObjectList, _ ...client.ListOption) error {
				asl := list.(*v1alpha1.AssignmentList)
				asl.Items = assignments

				return expectedErrList
			}).Once()
	}

	listDefault := func(t *testing.T) {
		t.Helper()

		c.EXPECT().
			List(mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, list client.ObjectList, _ ...client.ListOption) error {
				asl := list.(*v1alpha1.AssignmentList)
				asl.Items = defaultAssignments

				return expectedErrListDefault
			}).Once()
	}

	get := func(t *testing.T) {
		t.Helper()

		c.EXPECT().
			Get(mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, name k8stypes.NamespacedName, obj client.Object, _ ...client.GetOption) error {
				profile := obj.(*v1alpha1.Profile)
				return expectedErrGet
			})
	}

	t.Run("FindBySelectors", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Run("found", func(t *testing.T) {
				defer setup(t)()
				list(t)
				get(t)

				selectors := types.IpxeSelectors{
					UUID:      id,
					Buildarch: buildarch,
				}

				actual, err := p.FindBySelectors(ctx, selectors)
				assert.NoError(t, err)
			})

			t.Run("default assignment", func(t *testing.T) {
				defer setup(t)()
				list(t)
				listDefault(t)
				get(t)

			})
		})

		t.Run("err", func(t *testing.T) {
			t.Run("list assignment", func(t *testing.T) {
				defer setup(t)()
				expectedErrList = errors.New("expected err list")
				list(t)

				selectors := types.IpxeSelectors{}

				_, err := p.FindBySelectors(ctx, selectors)
				assert.ErrorIs(t, err, adapters.ErrFindAssignmentBySelectors)
				assert.ErrorIs(t, err, expectedErrList)
			})

			t.Run("list default assignment", func(t *testing.T) {
				defer setup(t)()
				list(t)

				expectedErrListDefault = errors.New("expected err list default")
				listDefault(t)

				_, err := p.FindBySelectors(ctx, types.IpxeSelectors{})
				assert.ErrorIs(t, err, adapters.ErrFindAssignmentBySelectors)
				assert.ErrorIs(t, err, expectedErrListDefault)
			})

			t.Run("cannot find default assignment", func(t *testing.T) {
				defer setup(t)()
				list(t)
				listDefault(t)

				_, err := p.FindBySelectors(ctx, types.IpxeSelectors{})
				assert.ErrorIs(t, err, adapters.ErrFindAssignmentBySelectors)
				assert.ErrorIs(t, err, adapters.ErrProfileNotFound)
			})

			t.Run("get profile", func(t *testing.T) {
				defer setup(t)()
				list(t)
				listDefault(t)

				expectedErrGet = errors.New("expected err get")
				get(t)

				_, err := p.FindBySelectors(ctx, types.IpxeSelectors{})
				assert.ErrorIs(t, err, adapters.ErrFindAssignmentBySelectors)
				assert.ErrorIs(t, err, expectedErrGet)
			})

			t.Run("convert to profile", func(t *testing.T) {
				t.Run("success", func(t *testing.T) {

				})

				t.Run("err", func(t *testing.T) {
					t.Run("get profile", func(t *testing.T) {
						defer setup(t)()

						expectedErrGet = errors.New("expected err get")
						get(t)

						_, err := p.Get(ctx, id)
						assert.ErrorIs(t, err, adapters.errProfileGet)
						assert.ErrorIs(t, err, expectedErrGet)
					})
				})
			})
		})
	})

	t.Run("Get", func(t *testing.T) {

	})
}
