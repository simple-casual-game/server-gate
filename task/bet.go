package task

import (
	"errors"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"github.com/simple-casual-game/server-gate/dao/clientDao"
	"github.com/simple-casual-game/server-gate/logger"
	"github.com/simple-casual-game/server-gate/protobuf/flipCoin"
	"github.com/simple-casual-game/server-gate/protobuf/game"
)

func Bet(ip uint32, username string, sessionID uint32, reqHeader *flipCoin.BetReq) (*game.GameMessage, error) {

	clientModel := clientDao.Get(username)

	betAmount, err := decimal.NewFromString(reqHeader.Betamount)
	if err != nil {
		return nil, err
	}

	balance, err := decimal.NewFromString(clientModel.Amount)
	if err != nil {
		return nil, err
	}

	if balance.LessThan(betAmount) {
		return nil, errors.New("balance not enough")
	}

	balance = balance.Sub(betAmount)
	clientDao.Modify(username, balance.String())

	rand.Seed(time.Now().UnixNano())

	result := rand.Int() % 2
	winAmount := decimal.Zero

	if result == 1 {
		// win
		winAmount = betAmount.Mul(decimal.NewFromInt(2))
		balance.Add(winAmount)
		clientDao.Modify(username, balance.String())
		logger.Infof("[Bet] win %s\n", betAmount.Mul(decimal.NewFromInt(2)).String())
	} else {
		// lose
		logger.Infof("[Bet] lose %s\n", betAmount.String())
	}

	logger.Infof("[Bet] %s 發送扣款回應給玩家\n", clientModel.Username)
	//發送回應給client
	res := &game.GameMessage{
		Payload: &game.GameMessage_FlipCoinMessage{
			&flipCoin.GameMessage{
				Payload: &flipCoin.GameMessage_BetRes{
					&flipCoin.BetRes{
						Betamount: betAmount.String(),
						Code:      uint32(0),
						Bettime:   uint64(time.Now().Unix()),
						Win:       winAmount.String(),
					},
				},
			},
		},
	}

	return res, nil
}
