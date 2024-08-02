package commander

import (
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"fmt"
	"strings"
)

var (
	currentMap    *maphandler.MapInfo
	currentDB     *playerhandler.Database
	currentPlayer *playerhandler.Player
)

// Give the commander the map, db, and player
func Init(initMap *maphandler.MapInfo, initDB *playerhandler.Database, initPlayer *playerhandler.Player) {
	// Initialise map
	currentMap = initMap

	// Load / Create the database for the map
	currentDB = initDB

	// Load the player from the database
	currentPlayer = initPlayer

	fmt.Println(currentPlayer.Name)
	fmt.Println(currentPlayer.Inventory)
}

// Get players inventory -- for displaying
func GetCurrPlayerInv() []string {
	return currentPlayer.Inventory
}

// Get players name -- for messages
func GetCurrPlayerName() string {
	return currentPlayer.Name
}

// When the player gives a command, this handles it and returns a string to be shown to the player in response
func PlayerCommand(cmd string) string {
    trimmedCmd := strings.TrimSpace(cmd)
    cleanedCmd := strings.Fields(trimmedCmd)
	if len(cleanedCmd) > 2 { // This is so dumb
		return "Hmm..."
	}

    // Command prefix switch
	switch cleanedCmd[0] {
    case "go": // Switch rooms
        if len(cleanedCmd) == 1 { return "Please specify a direction" }
		if !currentMap.MoveDirection(cleanedCmd[1], currentPlayer.Inventory) {
			return "Could not move there"
		}
		return "Moved to " + currentMap.CurrentRoom.Name
	case "look": // Give us the look of the room - will remain the same - Maybe this should be called immediatly on entering...
		return currentMap.CurrentRoom.Look
	case "get": // Get the item in the room
        if len(cleanedCmd) == 1 { return "Please specify an item" }
		if !currentMap.ItemInRoom(cleanedCmd[1]) {
			return "That doesn't appear to be here."
		}
		if currentPlayer.IsInInv(cleanedCmd[1]) {
			return "You already have " + cleanedCmd[1]
		}
        // Save whatever item we get from the room
		currentPlayer.AddToInv(cleanedCmd[1])
		currentDB.SavePlayerInfo(currentPlayer)
		return "Got " + cleanedCmd[1]
	case "whereami": // Duh
		return currentMap.CurrentRoom.Name
	default:
		return "Hmm..."
	}
}
