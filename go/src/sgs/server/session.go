package server

type session struct {
	cch     chan string
	clients clientMap
	running bool
}
