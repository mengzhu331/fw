package main

const (
	ST_FARM_INCOME            int = ST_C_INCOME | 0x01
	ST_FEED_SHEEP_INCOME      int = ST_C_INCOME | 0x02
	ST_TAKE_VACATION_INCOME   int = ST_C_INCOME | 0x03
	ST_DO_PARTTIME_JOB_INCOME int = ST_C_INCOME | 0x04
)

func farmIncomeCopy() siSimple {
	return siSimple{
		effect:  playerDF{0, 0, 2, 0, 0, 0, 0, 0, 0},
		counter: 1,
		siType:  ST_FARM_INCOME,
	}
}

func feedSheepIncomeCopy() siSimple {
	return siSimple{
		effect:  playerDF{0, 0, 0, 1, 1, 1, 0, 0, 0},
		counter: 1,
		siType:  ST_FEED_SHEEP_INCOME,
	}
}

func takeVacationIncomeCopy() siSimple {
	return siSimple{
		effect:  playerDF{1, 0, 0, 0, 0, 0, 0, 0, 0},
		counter: 1,
		siType:  ST_TAKE_VACATION_INCOME,
	}
}

func doParttimeJobIncomeCopy() siSimple {
	return siSimple{
		effect:  playerDF{0, 10, 0, 0, 0, 0, 0, 0, 0},
		counter: 1,
		siType:  ST_DO_PARTTIME_JOB_INCOME,
	}
}
