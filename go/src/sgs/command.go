package sgs

const (
	CMD_TIER_FRAMEWORK_INTERNAL    uint = 0x80000000
	CMD_TIER_FRAMEWORK_APPLICATION uint = 0x00010000
	CMD_ANY                        uint = 0x00000000
)

const (
	CMD_TICK uint = CMD_TIER_FRAMEWORK_APPLICATION | 0x0001
)

type Command struct {
	ID     uint
	Source string
	Target string
	Param  interface{}
}
