package main

import (
	"./tools"
)

func main() {
	gw := tools.CreateGroupper("2020-11-15", "d:/777", 400)
	gw.Start()
}