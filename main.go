package main

import (
	"explore/internal/commander"
	"explore/internal/tui"
)


func main() {
    // Menu
    // Get menu option

    commander.Init("./maps/prologue.json") // Run tui once commander is setup

    tui.Start()
}
