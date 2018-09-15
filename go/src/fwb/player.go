package main

import (
	"sgs"
)

type player interface {
	getName() string
	getDF() *playerDF
	SendCommand(sgs.Command) error
}

type playerDF struct {
	name    string
	heart   int
	cereals int
	gold    int
	meat    int

	woo      int
	leather  int
	clothe   int
	worker   int
	employee int
}

func (this *playerDF) apply(offset *playerDF) {
	this.cereals += offset.cereals
	this.clothe += offset.clothe
	this.employee += offset.employee
	this.gold += offset.gold
	this.heart += offset.heart
	this.leather += offset.leather
	this.meat += offset.meat
	this.woo += offset.woo
	this.worker += offset.worker
}

func (this *playerDF) afford(df playerDF) bool {
	return this.cereals+df.cereals >= 0 &&
		this.clothe+df.clothe >= 0 &&
		this.employee+df.employee >= 0 &&
		this.gold+df.gold >= 0 &&
		this.heart+df.heart >= 0 &&
		this.leather+df.leather >= 0 &&
		this.meat+df.meat >= 0 &&
		this.woo+df.woo >= 0
}
