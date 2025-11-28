package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureWebsiteExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_ensure_website")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	existingDir := filepath.Join(tempDir, "existing")
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create existing dir: %v", err)
	}

	nonExistentDir := filepath.Join(tempDir, "non_existent")

	regenerateCalled := false
	regenerate := func() error {
		regenerateCalled = true
		return nil
	}

	t.Run("Directory exists, watch=false, no regeneration", func(t *testing.T) {
		regenerateCalled = false
		err := ensureWebsiteExists(existingDir, false, regenerate)
		assert.NoError(t, err)
		assert.False(t, regenerateCalled)
	})

	t.Run("Directory exists, watch=true, regenerates", func(t *testing.T) {
		regenerateCalled = false
		err := ensureWebsiteExists(existingDir, true, regenerate)
		assert.NoError(t, err)
		assert.True(t, regenerateCalled)
	})

	t.Run("Directory does not exist, regenerates", func(t *testing.T) {
		regenerateCalled = false
		err := ensureWebsiteExists(nonExistentDir, false, regenerate)
		assert.NoError(t, err)
		assert.True(t, regenerateCalled)
	})

	t.Run("Regeneration fails", func(t *testing.T) {
		regenerateWithError := func() error {
			return errors.New("regeneration failed")
		}
		err := ensureWebsiteExists(nonExistentDir, false, regenerateWithError)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "regeneration failed")
	})
}

func TestCreateRegenerationFunc(t *testing.T) {
	// This is a complex function that calls multiple other functions
	// We'll test that it returns a callable function
	t.Run("Returns a function", func(t *testing.T) {
		fn := createRegenerationFunc("data", "output", "en", "default")
		assert.NotNil(t, fn)
		// We don't call it because it would require full setup
		// The actual generation logic is tested in other tests
	})
}
