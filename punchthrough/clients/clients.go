// Client storage for the holepunch service
package clients

import (
	"errors"
	"fmt"
	"net"
)

type ClientStorage struct {
	// Maybe some day this can be a *database* !!
	ClientMap map[int]Client
	NextId    int
}

type Client struct {
	Id      int
	Name    string
	Address net.Addr
}

func New() *ClientStorage {
	return &ClientStorage{ClientMap: make(map[int]Client), NextId: 0}
}

func (cs *ClientStorage) CreateClient(clientName string, addr net.Addr) int {
	// TODO: Should this fail if there's a name collision?
	clientId := cs.NextId
	cs.NextId++

	cs.ClientMap[clientId] = Client{Id: clientId, Name: clientName, Address: addr}
	return clientId
}

func (cs *ClientStorage) GetClient(clientId int) (Client, error) {
	if _, ok := cs.ClientMap[clientId]; !ok {
		return Client{}, errors.New(fmt.Sprintf("Invalid client ID: %d", clientId))
	}

	return cs.ClientMap[clientId], nil
}

func (cs *ClientStorage) DeleteClient(clientId int) error {
	if _, ok := cs.ClientMap[clientId]; !ok {
		return errors.New(fmt.Sprintf("Invalid client ID: %d", clientId))
	}

	delete(cs.ClientMap, clientId)
	return nil
}

func (cs *ClientStorage) DeleteAllClients() int {
	numDeleted := len(cs.ClientMap)
	for key := range cs.ClientMap {
		delete(cs.ClientMap, key)
	}

	return numDeleted
}

func (cs *ClientStorage) GetAllClients() []Client {
	clients := make([]Client, 0, len(cs.ClientMap))
	for _, value := range cs.ClientMap {
		clients = append(clients, value)
	}

	return clients
}

func (cs *ClientStorage) GetClientsByName(clientName string) []Client {
	clients := make([]Client, 0, len(cs.ClientMap))
	for _, value := range cs.ClientMap {
		if value.Name == clientName {
			clients = append(clients, value)
		}
	}

	return clients
}
