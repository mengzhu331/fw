package core

type phase int32

const (
	_P_GAME_START = iota << 24
	_P_GAME_ROUNDS
	_P_GAME_SETTLEMENT
	_P_GAME_FINISH
)

const (
	_P_ROUNDS_START = _P_GAME_ROUNDS | (iota << 16)
	_P_ROUNDS_TURNS
	_P_ROUNDS_SETTLEMENT
	_P_ROUNDS_FINISH
)
