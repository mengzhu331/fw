package main

import (
	"sgs"
)

const (
	PD_HEART    int = 0
	PD_GOLD     int = 1
	PD_CEREALS  int = 2
	PD_MEAT     int = 3
	PD_WOO      int = 4
	PD_LEATHER  int = 5
	PD_JACK     int = 6
	PD_WORKER   int = 7
	PD_EMPLOYEE int = 8

	PD_MAX_BASIC int = PD_WORKER
	PD_MAX       int = PD_EMPLOYEE + 1
)

type playerDF [PD_MAX]int

type player interface {
	getName() string
	getDF() *playerDF
	setDF(*playerDF)
	getSI() []settlementItem
	removeSI(int)
	SendCommand(sgs.Command) error
}

func (this *playerDF) apply(offset *playerDF) {
	for i := 0; i < PD_MAX_BASIC; i++ {
		this[i] += offset[i]
	}

	this[PD_WORKER] += offset[PD_WORKER]

	if offset[PD_EMPLOYEE] < 0 {
		this[PD_EMPLOYEE] += offset[PD_EMPLOYEE]
		if this[PD_EMPLOYEE] < 0 {
			this[PD_WORKER] += this[PD_EMPLOYEE]
			this[PD_EMPLOYEE] = 0
		}
	} else {
		this[PD_EMPLOYEE] += offset[PD_EMPLOYEE]
	}
}

func (this *playerDF) afford(df playerDF) bool {
	for i := 0; i < PD_MAX_BASIC; i++ {
		if this[i]+df[i] < 0 {
			return false
		}
	}

	if this[PD_WORKER]+df[PD_WORKER] < 0 {
		return false
	}

	if df[PD_EMPLOYEE]+this[PD_WORKER]+this[PD_EMPLOYEE] < 0 {
		return false
	}

	return true
}
