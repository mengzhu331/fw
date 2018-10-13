package er

const (
	//E_IMPORTANCE importance bits
	E_IMPORTANCE = 0xF << 24

	//E_TYPE type bits
	E_TYPE = 0xF << 20

	//E_CAUSE cause bits
	E_CAUSE = 0xF << 16
)

const (

	//IMPT_NONE no threat
	IMPT_NONE = 0

	//IMPT_ACCEPTIBLE the error is a sensible variation of normal situation
	IMPT_ACCEPTIBLE = (iota + 1) << 24

	//IMPT_RECOVERABLE system is able to recover from the error
	IMPT_RECOVERABLE

	//IMPT_REMARKABLE system is able to recover from the error, but it should be taken care of
	IMPT_REMARKABLE

	//IMPT_THREAT system is able to recover from the error, but it can potentially threaten the system
	IMPT_THREAT

	//IMPT_DEGRADE system has to degrade to recover from the error
	IMPT_DEGRADE

	//IMPT_UNRECOVERABLE system is not able to recover from the error
	IMPT_UNRECOVERABLE
)

const (

	//ET_INTERNAL program internal error
	ET_INTERNAL = (iota + 1) << 20

	//ET_INTERACTION error from systems interaction
	ET_INTERACTION

	//ET_SERVICE failed to call service
	ET_SERVICE

	//ET_SETTINGS error with sofware settings
	ET_SETTINGS
)

const (

	//EI_ILLEGAL_PARAMETER parameter is invalid for a request
	EI_ILLEGAL_PARAMETER = (iota + 1) << 16

	//EI_INVALID_STATE system internal state is invalid
	EI_INVALID_STATE

	//EI_INCORRECT_STATIC_DATA static data of the program is incomplete or incorrect
	EI_INCORRECT_STATIC_DATA

	//EI_INVALID_REQUEST request should not be received according to the system state
	EI_INVALID_REQUEST

	//EI_TIMEOUT have not obtained response or request within required time frame
	EI_TIMEOUT

	//EI_NETWORK network error
	EI_NETWORK

	//EI_MISSING_DATA mandatory data for sofware is missing
	EI_MISSING_DATA

	//EI_ILLEGAL_RESULT result is not valid
	EI_ILLEGAL_RESULT

	//EI_SYSTEM the based hardware or software system fail
	EI_SYSTEM
)
