package connection

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/simple-casual-game/server-gate/logger"
	"github.com/simple-casual-game/server-gate/packet"
)

const (
	//MaxStackSize buffer緩衝大小
	MaxStackSize = 1024 * 64
)

func NewConnection(connection *net.Conn) *Connection {
	c := &Connection{
		Connection:      connection,
		VirtualSessions: make(map[uint32]*VirtualSession),
	}

	c.BufferReader = bufio.NewReaderSize(*connection, MaxStackSize)
	c.BufferWriter = bufio.NewWriterSize(*connection, MaxStackSize)

	return c
}

type Connection struct {
	Connection      *net.Conn
	BufferReader    *bufio.Reader
	BufferWriter    *bufio.Writer
	Lock            sync.RWMutex
	VirtualSessions map[uint32]*VirtualSession
}

func (c *Connection) GetVirtualSession(sequence uint32) *VirtualSession {
	if vs, ok := c.VirtualSessions[sequence]; ok {
		return vs
	}
	return nil
}

func (c *Connection) OnData(data []byte, length int) error {
	logger.Infof("[Connection.OnData] get data len: %d", length)

	packet, err := packet.GetPacker().Unpack(data)
	if err != nil {
		logger.Warningf("[Connection.OnData] failed to unpack due to %v", err)
	}

	return c.OnPacket(packet)
}

func (c *Connection) OnPacket(p packet.Packet) error {

	switch p.GetCommand() {
	case packet.PackageCommand:
		logger.Infof("[Connection.OnPacket] get PackageCommand")

		packagePacket := p.(*packet.PackagePacket)
		subPacket, err := packagePacket.GetSubpacket()
		if err != nil {
			logger.Errorf("[Connection.OnPacket] failed to get sub packet")
			return err
		}

		if p.GetSequence() == 0 {
			return c.OnPacket(subPacket)
		}

		if vs, ok := c.VirtualSessions[p.GetSequence()]; ok {
			vs.OnPacket(subPacket)
			return nil
		}

		logger.Errorf("[Connection.OnPacket] virtual session %s not exist", p.GetSequence())
		return errors.New("no such virtual session")

	case packet.ClientEnterReqCommand:
		logger.Infof("[Connection.OnPacket] get ClientEnterReqCommand")

		clientEnterReqPacket := p.(*packet.ClientEnterReqPacket)

		if _, ok := c.VirtualSessions[p.GetSequence()]; ok {
			logger.Infof("[Connection.OnPacket] virtual session [%d] already exist", p.GetSequence())
			return nil
		}

		vsession := NewVirtualSession(c, p.GetSequence(), clientEnterReqPacket.GetIP())
		c.Lock.Lock()
		defer c.Lock.Unlock()
		c.VirtualSessions[p.GetSequence()] = vsession
		logger.Infof("[Connection.OnPacket] virtual session added to [%d]", p.GetSequence())
		return nil
	}
	return nil
}

func (c *Connection) SendPackage(seq uint32, body []byte) error {

	packer := packet.GetPacker()

	packagePacket, err := packer.Pack(packet.PackageCommand, seq, body)
	if err != nil {
		logger.Errorf("[Connection.SendPackage] fail to package seq[%d]: %v", seq, err)
		return err
	}

	_, err = c.BufferWriter.Write(packagePacket.ToByte())
	if err != nil {
		logger.Errorf("[Connection.SendPackage] fail to write buffer seq[%d]: %v", seq, err)
		return err
	}

	err = c.BufferWriter.Flush()
	if err != nil {
		logger.Errorf("[Connection.SendPackage] fail to flush seq[%d]: %v", seq, err)
		return err
	}
	return nil
}

func (c *Connection) Start(ctx context.Context) {

	fmt.Printf("Connection開始連線讀取\n")

	message := make([]byte, MaxStackSize)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Session停止連線讀取\n")
			return
		default:
		}

		//stcpsession.conn.SetReadDeadline(time.Now().Add(config.SocketTimeout * time.Second))
		//message := pbytes.GetLen(MaxStackSize)
		length, err := c.BufferReader.Read(message) //read 是 blocking  的操作，所以可以用在for loop中
		//fmt.Println("CTCPSession puller leng:", length)
		if err != nil {
			fmt.Printf("連線錯誤 %v\n", err)
			return
		}
		if length > 0 {
			c.OnData(message[:length], length)
		}
	}
}
