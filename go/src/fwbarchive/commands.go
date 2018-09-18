package main

import (
	"strings"
)

const (
	CMD_GAME_START                  uint = 0x0001
	CMD_GAME_START_ACK              uint = 0x0002
	CMD_ENTER                       uint = 0x0003
	CMD_GAME_END_PLAYER_NO_RESPONSE uint = 0x0004

	CMD_CLAIM_ACTION_CARD uint = 0x0010
)

const (
	TARGET_UNKNOWN string = "UNK:"
	TARGET_PLAYER  string = "PLY:"
	TARGET_GAME    string = "FWG:"
)

func extractCommandParticipant(participantUri string) (string, string) {
	protocol := strings.ToUpper(participantUri[0:4])

	switch protocol {
	case "FWG:":
		return protocol, ""
	case "PLY:":

		if strings.ToUpper(participantUri) == "PLY:%ALL" {
			return protocol, ""
		}

		name := participantUri[4:]
		return protocol, name
	case "SYS:":
		return participantUri, ""
	}
	return TARGET_UNKNOWN, ""
}

func makeCommandParticipantUri(participant string, name string) string {
	if participant == TARGET_PLAYER && name == "" {
		return TARGET_PLAYER + "%ALL"
	} else if participant == TARGET_PLAYER {
		return TARGET_PLAYER + name
	} else if participant == TARGET_GAME {
		return TARGET_GAME
	}

	return "UNK:"
}
