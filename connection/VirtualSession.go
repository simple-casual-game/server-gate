package connection

import (
	"errors"
	"fmt"

	"github.com/simple-casual-game/server-gate/dao/clientDao"
	"github.com/simple-casual-game/server-gate/logger"
	"github.com/simple-casual-game/server-gate/packet"
	"github.com/simple-casual-game/server-gate/protobuf/flipCoin"
	"github.com/simple-casual-game/server-gate/protobuf/game"
	"github.com/simple-casual-game/server-gate/task"
)

func NewVirtualSession(c *Connection, clientID uint32, ip uint32) *VirtualSession {
	return &VirtualSession{
		Connection: c,
		ClientID:   clientID,
		IP:         ip,
	}
}

type Client struct {
	IP        uint32 /**< \brief 客戶端當前 IP */
	UserName  string //不帶pid的帳號名稱
	SessionID uint32
}

type VirtualSession struct {
	Connection *Connection
	ClientID   uint32
	IP         uint32
	State      State
	Client     *Client
}

func (v *VirtualSession) OnPacket(p packet.Packet) error {
	logger.Infof("[VirtualSession.OnPacket] get packet %+v", p)
	if v.State == State_Disconnected {
		return errors.New("Disconnected")
	}

	switch p.GetCommand() {
	case packet.Command(packet.ProtobufCommand):
		if v.State != State_Connecting && v.State != State_Connected {
			logger.Errorf("[VirtualSession.OnPacket] wrong state %+v", v.State)
			return errors.New("state error")
		}

		protobufPacket := p.(*packet.ProtobufPacket)
		gameMessage, err := protobufPacket.GetProtobuf()
		if err != nil {
			logger.Errorf("[VirtualSession.OnPacket] failed to parse game message %v", err)
			return err
		}

		if v.State == State_Connecting {

			switch gameMessage.Payload.(type) {
			case *game.GameMessage_LoginReq:
				// TODO: add client
				header := gameMessage.GetLoginReq()
				client := &Client{
					UserName:  header.Name,
					IP:        v.IP,
					SessionID: p.GetSequence(),
				}
				v.Client = client

				clientDao.New(header.Name, "10000")
				logger.Infof("[VirtualSession.OnPacket] get login request packet. user [%s], seuqence [%d]", header.Name, p.GetSequence())

				res := &game.GameMessage{
					Payload: &game.GameMessage_LoginRes{
						&game.LoginRes{
							Code:     uint32(0),
							Username: header.Name,
						},
					},
				}

				resPacket := packet.NewProtobufPacket(p.GetSequence(), res)
				if resPacket == nil {
					return errors.New("[VirtualSession.OnPacket] fail to new protobuf packet")
				}

				if err = v.Connection.SendPackage(v.ClientID, resPacket.ToByte()); err != nil {
					return err
				}

				return nil
			}
		}

		v.handleProtobuf(gameMessage)

		logger.Infof("[VirtualSession.OnPacket] get proto: %+v", gameMessage)

		return nil

	}
	return nil
}

func (v *VirtualSession) handleProtobuf(header *game.GameMessage) error {
	fmt.Printf("handleProtobuf %+v\n", header)
	client := v.Client
	if client == nil {
		fmt.Printf("[Gate] session %#v 收到protobuf client是空，略過\n", v)
		return nil
	}
	//logger.Get().Infof("[Client] %s 收到protobuf %s", client.GetLoginname(), header.String())
	switch header.Payload.(type) {
	case *game.GameMessage_FlipCoinMessage:
		param := header.GetFlipCoinMessage()
		switch param.Payload.(type) {
		case *flipCoin.GameMessage_BetReq:
			//請求下注
			reqHeader := param.GetBetReq()
			res, err := task.Bet(client.IP, client.UserName, client.SessionID, reqHeader)

			if err != nil {
				//失敗了就提前回應，不然要等API扣款完成

				res := &game.GameMessage{
					Payload: &game.GameMessage_FlipCoinMessage{
						&flipCoin.GameMessage{
							Payload: &flipCoin.GameMessage_BetRes{
								&flipCoin.BetRes{
									Code: uint32(111),
								},
							},
						},
					},
				}

				resPacket := packet.NewProtobufPacket(0, res)

				if err = v.Connection.SendPackage(v.ClientID, resPacket.ToByte()); err != nil {
					return err
				}
			}

			resPacket := packet.NewProtobufPacket(0, res)

			if err = v.Connection.SendPackage(v.ClientID, resPacket.ToByte()); err != nil {
				logger.Errorf("[VirtualSession] %s 發送扣款回應給玩家失敗: %v", err)
				return err
			}

			return nil
		}
		return nil
	}
	return nil
}

func (v *VirtualSession) Write(body []byte) error {

	return v.Connection.SendPackage(v.ClientID, body)

}
