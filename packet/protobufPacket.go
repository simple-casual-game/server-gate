package packet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/simple-casual-game/server-gate/logger"
	"github.com/simple-casual-game/server-gate/protobuf/game"
	"google.golang.org/protobuf/proto"
)

const ProtobufCommand Command = 0x00030000

func init() {
	mainPacker.Register(ProtobufCommand, func(sequence uint32, body []byte) Packet {
		return protobufPackFunc(sequence, body)
	})
}

type ProtobufPacket struct {
	header PacketHeader
	data   []byte
}

func NewProtobufPacket(sequence uint32, message *game.GameMessage) *ProtobufPacket {
	packet := &ProtobufPacket{
		header: PacketHeader{
			command:  ProtobufCommand,
			sequence: sequence,
		},
	}

	data, err := proto.Marshal(message)
	if err != nil {
		logger.Errorf("[Packer.Pack] unable to marshal messaage %+v", message)
		return nil
	}

	packet.header.size = uint32(len(data)) + PacketHeaderSize
	packet.data = data

	return packet
}

func (p *ProtobufPacket) GetCommand() Command {
	return p.header.command
}

func (p *ProtobufPacket) GetSize() uint32 {
	return p.header.size
}

func (p *ProtobufPacket) GetBodySize() uint32 {
	return uint32(len(p.data))
}

func (p *ProtobufPacket) GetSequence() uint32 {
	return p.header.sequence
}

func (p *ProtobufPacket) GetData() []byte {
	return p.data
}

func (p *ProtobufPacket) GetProtobuf() (*game.GameMessage, error) {

	gameMessage := &game.GameMessage{}

	if err := proto.Unmarshal(p.data, gameMessage); err != nil {
		return nil, errors.New("[ProtobufPacket.GetProtobuf] fail to parse to protobuf")
	}

	return gameMessage, nil
}

func (p *ProtobufPacket) ToByte() []byte {

	byteSlice := make([]byte, p.header.size)

	buffer := bytes.NewBuffer(byteSlice)
	binary.Write(buffer, binary.BigEndian, p.header.command)
	binary.Write(buffer, binary.BigEndian, p.header.size)
	binary.Write(buffer, binary.BigEndian, p.header.sequence)
	buffer.Write(p.data)

	return buffer.Bytes()
}

func protobufPackFunc(sequence uint32, body []byte) *ProtobufPacket {
	return &ProtobufPacket{
		header: PacketHeader{
			command:  PackageCommand,
			size:     PacketHeaderSize + uint32(len(body)),
			sequence: sequence,
		},
		data: body,
	}
}
