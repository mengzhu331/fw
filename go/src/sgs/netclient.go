package sgs

//NetConn network interface independent of protocol
type NetConn interface {
	Send(cmd Command) error
	Run(ch chan Command, mch chan Command)
}

type netClient struct {
	id       int
	username string
	conn     NetConn
	s        *session
	mch      chan Command
}

func (me *netClient) send(cmd Command) error {
	return me.conn.Send(cmd)
}

func (me *netClient) run(ch chan Command) {
	go me.conn.Run(ch, me.mch)
}

func (me *netClient) close() {
	me.mch <- Command{
		ID: _CMD_CLOSE_NET_CLIENT,
	}
}
