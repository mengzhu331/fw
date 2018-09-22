package main

const (
	_P_HEART = iota
	_P_GOLD
	_P_CEREALS
	_P_MEAT

	_P_WOO
	_P_LEATHER
	_P_ROBE
	_P_HAT

	_P_MAX
)

type propertySet [_P_MAX]int

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
	cards []cardData
}
