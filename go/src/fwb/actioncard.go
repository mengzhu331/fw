package main

const (
	AT_C_ALL         int = 0xffffffff
	AT_C_PRODUCTION  int = 0x01000000
	AT_C_CONSUMPTION int = 0x02000000
	AT_C_TRADING     int = 0x04000000
	AT_C_FINANCING   int = 0x08000000
	AT_C_LOTTERY     int = 0x10000000
	AT_C_SPECIAL     int = 0x20000000
)

type actionCard struct {
	acType int
	slot   int
	fs     int
	wps    int

	cost playerDF
	ret  playerDF
	si   settlementItem
	prob float32
}

func (this *actionCard) claimedBy(p player) error {
	if this.fs >= this.wps {
		this.fs -= this.wps
	} else {
		return MakeFwErrorByCode(EC_USE_FULL_AC)
	}

	df := p.getDF()
	if !df.afford(this.cost) {
		return MakeFwErrorByCode(EC_USE_AC_NOT_AFFORD)
	}

	df.apply(&this.cost)
	df.apply(&this.ret)
	p.setDF(df)

	return nil
}
