package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

////////////////////////////////////////////////////////////////////////////////
// ROOMS
// TODO: put this in a package?
type Room struct {
	Name           string
	HostClientName string
}

type RoomStorage struct {
	RoomMap map[string]Room
}

func (rs *RoomStorage) CreateRoom(roomName string) error {
	if room, ok := rs.RoomMap[roomName]; ok {
		return errors.New(fmt.Sprintf("Room %s already exists", roomName))
	}

	rs.RoomMap[roomName] = Room{Name: roomName, HostClientName: "PlaceholderName"}
}

// TODO: Get/Delete rooms
////////////////////////////////////////////////////////////////////////////////

// Testing out some message parsing junk
// TODO: messages to:
// - Register a client (hello)
// - Disconnect
// - Host a room
// - Delete a room
// - Join a room
// - Leave a room
// - Begin a session

const MSG_HELLO = 0x01
const MSG_ECHO = 0x02
const MSG_BYE = 0x03

const Protocol = "udp"
const Port = 9999

func serve(packetConn net.PacketConn, addr net.Addr, buf []byte) {
	msgId := buf[0]
	switch msgId {
	case MSG_HELLO:
		fmt.Println("Hello")
	case MSG_ECHO:
		fmt.Println("Echo: ", string(buf[1:]))
	case MSG_BYE:
		fmt.Println("Bye")
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
