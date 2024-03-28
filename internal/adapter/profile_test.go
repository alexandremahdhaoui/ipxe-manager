//go:build unit

package adapter_test

import (
	"testing"
)

func TestProfile(t *testing.T) {
	var ()

	setup := func(t *testing.T) func() {
		return func() {
			t.Helper()
		}
	}

	t.Run("Get", func(t *testing.T) {
		defer setup(t)()
	})
}
