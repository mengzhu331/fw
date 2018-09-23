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

	//ET_ILLEGAL_PARAMETER parameter is invalid for a request
	ET_ILLEGAL_PARAMETER = (iota + 1) << 16

	//ET_INVALID_STATE system internal state is invalid
	ET_INVALID_STATE

	//ET_INVALID_REQUEST request should not be received according to the system state
	ET_INVALID_REQUEST

	//ET_TIMEOUT have not obtained response or request within required time frame
	ET_TIMEOUT

	//ET_NETWORK network error
	ET_NETWORK

	//ET_MISSING_DATA mandatory data for sofware is missing
	ET_MISSING_DATA
)
