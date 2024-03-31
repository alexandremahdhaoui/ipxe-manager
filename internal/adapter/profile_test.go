//go:build unit

package adapter_test

import (
	"context"
	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/util/mocks/mockclient"
	"github.com/alexandremahdhaoui/ipxer/internal/util/testutil"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	types2 "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestProfile(t *testing.T) {
	var (
		ctx       context.Context
		namespace string

		inputProfileName string

		v1alpha1Profile v1alpha1.Profile
		expectedErr     error

		cl      *mockclient.MockClient
		profile adapter.Profile
	)

	setup := func(t *testing.T) func() {
		t.Helper()

		ctx = context.Background()
		namespace = "test-profile"

		inputProfileName = "profile-name"

		v1alpha1Profile = testutil.NewV1alpha1Profile()

		cl = mockclient.NewMockClient(t)
		profile = adapter.NewProfile(cl, namespace)

		return func() {
			t.Helper()

			cl.AssertExpectations(t)
		}
	}

	get := func(t *testing.T) {
		t.Helper()

		cl.EXPECT().
			Get(ctx, types2.NamespacedName{
				Namespace: namespace,
				Name:      inputProfileName,
			}, mock.Anything).
			RunAndReturn(func(_ context.Context, _ types2.NamespacedName, obj client.Object, _ ...client.GetOption) error {
				p := obj.(*v1alpha1.Profile)
				*p = v1alpha1Profile

				return expectedErr
			})
	}

	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer setup(t)()

			expected := testutil.NewTypesProfile()

			get(t)

			actual, err := profile.Get(ctx, inputProfileName)
			assert.NoError(t, err)
			assert.Equal(t, expected, testutil.MakeProfileComparable(actual))
		})

		t.Run("Failure", func(t *testing.T) {
			t.Run("Get error", func(t *testing.T) {
				defer setup(t)()

				expectedErr = assert.AnError
				get(t)

				_, err := profile.Get(ctx, inputProfileName)
				assert.ErrorIs(t, err, assert.AnError)
			})
		})
	})
}
