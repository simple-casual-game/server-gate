package packet

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const ClientEnterReqCommand Command = 0x00010000

func init() {
	mainPacker.Register(ClientEnterReqCommand, func(sequence uint32, body []byte) Packet {
		return clientEnterPackFunc(sequence, body)
	})
}

type ClientEnterReqPacket struct {
	header PacketHeader
	ip     uint32
}

func NewClientEnterReqPacket(sequence uint32, ip uint32) *ClientEnterReqPacket {
	packet := &ClientEnterReqPacket{
		header: PacketHeader{
			command:  ClientEnterReqCommand,
			size:     uint32(unsafe.Sizeof(ip)) + PacketHeaderSize,
			sequence: sequence,
		},
		ip: ip,
	}

	return packet
}

func (p *ClientEnterReqPacket) GetCommand() Command {
	return p.header.command
}

func (p *ClientEnterReqPacket) GetSize() uint32 {
	return p.header.size
}

func (p *ClientEnterReqPacket) GetBodySize() uint32 {
	return p.header.size - PacketHeaderSize
}

func (p *ClientEnterReqPacket) GetSequence() uint32 {
	return p.header.sequence
}

func (p *ClientEnterReqPacket) ToByte() []byte {

	byteSlice := make([]byte, p.header.size)

	buffer := bytes.NewBuffer(byteSlice)
	binary.Write(buffer, binary.BigEndian, p.header.command)
	binary.Write(buffer, binary.BigEndian, p.header.size)
	binary.Write(buffer, binary.BigEndian, p.header.sequence)
	binary.Write(buffer, binary.BigEndian, p.ip)

	return buffer.Bytes()
}

func (p *ClientEnterReqPacket) GetIP() uint32 {
	return p.ip
}

func clientEnterPackFunc(sequence uint32, body []byte) *ClientEnterReqPacket {

	dataBuffer := bytes.NewBuffer(body)

	var ip uint32

	binary.Read(dataBuffer, binary.BigEndian, &ip)

	return &ClientEnterReqPacket{
		header: PacketHeader{
			command:  ClientEnterReqCommand,
			size:     PacketHeaderSize + uint32(len(body)),
			sequence: sequence,
		},
		ip: ip,
	}
}
