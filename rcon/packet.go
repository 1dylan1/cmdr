package rcon

type PacketType int32
const (
	SERVERDATA_RESPONSE_VALUE 	PacketType = 0 
	SERVERDATA_AUTH_RESPONSE 	PacketType = 2
	SERVERDATA_EXECCOMMAND 		PacketType = 2
	SERVERDATA_AUTH 			PacketType = 3
)

const (
	CUSTOM_ID int32 = 52
	EMPTY_PACKET_ID int32 = 1337
	PacketPadding int32 = 2
	PacketHeader int32 = 8

	MinPacket = PacketPadding + PacketHeader
	MaxPacket = 4096 + MinPacket
)

type Packet struct {
	Size int32
	Id int32
	Type PacketType
	Body []byte
}

func NewPacket(packetId int32, packetType PacketType, body string) *Packet {
	size := len([]byte(body)) + int(PacketHeader+PacketPadding)

	return &Packet{
		Size: int32(size),
		Type: packetType,
		Id: packetId,
		Body: []byte(body),
	}
}

func (p *Packet) GetBody() string{
	return string(p.Body)
}
