package packet

const PacketHeaderSize = 12

type PacketHeader struct {
	command  Command
	sequence uint32
	size     uint32
}
