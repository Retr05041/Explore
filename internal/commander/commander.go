package commander

import (
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"fmt"
	"strings"
)

var (
	currentMap    maphandler.MapInfo
	currentDB     playerhandler.Database
	currentPlayer playerhandler.Player
)

func Init(initMap maphandler.MapInfo, initDB playerhandler.Database, initPlayer playerhandler.Player) {
	// Initialise map
	currentMap = initMap

	// Load / Create the database for the map
	currentDB = initDB

	// Load the player from the database
	currentPlayer = initPlayer

	fmt.Println(currentPlayer.Name)
	fmt.Println(currentPlayer.Inventory)
}

func GetCurrPlayerInv() []string {
	return currentPlayer.Inventory
}

func PlayerCommand(cmd string) string {
	splitCmd := strings.Split(cmd, " ")
	if len(splitCmd) > 2 {
		return "Hmm..."
	}

	for _, token := range splitCmd {
		strings.ToLower(strings.ReplaceAll(token, " ", ""))
		if token == "" {
			splitCmd[len(splitCmd)-1] = "UNKNOWN"
		} // Bogus check to see if there is a whitespace element - needs to be after the length check to work :(
	}

	switch splitCmd[0] {
	case "go":
		if !currentMap.MoveDirection(splitCmd[1]) {
			return "Could not move there"
		}
		return "Moved to " + currentMap.CurrentRoom.Name
	case "look":
		return currentMap.CurrentRoom.Look
	case "whereami":
		return currentMap.CurrentRoom.Name
	default:
		return "Hmm..."
	}
}
