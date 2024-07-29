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
	splitCmd := strings.Split(cmd, " ")
	if len(splitCmd) > 2 || len(splitCmd) <= 1 { // This is so dumb
		return "Hmm..."
	}

    for _, token := range splitCmd { // So is this lol
		strings.ToLower(strings.ReplaceAll(token, " ", ""))
		if token == "" {
			splitCmd[len(splitCmd)-1] = "UNKNOWN"
		} // Bogus check to see if there is a whitespace element - needs to be after the length check to work :(
	}

    // Command prefix switch
	switch splitCmd[0] {
    case "go": // Switch rooms
		if !currentMap.MoveDirection(splitCmd[1]) {
			return "Could not move there"
		}
		return "Moved to " + currentMap.CurrentRoom.Name
	case "look": // Give us the look of the room - will remain the same - Maybe this should be called immediatly on entering...
		return currentMap.CurrentRoom.Look
	case "get": // Get the item in the room
		if !currentMap.ItemInRoom(splitCmd[1]) {
			return "That doesn't appear to be here."
		}
		if currentPlayer.IsInInv(splitCmd[1]) {
			return "You already have " + splitCmd[1]
		}
        // Save whatever item we get from the room
		currentPlayer.AddToInv(splitCmd[1])
		currentDB.SavePlayerInfo(currentPlayer)
		return "Got " + splitCmd[1]
	case "whereami": // Duh
		return currentMap.CurrentRoom.Name
	default:
		return "Hmm..."
	}
}
