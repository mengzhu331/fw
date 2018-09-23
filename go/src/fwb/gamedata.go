package main

const (
	_PT_HEART = iota
	_PT_GOLD
	_PT_CEREALS
	_PT_MEAT

	_PT_WOO
	_PT_LEATHER
	_PT_ROBE
	_PT_HAT

	_PT_MAX
)

type propertySet [_PT_MAX]int

type playerData struct {
	playerID int
	property propertySet
	worker   int
	employee int
	houseLv  int
}

func (me propertySet) covers(p propertySet) bool {
	for i := 0; i < len(me); i++ {
		if me[i] < p[i] {
			return false
		}
	}
	return true
}

func (me propertySet) add(p propertySet) propertySet {
	sum := propertySet{}
	for i := 0; i < len(me); i++ {
		sum[i] = me[i] + p[i]
	}
	return sum
}

func (me propertySet) sub(p propertySet) propertySet {
	delta := propertySet{}
	for i := 0; i < len(me); i++ {
		delta[i] = me[i] - p[i]
	}
	return delta
}

type cardData struct {
	cardID   int
	slot     int
	usedSlot int
}

type gameData struct {
	round int
	pData map[int]playerData
	cards []card
}

func initPlayerData() playerData {
	data := playerData{}

	data.worker = 5
	data.employee = 0
	data.houseLv = 0

	data.property[_PT_HEART] = 0
	data.property[_PT_GOLD] = 30
	data.property[_PT_CEREALS] = 5
	data.property[_PT_MEAT] = 0
	data.property[_PT_WOO] = 0
	data.property[_PT_LEATHER] = 0
	data.property[_PT_HAT] = 0
	data.property[_PT_ROBE] = 0

	return data
}
