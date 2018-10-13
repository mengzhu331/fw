package fwb

import "er"

const (
	_E_FWB = 0x8000

	//E_CMD_INVALID_CLIENT the client ID is not found in current game, or the client should not involve current process
	E_CMD_INVALID_CLIENT = _E_FWB | er.IMPT_REMARKABLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x1

	//E_INVALID_CMD the command is not recognized
	E_INVALID_CMD = _E_FWB | er.IMPT_REMARKABLE | er.ET_INTERNAL | er.EI_INVALID_REQUEST | 0x2

	//E_CMD_NOT_EXECUTABLE the command is not supposed to be received and/or executed
	E_CMD_NOT_EXECUTABLE = _E_FWB | er.IMPT_RECOVERABLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x3

	//E_CMD_PAYLOAD_NOT_ENCODABLE failed to encode the payload of the command
	E_CMD_PAYLOAD_NOT_ENCODABLE = _E_FWB | er.IMPT_THREAT | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x4

	//E_CMD_PAYLOAD_NOT_DECODABLE failed to decode the payload of the command
	E_CMD_PAYLOAD_NOT_DECODABLE = _E_FWB | er.IMPT_THREAT | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x5

	//E_MISSING_GAME_SETTINGS mandatory game settings are not provided
	E_MISSING_GAME_SETTINGS = _E_FWB | er.IMPT_UNRECOVERABLE | er.ET_SETTINGS | er.EI_MISSING_DATA | 0x6

	//E_INVALID_SETTINGS mandatory game settings are not valid
	E_INVALID_SETTINGS = _E_FWB | er.IMPT_UNRECOVERABLE | er.ET_SETTINGS | er.EI_ILLEGAL_PARAMETER | 0x7

	//E_PROFILE_MISMATCH the game profile does not match game settings
	E_PROFILE_MISMATCH = _E_FWB | er.IMPT_UNRECOVERABLE | er.ET_SETTINGS | er.EI_ILLEGAL_PARAMETER | 0x8

	//E_ACTION_PARSER_NOT_FOUND no parser defined for the action ID
	E_ACTION_PARSER_NOT_FOUND = _E_FWB | er.IMPT_UNRECOVERABLE | er.ET_INTERNAL | er.EI_INCORRECT_STATIC_DATA | 0x9

	//E_INVALID_ACTION the action is not supposed to be received and/or executed
	E_INVALID_ACTION = _E_FWB | er.IMPT_THREAT | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x10

	//E_DOACTION_INVALID_CLIENTID when do action the specified player ID is illegal
	E_DOACTION_INVALID_CLIENTID = _E_FWB | er.IMPT_DEGRADE | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x11
)
