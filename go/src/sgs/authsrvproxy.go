package sgs

type authSrvPrx interface {
	vclient(clientID int, token string) string
}

type sgasPrx struct{}

func (me *sgasPrx) vclient(clientID int, token string) string {
	return ""
}

func createSgasPrx() *sgasPrx {
	return &sgasPrx{}
}
