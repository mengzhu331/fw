package sgs

import (
	"fmt"
)

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

	//CMD_C_APP_PRIVATE the command is app specific
	CMD_C_APP_PRIVATE = CMD_C_APP | 0x00020000
)

//Command common command object
type Command struct {
	ID      int
	Source  int
	Payload interface{}
}

//InCategory if the command category includes the command
func (me *Command) InCategory(c int) bool {
	return (me.ID & c) == c
}

//HexID make hex string for ID of the command
func (me *Command) HexID() string {
	return fmt.Sprintf("%x", me.ID)
}
