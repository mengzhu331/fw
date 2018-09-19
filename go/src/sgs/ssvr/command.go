package ssvr

const (

	//CMD_C_SYSTEM system command category
	CMD_C_SYSTEM = 0x08000000

	//CMD_C_CLIENT client command category
	CMD_C_CLIENT = 0x04000000

	//CMD_C_APP application command category
	CMD_C_APP = 0x02000000
)

const (

	//CMD_TICK tick command
	CMD_TICK = CMD_C_SYSTEM | (iota + 1)
)

const (
	//CMD_TO_CLIENT app send to client
	CMD_TO_CLIENT = CMD_C_APP | (iota + 1)
)

const (
	//CMD_FROM_CLIENT client send to app
	CMD_FROM_CLIENT = CMD_C_CLIENT | (iota + 1)
)

//Command common command object
type Command struct {
	ID    int
	Param map[string]interface{}
}

//InCate if command is in category
func (me *Command) InCate(c int) bool {
	return (me.ID & c) == c
}
