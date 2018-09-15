package main

type simpleBuff interface {
	getBuffType() int
	getTargetPhase() uint64
	attach(*player)
	trigger() (playerDF, bool)
	getEffect() playerDF
}

type actionCardBuff interface {
	getBuffType() int
	getACType() int
	attach(*player)
	trigger() (playerDF, playerDF, bool)
	getEffect() (playerDF, playerDF)
}

type buffTriggerBuff interface {
	getBuffType() int
	getTargetBuffType() int
	attach()
}
