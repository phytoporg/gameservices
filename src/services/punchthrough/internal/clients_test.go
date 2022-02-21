package clients

import (
	"fmt"
	"testing"
)

func TestNewClientStorage(t *testing.T) {
	fmt.Println("TESTING")
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
