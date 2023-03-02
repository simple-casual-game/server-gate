package packet

import (
	"bytes"
	"encoding/binary"
)

const PackageCommand Command = 0x00020000

func init() {
	mainPacker.Register(PackageCommand, func(sequence uint32, body []byte) Packet {
		return packagePackFunc(sequence, body)
	})
}

type PackagePacket struct {
	header PacketHeader
	body   []byte
}

func NewPackagePacket(sequence uint32, body []byte) *PackagePacket {
	packet := &PackagePacket{
		header: PacketHeader{
			command:  ClientEnterReqCommand,
			size:     uint32(len(body)) + PacketHeaderSize,
			sequence: sequence,
		},
		body: body,
	}

	return packet
}

func (p *PackagePacket) GetCommand() Command {
	return p.header.command
}

func (p *PackagePacket) GetSize() uint32 {
	return p.header.size
}

func (p *PackagePacket) GetBodySize() uint32 {
	return uint32(len(p.body))
}

func (p *PackagePacket) GetSequence() uint32 {
	return p.header.sequence
}

func (p *PackagePacket) GetBody() []byte {
	return p.body
}

func (p *PackagePacket) GetSubpacket() (Packet, error) {
	return GetPacker().Unpack(p.body)
}

func (p *PackagePacket) ToByte() []byte {

	byteSlice := make([]byte, p.header.size)

	buffer := bytes.NewBuffer(byteSlice)
	binary.Write(buffer, binary.BigEndian, p.header.command)
	binary.Write(buffer, binary.BigEndian, p.header.size)
	binary.Write(buffer, binary.BigEndian, p.header.sequence)
	buffer.Write(p.body)

	return buffer.Bytes()
}

func packagePackFunc(sequence uint32, body []byte) *PackagePacket {
	return &PackagePacket{
		header: PacketHeader{
			command:  PackageCommand,
			size:     PacketHeaderSize + uint32(len(body)),
			sequence: sequence,
		},
		body: body,
	}
}
