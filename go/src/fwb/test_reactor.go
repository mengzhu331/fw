package main

import "sgs"

type trigger struct {
	state       int
	actionCmd   sgs.Command
	cts         []CommandTarget
	targetState int
}
type reactor struct {
	cm    map[sgs.Command]*trigger
	state int
}

func (this *reactor) sendCommand(cmd sgs.Command) error {
	t := this.cm[cmd]
	if t != nil && t.state == this.state {
		err := this.dispatchCommand(t.actionCmd, t.cts)
		if err != nil {
			return err
		}
		this.state = t.targetState
	}
	return nil
}

func (this *reactor) dispatchCommand(command sgs.Command, cts []CommandTarget) error {
	for _, t := range cts {
		err := t.SendCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}
