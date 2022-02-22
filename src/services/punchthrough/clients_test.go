package clients

import (
	"fmt"
	"testing"
)

func TestNewClientStorage(t *testing.T) {
	storage := New()
	if storage == nil {
		t.Fatal("Failed to create client storage")
	}

	if len(storage.clientMap) > 0 {
		t.Fatal("Newly allocated client storage is not empty")
	}

	if storage.nextId > 0 {
		t.Fatal("Newly allocated client storage has non-zero next ID")
	}
}

func TestCreateClient(t *testing.T) {
	storage := New()
	firstClientId := storage.CreateClient("test_client_1")
	if firstClientId > 0 {
		t.Fatal("First client has non-zero ID")
	}
}

func TestCreatesUniqueClients(t *testing.T) {
	const NumToCreate = 10

	storage := New()

	// No sets in golang? We'll just n^2 it over an array. This thing's small anyway.
	clientIds := make([]int, 0, NumToCreate)
	for i := 0; i < NumToCreate; i++ {
		clientId := storage.CreateClient(fmt.Sprintf("test_client_%d", i))
		clientIds = append(clientIds, clientId)
	}

	for i := 0; i < len(clientIds)-1; i++ {
		for j := i + 1; j < len(clientIds); j++ {
			if clientIds[i] == clientIds[j] {
				t.Fatalf("Created two clients with identical IDs: %d == %d", clientIds[i], clientIds[j])
			}
		}
	}
}

func TestCreateAndGetSingleClient(t *testing.T) {
	storage := New()
	clientId := storage.CreateClient("test_client")

	client, err := storage.GetClient(clientId)
	if err != nil {
		t.Fatal("Failed to get newly created client")
	}

	if client.id != clientId {
		t.Fatal("Retrieved client with mismatched client ID")
	}
}

func TestCreateAndDeleteSingleClient(t *testing.T) {
	storage := New()
	clientId := storage.CreateClient("test_client")

	_, err := storage.GetClient(clientId)
	if err != nil {
		t.Fatal("Failed to get newly created client.")
	}

	err = storage.DeleteClient(clientId)
	if err != nil {
		t.Fatal("Failed to delete newly added client.")
	}
}

func TestCreateAndDeleteSingleNonExistingClient(t *testing.T) {
	storage := New()
	clientId := storage.CreateClient("test_client")

	_, err := storage.GetClient(clientId)
	if err != nil {
		t.Fatal("Failed to get newly created client.")
	}

	err = storage.DeleteClient(clientId + 1)
	if err == nil {
		t.Fatal("Deleting non-existent client produced no error.")
	}
}

func TestGetAllClients(t *testing.T) {
	const NumToCreate = 10
	storage := New()

	for i := 0; i < NumToCreate; i++ {
		storage.CreateClient(fmt.Sprintf("test_client_%d", i))
	}

	clients := storage.GetAllClients()
	if len(clients) != NumToCreate {
		t.Fatalf("Unexpected number of clients: %d != %d", len(clients), NumToCreate)
	}

	for i := 0; i < NumToCreate; i++ {
		expectedClientName := fmt.Sprintf("test_client_%d", clients[i].id)
		if clients[i].name != expectedClientName {
			t.Fatalf("Unexpected client: %s != %s", clients[i].name, expectedClientName)
		}
	}
}

func TestDeleteAllClients(t *testing.T) {
	const NumToCreate = 10
	storage := New()

	for i := 0; i < NumToCreate; i++ {
		storage.CreateClient(fmt.Sprintf("test_client_%d", i))
	}

	clients := storage.GetAllClients()
	if len(clients) != NumToCreate {
		t.Fatalf("Unexpected number of clients: %d != %d", len(clients), NumToCreate)
	}

	numDeleted := storage.DeleteAllClients()
	if numDeleted != NumToCreate {
		t.Fatalf("Deleted unexpected number of clients: %d != %d", numDeleted, NumToCreate)
	}

	clients = storage.GetAllClients()
	if len(clients) > 0 {
		t.Fatalf("Non-zero clients remaining after deleting all: %d", len(clients))
	}
}

func TestGetClientsByName(t *testing.T) {
	// TODO
}
