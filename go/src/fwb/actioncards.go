package main

const (
	AT_FARM            int = AT_C_PRODUCTION | 0x01
	AT_FEED_SHEEP      int = AT_C_PRODUCTION | 0x02
	AT_DO_PARTTIME_JOB int = AT_C_PRODUCTION | 0x03
	AT_TAKE_VACATION   int = AT_C_CONSUMPTION | 0x01
)

var farmIncome = farmIncomeCopy()

var acFarm = actionCard{
	acType: AT_FARM,
	slot:   3,
	fs:     3,
	wps:    1,

	cost: playerDF{0, 0, 0, 0, 0, 0, 0, 0, -1},
	ret:  playerDF{0, 0, 0, 0, 0, 0, 0, 0, 0},
	si:   &farmIncome,
}

var feedSheepIncome = feedSheepIncomeCopy()

var acFeedSheep = actionCard{
	acType: AT_FEED_SHEEP,
	slot:   3,
	fs:     3,
	wps:    1,

	cost: playerDF{0, 0, 1, 0, 0, 0, 0, 0, -1},
	ret:  playerDF{0, 0, 0, 0, 0, 0, 0, 0, 0},
	si:   &feedSheepIncome,
}

var takeVacationIncome = takeVacationIncomeCopy()

var acTakeVacation = actionCard{
	acType: AT_TAKE_VACATION,
	slot:   5,
	fs:     5,
	wps:    1,

	cost: playerDF{0, 5, 0, 0, 0, 0, 0, 0, -1},
	ret:  playerDF{0, 0, 0, 0, 0, 0, 0, 0, 0},
	si:   &takeVacationIncome,
}

var doParttimeJobIncome = doParttimeJobIncomeCopy()

var acDoParttimeJob = actionCard{
	acType: AT_DO_PARTTIME_JOB,
	slot:   3,
	fs:     3,
	wps:    1,

	cost: playerDF{0, 0, 0, 0, 0, 0, 0, 0, -1},
	ret:  playerDF{0, 0, 0, 0, 0, 0, 0, 0, 0},
	si:   &doParttimeJobIncome,
}
