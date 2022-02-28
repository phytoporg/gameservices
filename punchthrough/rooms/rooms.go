// Room storage for the holepunch/lobby service
package rooms

import (
	"errors"
	"fmt"
)

type Room struct {
	Id              int
	Name            string
	HostClientName  string
	HostClientId    int
	ClientIdsInRoom []int
}

type RoomStorage struct {
	RoomMap map[int]Room
	NextId  int
}

func New() *RoomStorage {
	return &RoomStorage{RoomMap: make(map[int]Room), NextId: 0}
}

func (rs *RoomStorage) CreateRoom(roomName string, clientName string, clientId int) Room {
	roomId := rs.NextId
	clientIdArray := make([]int, 1, 2)
	clientIdArray[0] = clientId
	rs.RoomMap[roomId] = Room{Name: roomName, HostClientName: clientName, HostClientId: clientId, ClientIdsInRoom: clientIdArray}
	rs.NextId++

	return rs.RoomMap[roomId]
}

func (rs *RoomStorage) GetRoom(roomId int) (Room, error) {
	if _, ok := rs.RoomMap[roomId]; !ok {
		return Room{}, errors.New(fmt.Sprintf("Could not get room %d", roomId))
	}

	return rs.RoomMap[roomId], nil
}

func (rs *RoomStorage) DeleteRoom(roomId int) (Room, error) {
	if _, ok := rs.RoomMap[roomId]; !ok {
		return Room{}, errors.New(fmt.Sprintf("Could not find room %d for deletion", roomId))
	}

	return rs.RoomMap[roomId], nil
}

func (rs *RoomStorage) AddClientToRoom(clientId int, roomId int) error {
	room, ok := rs.RoomMap[roomId]
	if !ok {
		return errors.New(fmt.Sprintf("Could not find room %d", roomId))
	}

	for id := range room.ClientIdsInRoom {
		if id == clientId {
			return errors.New(fmt.Sprintf("Client is already in room %d", clientId))
		}
	}

	room.ClientIdsInRoom = append(room.ClientIdsInRoom, clientId)
	rs.RoomMap[roomId] = room
	return nil
}

func (rs *RoomStorage) RemoveClientFromRoom(clientId int, roomId int) error {
	room, ok := rs.RoomMap[roomId]
	if !ok {
		return errors.New(fmt.Sprintf("Could not find room %d", roomId))
	}

	if room.HostClientId == clientId {
		// Dispand room instead!
		return errors.New(fmt.Sprintf("Cannot remove host client ID %d from room %d", roomId, roomId))
	}

	var foundClient = false
	for i := range room.ClientIdsInRoom {
		id := room.ClientIdsInRoom[i]
		if id == clientId {
			// Found it. Delete and bail.
			copy(room.ClientIdsInRoom[:i], room.ClientIdsInRoom[i+1:])
			room.ClientIdsInRoom = room.ClientIdsInRoom[:len(room.ClientIdsInRoom)-1]
			foundClient = true
			break
		}
	}

	if !foundClient {
		return errors.New(fmt.Sprintf("Could not find client %d in room %d", clientId, roomId))
	}

	rs.RoomMap[roomId] = room
	return nil
}

func (rs *RoomStorage) GetAllRooms() []Room {
	rooms := make([]Room, 0, len(rs.RoomMap))
	for _, value := range rs.RoomMap {
		rooms = append(rooms, value)
	}

	return rooms
}
