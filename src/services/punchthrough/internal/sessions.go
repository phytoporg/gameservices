// TODO: write even just a little bit of testing, I have no idea how to golang!!!!
package sessions

import (
	"errors"
	"fmt"
)

const MaxClientsPerSession = 2

type SessionStorage struct {
	// No databases here either.
	sessionMap map[int]Session
	nextId     int
}

type Session struct {
	id      int
	clients []clients.Client
}

// Session storage for the holepunch service
func New() *SessionStorage {
	return &SessionStorage{sessionMap: make(map[int]Session), nextId: 0}
}

// Returns session ID
func (ss *SessionStorage) CreateSession() int {
	ss.sessionMap[ss.nextId] = Session{id: ss.nextId, make([]clients, 0, MaxClientsPerSession)}
	sessionId = ss.nextId
	ss.nextId++

	return sessionId
}

func (ss *SessionStorage) GetSession(sessionId int) (*Session, error) {
	if session, ok := ss.sessionMap[sessionId]; !ok {
		return nil, errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	return &session, nil
}

func (ss *SessionStorage) AddClientToSession(sessionId int, clientToAdd Client) error {
	if session, ok := ss.sessionMap[sessionId]; !ok {
		return errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	clientId := clientToAdd.id

	// Make sure we're not adding duplicates
	for client := range ss.clients {
		if client.id == clientId {
			return errors.New(fmt.Sprintf("Client %d is already in session %d"))
		}
	}

	// Don't exceed the max number of clients allowed per session
	if len(ss.clients) == MaxClientsPerSession {
		return errors.New(fmt.Sprintf("Cannot add more clients to session %d of size %d"), session.id, MaxClientsPerSession)
	}

	append(session.clients, clientToAdd)
	return nil
}

func (ss *SessionStorage) RemoveClientFromSession(sessionId int, clientId int) error {
	if session, ok := &ss.sessionMap[sessionId]; !ok {
		return errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	// Does this client exist in the session?
	for i := 0; i < len(ss.clients); i++ {
		client := session.clients[i]
		if client.id == clientId {
			// Found it. Delete and get out.
			copy(session.clients[:i], session.clients[i+1:])
			session.clients = session.clients[:len(session.clients)-1]
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Could not find client %d in session %d"), clientId, sessionId)
}

func (ss *SessionStorage) GetSessionContainingClient(clientId int) (*Session, error) {
	for key, value := range ss.sessionMap {
		if value.id == clientId {
			return ss.sessionMap[key], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Could find no session containing client %d", clientId))
}
