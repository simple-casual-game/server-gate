package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/simple-casual-game/server-gate/logger"
)

func GetPacker() *Packer {
	return &mainPacker
}

type Packer struct {
	lock      sync.RWMutex
	packFuncs map[Command]func(sequence uint32, body []byte) Packet
}

var mainPacker Packer = Packer{
	packFuncs: make(map[Command]func(sequence uint32, body []byte) Packet),
}

func (p *Packer) Register(command Command, packFunc func(sequence uint32, body []byte) Packet) error {

	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.packFuncs[command]; ok {
		logger.Warningf("[Packer.Register] command %d already exist", command)
	}

	p.packFuncs[command] = packFunc
	return nil
}

func (p *Packer) Pack(command Command, sequence uint32, body []byte) (Packet, error) {

	p.lock.RLock()
	defer p.lock.RUnlock()

	if packFunc, ok := p.packFuncs[command]; !ok {
		logger.Warningf("[Packer.Pack] command %d has no packer", command)
		return nil, errors.New("no packer")
	} else {
		return packFunc(sequence, body), nil
	}
}

func (p *Packer) Unpack(packetBytes []byte) (Packet, error) {

	dataBuffer := bytes.NewBuffer(packetBytes)

	// 讀標頭
	var command Command
	var size uint32
	var sequence uint32
	var body []byte

	binary.Read(dataBuffer, binary.BigEndian, &command)
	binary.Read(dataBuffer, binary.BigEndian, &size)
	binary.Read(dataBuffer, binary.BigEndian, &sequence)

	body = dataBuffer.Next(int(size) - PacketHeaderSize)

	return p.Pack(command, sequence, body)
}
