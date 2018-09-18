package main

type acBuff interface {
	getEffect(*actionCard) (bool, playerDF, playerDF)
	countDown() bool
	attach(*player)
}

type acBuffSimple struct {
	acTypeMask int
	retOffset  playerDF
	costOffset playerDF
	counter    int
}

func (this *acBuffSimple) getEffect(ac *actionCard) (bool, playerDF, playerDF) {
	if (this.acTypeMask & ac.acType) != ac.acType {
		return false, playerDF{}, playerDF{}
	}

	return true, this.retOffset, this.costOffset
}

func (this *acBuffSimple) countDown() bool {
	this.counter--
	return this.counter != 0
}

func (this *acBuffSimple) attach(p *player) {}
