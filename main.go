package main

import (
	"explore/internal/commander"
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
)

func main() {
	// Menu ?

	// Get menu option
	// 1. List map(s)
	chosenMap := "prologue"
	initMap, err := maphandler.InitNewMap("maps/"+chosenMap+".json")
	if err != nil {
		return
	}

	// 2. Open Database for that map and list players + add a "create new"
	DB := playerhandler.LoadDatabase("databases/"+chosenMap)

	// 3. Get player name / create a player name
	chosenName := "player1"
	player, err := DB.LoadPlayer(chosenName)
	if err != nil {
		return
	}

	commander.Init(*initMap, *DB, *player) // Run tui once commander is setup
	//tui.Start()
}
