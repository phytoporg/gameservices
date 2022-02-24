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

func TestCreateSession(t *testing.T) {
	storage := New()
	sessionId := storage.CreateSession()
	if sessionId != 0 {
		t.Fatal("Newly allocated client session has non-zero ID")
	}
}

func TestGetSession(t *testing.T) {
	storage := New()
	sessionId := storage.CreateSession()

	session, err := storage.GetSession(sessionId)
	if err != nil {
		t.Fatal("Could not retrieve newly created session by ID")
	}

	if session.id != sessionId {
		t.Fatal("Could not retrieve newly created session by ID")
	}
}

// TODO: More tests
