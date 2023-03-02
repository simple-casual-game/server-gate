package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/simple-casual-game/server-gate/connection"
	"github.com/simple-casual-game/server-gate/global"
)

func main() {

	ip := "127.0.0.1"
	port := "8002"

	ipaddr := ip + ":" + port
	listener, err := net.Listen("tcp", ipaddr)
	if err != nil {
		fmt.Printf("fail to listen tcp %v", err)
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("fail to accept listener %v", err)
		return
	}

	webConnection := connection.NewConnection(&conn)
	connection.GetConnectionManager().AddConnection(global.WEB, webConnection)

	ctx, cancel := context.WithCancel(context.Background())
	go webConnection.Start(ctx)

	ip = "127.0.0.1"
	port = "8001"
	ipaddr = ip + ":" + port

	d := net.Dialer{Timeout: time.Second * 5}
	conn, err = d.Dial("tcp", ipaddr)
	if err != nil {
		fmt.Printf("fail to dial game server %v", err)
		return
	}

	flipCoinConnection := connection.NewConnection(&conn)
	connection.GetConnectionManager().AddConnection(global.GAME_FLIPCOIN, flipCoinConnection)
	go flipCoinConnection.Start(ctx)

	cancelChan := make(chan int, 0)
	for {
		select {
		case <-cancelChan:
			cancel()
			return
		}
	}
}
