// Client storage for the holepunch service
package clients

import (
	"errors"
	"fmt"
)

type ClientStorage struct {
	// Maybe some day this can be a *database* !!
	clientMap map[int]Client
	nextId    int
}

type Client struct {
	id   int
	name string
}

func New() *ClientStorage {
	return &ClientStorage{clientMap: make(map[int]Client), nextId: 0}
}

func (cs *ClientStorage) CreateClient(clientName string) int {
	// TODO: Should this fail if there's a name collision?
	clientId := cs.nextId
	cs.nextId++

	cs.clientMap[clientId] = Client{id: clientId, name: clientName}
	return clientId
}

func (cs *ClientStorage) GetClient(clientId int) (Client, error) {
	if _, ok := cs.clientMap[clientId]; !ok {
		return Client{}, errors.New(fmt.Sprintf("Invalid client ID: %d", clientId))
	}

	return cs.clientMap[clientId], nil
}

func (cs *ClientStorage) DeleteClient(clientId int) error {
	if _, ok := cs.clientMap[clientId]; !ok {
		return errors.New(fmt.Sprintf("Invalid client ID: %d", clientId))
	}

	delete(cs.clientMap, clientId)
	return nil
}

func (cs *ClientStorage) DeleteAllClients() int {
	numDeleted := len(cs.clientMap)
	for key := range cs.clientMap {
		delete(cs.clientMap, key)
	}

	return numDeleted
}

func (cs *ClientStorage) GetAllClients() []Client {
	clients := make([]Client, 0, len(cs.clientMap))
	for _, value := range cs.clientMap {
		clients = append(clients, value)
	}

	return clients
}

func (cs *ClientStorage) GetClientsByName(clientName string) []Client {
	clients := make([]Client, 0, len(cs.clientMap))
	for _, value := range cs.clientMap {
		if value.name == clientName {
			clients = append(clients, value)
		}
	}

	return clients
}
