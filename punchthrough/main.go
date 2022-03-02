package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/phytoporg/gameservices/punchthrough/clients"
	"github.com/phytoporg/gameservices/punchthrough/messages"
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
var messageDispatcher = messages.NewDispatcher()

const Protocol = "udp"
const Port = 9999

func serve(packetConn net.PacketConn, addr net.Addr, buf []byte) {
	fmt.Println("MsgId: ", buf[0])

	messageDispatcher.ParseAndDispatchMessage(buf, packetConn, addr)
}

func HandleClientHello(m messages.Message, packetConn net.PacketConn, addr net.Addr) {
	helloMessage := m.(*messages.MessageClientHello)

	clientId := clientStorage.CreateClient(helloMessage.ClientName, addr)
	fmt.Println("Hello from ", helloMessage.ClientName, " (", clientId, ")")

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
}

func HandleClientGoodbye(m messages.Message, packetConn net.PacketConn, addr net.Addr) {
	goodbyeMessage := m.(*messages.MessageClientGoodbye)

	client, err := clientStorage.GetClient(goodbyeMessage.ClientId)
	if err != nil {
		log.Fatal(err)
		// TODO: Send an error back to the client
		return
	}

	fmt.Println("Goodbye to client: ", client.Name, " (", client.Id, ")")
}

func HandleClientHostRoom(m messages.Message, packetConn net.PacketConn, addr net.Addr) {
	hostRoomMessage := m.(*messages.MessageClientHostRoom)

	client, err := clientStorage.GetClient(hostRoomMessage.ClientId)
	if err != nil {
		log.Fatal(err)
		// TODO: Send an error back to the client
		return
	}

	room := roomStorage.CreateRoom(hostRoomMessage.RoomName, client.Name, client.Id)

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
}

func HandleClientJoinRoom(m messages.Message, packetConn net.PacketConn, addr net.Addr) {
	joinRoomMessage := m.(*messages.MessageClientJoinRoom)

	_, err := clientStorage.GetClient(joinRoomMessage.ClientId)
	if err != nil {
		log.Fatal(err)
		// TODO: Send an error back to the client
		return
	}

	room, err := roomStorage.GetRoom(joinRoomMessage.RoomId)
	if err != nil {
		log.Fatal(err)
		// TODO: Send an error back to the client
		return
	}

	fmt.Println(
		"Add client ID ",
		joinRoomMessage.ClientId,
		" to room \"", room.Name, "\" (", room.Id, ")")
	err = roomStorage.AddClientToRoom(joinRoomMessage.ClientId, joinRoomMessage.RoomId)
	if err != nil {
		log.Fatal(err)
		// TODO: Send an error back to the client
		return
	}

	fmt.Println(
		"Client ID ",
		joinRoomMessage.ClientId,
		" joined room \"", room.Name, "\" (", room.Id, ")")
}

func main() {
	packetConn, err := net.ListenPacket(Protocol, fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatal(err)
	}
	defer packetConn.Close()

	messageDispatcher.RegisterHandler(messages.ClientHello, HandleClientHello)
	messageDispatcher.RegisterHandler(messages.ClientGoodbye, HandleClientGoodbye)
	messageDispatcher.RegisterHandler(messages.ClientHostRoom, HandleClientHostRoom)
	messageDispatcher.RegisterHandler(messages.ClientJoinRoom, HandleClientJoinRoom)

	for {
		buf := make([]byte, 1024)
		n, addr, err := packetConn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		go serve(packetConn, addr, buf[:n])
	}
}
