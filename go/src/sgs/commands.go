package sgs

const (

	//CMD_TICK tick command
	CMD_TICK = CMD_C_SYSTEM | (iota + 1)

	//CMD_FORWARD_TO_APP forward client sent data
	CMD_FORWARD_TO_APP

	//CMD_FORWARD_TO_CLIENT data is to be forwarded
	CMD_FORWARD_TO_CLIENT

	//CMD_APP_RUN allow app to run
	CMD_APP_RUN

	_CMD_CLOSE_NET_CLIENT

	_CMD_CLIENT_RECONNECT
)

//PlForwardToClient payload for forward to command
type PlForwardToClient struct {
	ClientID int
	Cmd      Command
}
