package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type RconClient struct {
	connection net.Conn
}

func NewRconClient(address string, password string) (*RconClient, error) {
	conn, err := net.DialTimeout("tcp", address, 20*time.Second)
	if err != nil {
		return nil, err
	}

	client := &RconClient{connection: conn}

	err = client.AuthenticateToRcon(CUSTOM_ID, SERVERDATA_AUTH, password)
	if err != nil {
		return nil, err
	}

	return &RconClient{connection: conn}, nil
}

func (c *RconClient) Close() error {
	if c.connection != nil {
		err := c.connection.Close()
		c.connection = nil
		return err
	}
	return nil
}

func (c *RconClient) AuthenticateToRcon(packetId int32, packetType PacketType, body string) error {
	_, err := c.WritePacket(packetId, packetType, body)
	if err != nil {
		return fmt.Errorf("error authenticating: %v", err)
	}

	responseHeader, err := c.ReadHeader()
	if err != nil {
		return fmt.Errorf("error reading response header: %v", err)
	}
	size := responseHeader.Size - PacketHeader
	buffer := make([]byte, size)
	if _, err := c.connection.Read(buffer); err != nil {
		return fmt.Errorf("error authenticating - body: %v", err)
	}

	if responseHeader.Id != CUSTOM_ID {
		return fmt.Errorf("error authenticating - response header id did not match custom id sent: expected %d -> got %d", CUSTOM_ID, responseHeader.Id)
	}

	if responseHeader.Type != SERVERDATA_AUTH_RESPONSE {
		return fmt.Errorf("error authenticating - response header was not expected success code: expected %d -> got %d", SERVERDATA_AUTH_RESPONSE, responseHeader.Type)
	}

	return nil
}

func (c *RconClient) ExecuteCommand(command string) (string, error) {
	c.WritePacket(CUSTOM_ID, SERVERDATA_EXECCOMMAND, command)
	responsePacket, err := c.ReadPacket()
	if err != nil {
		return "", fmt.Errorf("error reading packet(1): %w", err)
	}

	if responsePacket.Id != CUSTOM_ID {
		return "", fmt.Errorf("error executing command - response header id did not match custom id sent: expected %d -> got %d", CUSTOM_ID, responsePacket.Id)
	}

	if responsePacket.Type != SERVERDATA_RESPONSE_VALUE {
		return "", fmt.Errorf("error executing command - response header was not expected success code: expected %d -> got %d", SERVERDATA_RESPONSE_VALUE, responsePacket.Type)
	}
	return responsePacket.GetBody(), nil
}

func (c *RconClient) ReadPacket() (Packet, error) {
	r := c.connection
	var incomingPacket Packet
	var size int32 // size of the incoming packet, minus the header.

	err := binary.Read(r, binary.LittleEndian, &incomingPacket.Size)
	if err != nil {
		return incomingPacket, err
	}

	if incomingPacket.Size < 10 {
		return incomingPacket, fmt.Errorf("packet size was below minimum (had: %d, needs at least 10)", &incomingPacket.Size)
	}

	err = binary.Read(r, binary.LittleEndian, &incomingPacket.Id)
	if err != nil {
		return incomingPacket, err
	}
	size += 4

	err = binary.Read(r, binary.LittleEndian, &incomingPacket.Type)
	if err != nil {
		return incomingPacket, err
	}
	size += 4

	calculatedBodySize := incomingPacket.Size - size
	buffer := make([]byte, calculatedBodySize)
	err = binary.Read(r, binary.LittleEndian, &buffer)
	if err != nil {
		return incomingPacket, err
	}
	incomingPacket.Body = buffer
	size += calculatedBodySize

	if incomingPacket.Size != size {
		return incomingPacket, fmt.Errorf("calculated size from reading did not match packet size given: (had: %d, expected: %d)", size, incomingPacket.Size)
	}

	return incomingPacket, nil
}

func (c *RconClient) ReadHeader() (Packet, error) {
	r := c.connection
	var packet Packet
	if err := binary.Read(r, binary.LittleEndian, &packet.Size); err != nil {
		return packet, fmt.Errorf("rcon: read packet size: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &packet.Id); err != nil {
		return packet, fmt.Errorf("rcon: read packet id: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &packet.Type); err != nil {
		return packet, fmt.Errorf("rcon: read packet type: %w", err)
	}

	return packet, nil
}

func (c *RconClient) WritePacket(packetId int32, packetType PacketType, body string) (int64, error) {
	packet := NewPacket(packetId, packetType, body)
	w := c.connection
	buffer := bytes.NewBuffer(make([]byte, 0, packet.Size+4))

	_ = binary.Write(buffer, binary.LittleEndian, packet.Size)
	_ = binary.Write(buffer, binary.LittleEndian, packet.Id)
	_ = binary.Write(buffer, binary.LittleEndian, packet.Type)

	buffer.Write(append(packet.Body, 0x00, 0x00))

	return buffer.WriteTo(w)
}
