package fwb

const (
	//GAME_OVER_PLAYER_TIMEOUT is a gameover status code, which means the game is terminated because of player response timeout
	GAME_OVER_PLAYER_TIMEOUT = iota + 1

	//GAME_OVER_NORMAL is a gameover status code, which means the game finished after normal procedures
	GAME_OVER_NORMAL
)

const (
	//PROFILE_3PVP game profile for 3 players
	PROFILE_3PVP = "3pvp"

	//PROFILE_2PVP game profile for 2 players dual
	PROFILE_2PVP = "2pvp"
)
