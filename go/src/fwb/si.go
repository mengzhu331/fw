package main

const (
	ST_C_INCOME    int = 0x01000000
	ST_C_FINANCING int = 0x02000000
)

type settlementItem interface {
	getEffect() playerDF
	countDown() bool
	getType() int
}

type siSimple struct {
	effect  playerDF
	counter int
	siType  int
}

func (this *siSimple) getEffect() playerDF {
	return this.effect
}

func (this *siSimple) countDown() bool {
	this.counter--
	return this.counter != 0
}

func (this *siSimple) getType() int {
	return this.siType
}
