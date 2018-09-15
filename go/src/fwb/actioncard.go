package main

type actionCard struct {
	slot int
	fs   int

	cost playerDF
	ret  playerDF
	prob float32
}

func (this *actionCard) usedBy(pdf *playerDF) error {
	if this.fs > 0 {
		this.fs--
	} else {
		return MakeFwErrorByCode(EC_USE_FULL_AC)
	}

	if !pdf.afford(this.cost) {
		return MakeFwErrorByCode(EC_USE_AC_NOT_AFFORD)
	}

	pdf.apply(&this.cost)
	pdf.apply(&this.ret)
	return nil
}
