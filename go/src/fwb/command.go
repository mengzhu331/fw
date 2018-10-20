package fwb

import (
	"sgs"
)

const (
	_CMD_WHO_GAME = 0x80000000
)

const (
	//CMD_TIMER command for programatically set timeout
	CMD_TIMER = sgs.CMD_C_APP_PRIVATE + iota + 1
)

const (
	//CMD_GAME_START command for notifying players that game started
	CMD_GAME_START = sgs.CMD_C_APP_TO_CLIENT + iota + 1

	//CMD_GAME_OVER command for notifying players that game is up
	CMD_GAME_OVER

	//CMD_ROUND_START command for notifying players that a new round started
	CMD_ROUND_START

	//CMD_START_TURN command for notifying players that a new turn started
	CMD_START_TURN

	//CMD_ACTION_REJECTED command for notifying players that the requested action is rejected
	CMD_ACTION_REJECTED

	//CMD_ACTION_COMMITTED command for notifying players that the action has been committed
	CMD_ACTION_COMMITTED

	//CMD_ROUND_SETTLEMENT command for notifying players to make settlement
	CMD_ROUND_SETTLEMENT

	//CMD_ROUND_SETTLEMENT_INVALID command for notifying a player that the data with the requested settlement is invalid
	CMD_ROUND_SETTLEMENT_INVALID

	//CMD_ROUND_SETTLEMENT_UPDATE command for notifying players that the settlement has updated the game data
	CMD_ROUND_SETTLEMENT_UPDATE

	//CMD_GAME_FINISH command for notifying players that game has finished with player rank ready
	CMD_GAME_FINISH

	//CMD_SYNC_GAME_STATE command for synchronize game state between components
	CMD_SYNC_GAME_STATE
)

const (
	//CMD_GAME_START_ACK command for notifying the game system that the player acknowledged game having started
	CMD_GAME_START_ACK = sgs.CMD_C_CLIENT_TO_APP + iota + 1

	//CMD_ROUND_START_ACK command for notifying the game system that the player ackowledged a new round having started
	CMD_ROUND_START_ACK

	//CMD_ACTION command for requesting an action
	CMD_ACTION

	//CMD_COMMIT_ROUND_SETTLEMENT command for commiting round settlement
	CMD_COMMIT_ROUND_SETTLEMENT

	//CMD_REMATCH command for requesting next match with the same players
	CMD_REMATCH
)

const (
	//CMD_WHO_APP the command is sent from the game system to the players or other systems
	CMD_WHO_APP = (_CMD_WHO_GAME | 0x8000) + iota + 1
)
