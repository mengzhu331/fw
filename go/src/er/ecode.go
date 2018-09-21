package er

const (
	//ImptAcceptible the error is a sensible variation of normal situation
	ImptAcceptible = (iota + 1) << 24

	//ImptRecoverable system is able to recover from the error
	ImptRecoverable

	//ImptRemarkable system is able to recover from the error, but it should be taken care of
	ImptRemarkable

	//ImptThreat system is able to recover from the error, but it can potentially threaten the system
	ImptThreat

	//ImptDegrade system has to degrade to recover from the error
	ImptDegrade

	//ImptUnrecoverable system is not able to recover from the error
	ImptUnrecoverable
)

const (

	//ETInternal program internal error
	ETInternal = (iota + 1) << 20

	//ETInteraction error from systems interaction
	ETInteraction

	//ETConsumeService failed to consume service
	ETConsumeService
)

const (

	//ETIllegalParameter parameter is invalid for a request
	ETIllegalParameter = (iota + 1) << 16

	//ETInvalidState system internal state is invalid
	ETInvalidState

	//ETInvalidRequest request should not be received according to the system state
	ETInvalidRequest

	//ETTimeOut have not obtained response or request within required time frame
	ETTimeOut
)
