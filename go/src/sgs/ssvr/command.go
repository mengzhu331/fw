package ssvr

const (
	//CMD_CATEGORY command category bits
	CMD_CATEGORY = 0xff000000

	//CMD_C_SYSTEM system command category
	CMD_C_SYSTEM = 0x08000000

	//CMD_C_CLIENT client command category
	CMD_C_CLIENT = 0x04000000

	//CMD_C_APP application command category
	CMD_C_APP = 0x02000000
)

const (
	//CMD_C_CLIENT_TO_APP the command is sent to backend app from client
	CMD_C_CLIENT_TO_APP = CMD_C_CLIENT | 0x00010000

	//CMD_C_APP_TO_CLIENT the command is sent to client from backend app
	CMD_C_APP_TO_CLIENT = CMD_C_APP | 0x00010000
)

const (

	//CMD_TICK tick command
	CMD_TICK = CMD_C_SYSTEM | (iota + 1)

	//CMD_FORWARD_TO_APP forward client sent data
	CMD_FORWARD_TO_APP

	//CMD_FORWARD_TO_CLIENT data is to be forwarded
	CMD_FORWARD_TO_CLIENT

	//CMD_APP_RUN allow app to run
	CMD_APP_RUN
)

//Command common command object
type Command struct {
	ID      int
	Source  int
	Payload interface{}
}

//PlForwardToClient payload for forward to command
type PlForwardToClient struct {
	ClientID int
	Payload  []byte
}

//CmdInCategory if the command category includes the command
func CmdInCategory(command Command, c int) bool {
	return (command.ID & c) == c
}
