package main

import "sgs"

const (
	P_GAME_START      uint64 = 0x01000000
	P_GAME_ROUNDS     uint64 = 0x02000000
	P_GAME_SETTLEMENT uint64 = 0x03000000
	P_GAME_FINISH     uint64 = 0x04000000

	P_ROUNDS_START      uint64 = P_GAME_ROUNDS | 0x010000
	P_ROUNDS_TURNS      uint64 = P_GAME_ROUNDS | 0x020000
	P_ROUNDS_SETTLEMENT uint64 = P_GAME_ROUNDS | 0x030000
	P_ROUNDS_FINISH     uint64 = P_GAME_ROUNDS | 0x040000
)

type fwGameExecutor func(*fwGame, sgs.Command) error

type phaseMap map[uint64]fwGameExecutor

func phaseDispatcher(this *fwGame, pm phaseMap, phase uint64, command sgs.Command) error {
	var mask uint64 = 0xff
	invertmask := (^mask)
	for mask = 0xff; mask != 0xff00000000000000; mask <<= 8 {
		if (phase&mask) != 0 && pm[phase] != nil {
			err := pm[phase](this, command)
			if err != nil {
				return err
			}
		}
		phase = phase & invertmask
		invertmask <<= 8
	}
	return nil
}

func phaseDispatcherSimple(this *fwGame, pm phaseMap, phase uint64, command sgs.Command) error {
	exec := pm[phase]
	return exec(this, command)
}
