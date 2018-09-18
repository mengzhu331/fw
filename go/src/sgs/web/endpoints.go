package web

const (

	//EP_LOGIN login endpoint
	EP_LOGIN int = 1

	//EP_WEBSOCKET open websocket endpoint
	EP_WEBSOCKET int = 2

	//EP_JOINSESSION join session endpoint
	EP_JOINSESSION int = 3
)

//EndPointMap configuration values for endpoints
type EndPointMap map[int]string

var _depm EndPointMap

func demp() EndPointMap {
	if _depm == nil {
		_depm = make(EndPointMap)
		_depm[EP_LOGIN] = "/login"
		_depm[EP_JOINSESSION] = "/join_session"
		_depm[EP_WEBSOCKET] = "/ws"
	}
	return _depm
}
