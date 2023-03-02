//go:build test
// +build test

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/simple-casual-game/server-gate/protobuf"
	"github.com/simple-casual-game/server-gate/protobuf/diamonds"
	"google.golang.org/protobuf/proto"
)

/*

go run main2.go connection.go command.go const.go packet.go VirtualSession.go client.go state.go bet.go connectionManager.go

*/

func main() {

	addr := "127.0.0.1:8001"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial error %v\n", err)
		return
	}
	defer conn.Close()
	bufferWriter := bufio.NewWriterSize(conn, MaxStackSize)

	go func() {
		for {
			bs := make([]byte, 1024)
			len, err := conn.Read(bs)
			if err != nil {
				fmt.Printf("read error %v\n", err)
			} else {
				fmt.Printf("read: %s\n", bs[:len])
			}
		}
	}()

	//enter

	msg := PkClientEnter{}

	msg.Head.Cmd = CmdClientEnter
	msg.Head.Size = ClientEnterSize
	msg.Head.Seq = 123

	tmpData := make([]byte, 0, ClientEnterSize)
	buffer := bytes.NewBuffer(tmpData)
	binary.Write(buffer, binary.BigEndian, msg)

	if err = sendEnter(bufferWriter, buffer); err != nil {
		fmt.Printf("sendEnter error %v\n", err)
		return
	}

	//login

	req := &protobuf.Header{
		Payload: &protobuf.Header_GateLoginRequest{
			GateLoginRequest: &protobuf.GateLoginRequest{
				Name: "testtest",
			},
		},
	}

	dataBuffer, err := proto.Marshal(req)
	if err != nil {
		fmt.Printf("proto marshal error %v\n", err)
		return
	}

	bufferData := ProtoToPackets(dataBuffer, len(dataBuffer))

	if err = sendPackage(bufferWriter, 123, bufferData, len(bufferData)); err != nil {
		fmt.Printf("sendPackage error %v\n", err)
		return
	}

	//bet

	betReq := &protobuf.Header{
		Payload: &protobuf.Header_DiamondsCHeader{
			&diamonds.CHeader{
				Payload: &diamonds.CHeader_DiamondsTsBet{
					&diamonds.DIAMONDS_TS_BET{
						Betamount: "10",
					},
				},
			},
		},
	}

	dataBuffer, err = proto.Marshal(betReq)
	if err != nil {
		fmt.Printf("proto marshal error %v\n", err)
		return
	}

	bufferData = ProtoToPackets(dataBuffer, len(dataBuffer))

	if err = sendPackage(bufferWriter, 123, bufferData, len(bufferData)); err != nil {
		fmt.Printf("sendPackage error %v\n", err)
		return
	}

	for {
	}

}

func sendTCP(addr, msg string) (string, error) {
	// connect to this socket
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// send to socket
	conn.Write([]byte(msg))

	// listen for reply
	bs := make([]byte, 1024)
	len, err := conn.Read(bs)
	if err != nil {
		return "", err
	} else {
		return string(bs[:len]), err
	}

}

func sendEnter(bufferWriter *bufio.Writer, buffer *bytes.Buffer) error {

	_, err := bufferWriter.Write(buffer.Bytes())

	if err != nil {
		return err
	}
	err = bufferWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

func sendPackage(bufferWriter *bufio.Writer, seq int, body []byte, size int) error {
	pkPackage := PkPackage{}

	pkPackage.Head.Cmd = CmdPackage
	pkPackage.Head.Seq = uint32(seq)
	pkPackage.Head.Size = uint32(PackageSize + size)

	//tmpData := pbytes.GetCap(gatecommand.PackageSize + size)
	tmpData := make([]byte, 0, PackageSize+uint32(PackageSize+size))
	buffer := bytes.NewBuffer(tmpData)

	binary.Write(buffer, binary.BigEndian, pkPackage)

	buffer.Write(body[:size])

	_, err := bufferWriter.Write(buffer.Bytes())

	if err != nil {
		return err
	}
	err = bufferWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}
