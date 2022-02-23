// TODO: write even just a little bit of testing, I have no idea how to golang!!!!
package sessions

import (
	"errors"
	"fmt"

	"github.com/phytoporg/gameservices/punchthrough/clients"
)

const MaxClientsPerSession = 2

type SessionStorage struct {
	// No databases here either.
	sessionMap map[int]Session
	nextId     int
}

type Session struct {
	id           int
	clientsArray []clients.Client
}

// Session storage for the holepunch service
func New() *SessionStorage {
	return &SessionStorage{sessionMap: make(map[int]Session), nextId: 0}
}

// Returns session ID
func (ss *SessionStorage) CreateSession() int {
	ss.sessionMap[ss.nextId] = Session{id: ss.nextId, clientsArray: make([]clients.Client, 0, MaxClientsPerSession)}
	sessionId := ss.nextId
	ss.nextId++

	return sessionId
}

func (ss *SessionStorage) GetSession(sessionId int) (*Session, error) {
	session, ok := ss.sessionMap[sessionId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	return &session, nil
}

func (ss *SessionStorage) AddClientToSession(sessionId int, clientToAdd clients.Client) error {
	session, ok := ss.sessionMap[sessionId]
	if !ok {
		return errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	clientId := clientToAdd.Id

	// Make sure we're not adding duplicates
	for i := range session.clientsArray {
		client := session.clientsArray[i]
		if client.Id == clientId {
			return errors.New(fmt.Sprintf("Client %d is already in session %d", clientId, session.id))
		}
	}

	// Don't exceed the max number of clients allowed per session
	if len(session.clientsArray) == MaxClientsPerSession {
		return errors.New(fmt.Sprintf("Cannot add more clients to session %d of size %d", session.id, MaxClientsPerSession))
	}

	session.clientsArray = append(session.clientsArray, clientToAdd)
	return nil
}

func (ss *SessionStorage) RemoveClientFromSession(sessionId int, clientId int) error {
	session, ok := ss.sessionMap[sessionId]
	if !ok {
		return errors.New(fmt.Sprintf("Invalid session ID: %d", sessionId))
	}

	// Does this client exist in the session?
	for i := range session.clientsArray {
		client := session.clientsArray[i]
		if client.Id == clientId {
			// Found it. Delete and get out.
			copy(session.clientsArray[:i], session.clientsArray[i+1:])
			session.clientsArray = session.clientsArray[:len(session.clientsArray)-1]
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Could not find client %d in session %d", clientId, sessionId))
}

func (ss *SessionStorage) GetSessionContainingClient(clientId int) (Session, error) {
	for key, value := range ss.sessionMap {
		if value.id == clientId {
			return ss.sessionMap[key], nil
		}
	}

	return Session{id: 0, clientsArray: make([]clients.Client, 0, MaxClientsPerSession)}, errors.New(fmt.Sprintf("Could find no session containing client %d", clientId))

}
