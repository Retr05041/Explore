package main

import (
	"explore/internal/commander"
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"explore/internal/tui"
)

func main() {
	Run()
}

func Run() error {
	// Menu ?

	// Get menu option
	// 1. List map(s)
	chosenMap := "prologue" // Placeholder
	initMap, err := maphandler.InitNewMap("maps/" + chosenMap + ".json")
	if err != nil {
		return err
	}

	// 2. Open Database for that map and list players + add a "create new"
	DB := playerhandler.LoadDatabase("databases/" + chosenMap) // Need to send metadata here...

	// 3. Get player name / create a player name
	chosenName := "player1" // Placeholder
	DB.CreatePlayer(chosenName) // This is here cause we don't have a menu with a "create player" option, so we need to create one every time
	player, err := DB.LoadPlayer(chosenName)
	if err != nil {
		return err
	}

    // Commander
	commander.Init(initMap, DB, player) // Run tui once commander is setup
    // TUI
	tui.Start()

	return nil
}
