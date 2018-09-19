package main

import (
	"encoding/json"
	"sgs"
)

func main() {
	s := "{\"ID\": 1,\"Param\":{\"t\": \"s\"}}"

	type Cmd struct {
		ID    int
		Param interface{}
	}

	c := Cmd{}

	json.Unmarshal([]byte(s), &c)

	sgs.Run(fwAppBuildFunc)
}
