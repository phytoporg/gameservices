// Message parsing and handling helper package
package messages

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type MessageId uint8

const (
	ClientHello     MessageId = 0x01
	ClientGoodbye             = 0x02
	ClientHostRoom            = 0x03
	ClientJoinRoom            = 0x04
	ClientLeaveRoom           = 0x05
	MaxMessageId
)

type MessageHeader struct {
	Id MessageId
}

type Message interface {
	Parse([]byte) error
}

type MessageHandler func(m Message, packetConn net.PacketConn, addr net.Addr)

// Hello
type MessageClientHello struct {
	Header     MessageHeader
	ClientName string
}

func (m *MessageClientHello) Parse(buf []byte) error {
	m.ClientName = string(buf[1:])
	return nil
}

// Goodbye
type MessageClientGoodbye struct {
	Header   MessageHeader
	ClientId int
}

func (m *MessageClientGoodbye) Parse(buf []byte) error {
	m.ClientId = int(binary.LittleEndian.Uint32(buf[:4]))
	return nil
}

// Host room
type MessageClientHostRoom struct {
	Header   MessageHeader
	ClientId int
	RoomName string
}

func (m *MessageClientHostRoom) Parse(buf []byte) error {
	m.ClientId = int(binary.LittleEndian.Uint32(buf[:4]))
	m.RoomName = string(buf[4:])
	return nil
}

// Join room
type MessageClientJoinRoom struct {
	Header   MessageHeader
	ClientId int
	RoomId   int
}

func (m *MessageClientJoinRoom) Parse(buf []byte) error {
	m.ClientId = int(binary.LittleEndian.Uint32(buf[:4]))
	m.RoomId = int(binary.LittleEndian.Uint32(buf[4:8]))
	return nil
}

// Leave room
type MessageClientLeaveRoom struct {
	// TODO
}

// Message dispatcher
type MessageDispatcher struct {
	handlerRegistry map[MessageId]MessageHandler
}

func NewDispatcher() *MessageDispatcher {
	return &MessageDispatcher{handlerRegistry: make(map[MessageId]MessageHandler)}
}

func (md *MessageDispatcher) RegisterHandler(id MessageId, handler MessageHandler) error {
	if id <= 0 || id >= MaxMessageId {
		return errors.New(fmt.Sprintf("Invalid message ID: %d", id))
	}

	md.handlerRegistry[id] = handler
	return nil
}

func (md *MessageDispatcher) ParseAndDispatchMessage(buf []byte, packetConn net.PacketConn, addr net.Addr) error {
	if len(buf) < 1 {
		return errors.New("Buffer is empty")
	}

	messageId := MessageId(buf[0])
	if _, ok := md.handlerRegistry[messageId]; !ok {
		return errors.New(fmt.Sprintf("No handler for message ID %d", messageId))
	}

	header := MessageHeader{Id: messageId}
	var message Message

	switch messageId {
	case ClientHello:
		message = &MessageClientHello{Header: header}
	case ClientGoodbye:
		message = &MessageClientGoodbye{Header: header}
	case ClientHostRoom:
		message = &MessageClientHostRoom{Header: header}
	case ClientJoinRoom:
		message = &MessageClientJoinRoom{Header: header}
	case ClientLeaveRoom:
		return errors.New("Not implemented: MessageClientLeaveRoom")
	}

	err := message.Parse(buf[1:])
	if err != nil {
		return err
	}

	handler := md.handlerRegistry[messageId]
	handler(message, packetConn, addr)

	return nil
}
