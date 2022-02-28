package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/phytoporg/gameservices/punchthrough/clients"
	"github.com/phytoporg/gameservices/punchthrough/rooms"
)

// Testing out some message parsing junk
// TODO: messages to:
// - Register a client (hello) *
// - Disconnect a client (goodbye) *
// - Host a room
// - Delete a room
// - Join a room
// - Leave a room
// - List rooms
// - Begin a session

var clientStorage = clients.New()
var roomStorage = rooms.New()

const MSG_CLIENT_HELLO = 0x01
const MSG_CLIENT_GOODBYE = 0x02
const MSG_CLIENT_HOST_ROOM = 0x03
const MSG_CLIENT_JOIN_ROOM = 0x04
const MSG_CLIENT_LEAVE_ROOM = 0x05

const Protocol = "udp"
const Port = 9999

func serve(packetConn net.PacketConn, addr net.Addr, buf []byte) {
	msgId := buf[0]
	fmt.Println("MsgId: ", msgId)
	switch msgId {
	case MSG_CLIENT_HELLO:
		clientName := string(buf[1:])
		clientId := clientStorage.CreateClient(clientName, addr)
		fmt.Println("Hello from ", clientName, " (", clientId, ")")

		// Send the client their assigned ID
		bytesToSend := make([]byte, 4)

		// Who needs network byte order??
		binary.LittleEndian.PutUint32(bytesToSend, uint32(clientId))
		_, err := packetConn.WriteTo(bytesToSend, addr)
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

	case MSG_CLIENT_GOODBYE:
		clientId := binary.LittleEndian.Uint32(buf[1:])
		client, err := clientStorage.GetClient(int(clientId))
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		fmt.Println("Goodbye to client: ", client.Name, " (", client.Id, ")")

	case MSG_CLIENT_HOST_ROOM:
		clientId := binary.LittleEndian.Uint32(buf[1:5])
		client, err := clientStorage.GetClient(int(clientId))
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		roomName := string(buf[5:])
		room := roomStorage.CreateRoom(roomName, client.Name, client.Id)

		// Send the client the room ID
		bytesToSend := make([]byte, 4)

		// Who needs network byte order??
		binary.LittleEndian.PutUint32(bytesToSend, uint32(room.Id))
		_, err = packetConn.WriteTo(bytesToSend, addr)
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		fmt.Println(
			"Created room \"", room.Name, "\" hosted by ", room.HostClientName,
			" (", room.HostClientId, ")")
	case MSG_CLIENT_JOIN_ROOM:
		clientId := binary.LittleEndian.Uint32(buf[1:5])
		_, err := clientStorage.GetClient(int(clientId))
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		roomId := binary.LittleEndian.Uint32(buf[5:9])
		room, err := roomStorage.GetRoom(int(roomId))
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		err = roomStorage.AddClientToRoom(int(clientId), int(roomId))
		if err != nil {
			log.Fatal(err)
			// TODO: Send an error back to the client
			return
		}

		fmt.Println("Client ID ", clientId, " joined room \"", room.Name, "\" (", room.Id, ")")
	}
}

func main() {
	packetConn, err := net.ListenPacket(Protocol, fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatal(err)
	}
	defer packetConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := packetConn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		go serve(packetConn, addr, buf[:n])
	}
}
