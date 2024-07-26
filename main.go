package main

import (
	"explore/internal/commander"
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"fmt"
)

func main() {
	Run()
}

func Run() error {
	// Menu ?

	// Get menu option
	// 1. List map(s)
	chosenMap := "prologue"
	initMap, err := maphandler.InitNewMap("maps/" + chosenMap + ".json")
	if err != nil {
		fmt.Println("THERE WAS AN ERROR - MAP")
		return err
	}

	// 2. Open Database for that map and list players + add a "create new"
	DB := playerhandler.LoadDatabase("databases/" + chosenMap)

	// 3. Get player name / create a player name
	chosenName := "player1"
	DB.CreatePlayer(chosenName) // This is here cause we don't have a menu with a "create player" option, so we need to create one every time
	player, err := DB.LoadPlayer(chosenName)
	if err != nil {
		return err
	}

    // Commander
	commander.Init(initMap, DB, player) // Run tui once commander is setup

    // TUI
	//tui.Start()

	return nil
}
