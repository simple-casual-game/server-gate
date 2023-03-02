package packet

type Command uint32

type Packet interface {
	GetCommand() Command
	GetSize() uint32
	GetBodySize() uint32
	GetSequence() uint32
	ToByte() []byte
}
