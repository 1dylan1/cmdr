package rcon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type RconClient struct {
	connection net.Conn
}

func NewRconClient(address string, password string) (*RconClient, error) {
	conn, err := net.Dial("tcp", address)
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
		return "" , fmt.Errorf("error reading packet(1): %w", err)
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
	var packet Packet
	var empty Packet
	var n int64

	err := binary.Read(r, binary.LittleEndian, &packet.Size)
	if err != nil {
		log.Printf("Error inside packet size")
		return packet, err
	}
	n += 4
	if packet.Size < MinPacket {
		return empty, fmt.Errorf("packet size too small")
	}
	n += 4
	err = binary.Read(r, binary.LittleEndian, &packet.Id)
	if err != nil {
		log.Printf("Error inside packet id")
		return empty, err
	}
	n += 4
	err = binary.Read(r, binary.LittleEndian, &packet.Type)
	if err != nil {
		log.Printf("Error inside packet type")
		return empty, err
	}
	n += 4	
	packet.Body = make([]byte, packet.Size - PacketHeader)
	
	var i int64
	for i < int64(packet.Size - PacketHeader) {
		var m int
		var err error

		if m, err = r.Read(packet.Body[i:]); err != nil {
			return empty, fmt.Errorf("read err: %w", err)
		}

		i += int64(m)
	}
	n += 1

	if !bytes.Equal(packet.Body[len(packet.Body) - int(PacketPadding):], []byte{0x00,0x00}) {
		return empty, fmt.Errorf("invalid packet padding")
	}

	packet.Body = packet.Body[0 : len(packet.Body) - int(PacketPadding)]
	return packet, nil
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

