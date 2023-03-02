package clientDao

import "errors"

var clients map[string]*ClientModel = make(map[string]*ClientModel)

type ClientModel struct {
	Username string
	Amount   string
}

func New(username string, amount string) error {
	clients[username] = &ClientModel{
		Username: username,
		Amount:   amount,
	}
	return nil
}

func Get(username string) *ClientModel {
	return clients[username]
}

func Modify(username string, amount string) error {
	if _, ok := clients[username]; !ok {
		return errors.New("no such client")
	}
	clients[username].Amount = amount
	return nil
}
