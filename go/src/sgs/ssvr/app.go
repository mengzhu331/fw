package ssvr

//App interface for applications
type App interface {
	Init(chan string, []NetClient) error
	SendCommand(Command) error
}

//AppBuildFunc function to build an application
type AppBuildFunc func() App
