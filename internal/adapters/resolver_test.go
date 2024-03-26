//go:build unit

package adapters_test

import "testing"

func TestInlineResolver(t *testing.T) {
	var ()

	setup := func(t *testing.T) func() {
		t.Helper()
		return func() {
			t.Helper()
		}
	}

	t.Run("Resolve", func(t *testing.T) {
		defer setup(t)()
	})
}

func TestObjectRefResolver(t *testing.T) {
	var ()

	setup := func(t *testing.T) func() {
		t.Helper()
		return func() {
			t.Helper()
		}
	}

	t.Run("Resolve", func(t *testing.T) {
		defer setup(t)()
	})
}

func TestWebhookResolver(t *testing.T) {
	var ()

	setup := func(t *testing.T) func() {
		t.Helper()
		return func() {
			t.Helper()
		}
	}

	t.Run("Resolve", func(t *testing.T) {
		defer setup(t)()
	})
}
