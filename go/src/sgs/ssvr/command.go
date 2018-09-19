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

//Command common command object
type Command struct {
	ID    int
	Param interface{}
}

//TickParam parameter with tick
type TickParam struct {
	DeltaMs int
}
