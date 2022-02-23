package sessions

import (
	"testing"
)

func TestNewSessionStorage(t *testing.T) {
	storage := New()
	if storage == nil {
		t.Fatal("Failed to create session storage")
	}

	if len(storage.sessionMap) > 0 {
		t.Fatal("Newly allocated session storage is not empty")
	}

	if storage.nextId > 0 {
		t.Fatal("Newly allocated client storage has non-zero next ID")
	}
}

// TODO: All of the other test
