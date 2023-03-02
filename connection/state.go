package connection

type State int

const (
	State_Connecting    State = iota //State_Connecting 連線中
	State_Authorization              //State_Authorization 授權中
	State_Connected                  //State_Connected 已連線
	State_Disconnected               //State_Disconnected 已斷線
)
