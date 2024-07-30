package commander

import (
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"strings"
)

type Commander struct {
	currentMap    *maphandler.MapInfo
	currentDB     *playerhandler.Database
	currentPlayer *playerhandler.Player
}

// Give the commander the map, db, and player
func Init(initMap *maphandler.MapInfo, initDB *playerhandler.Database, initPlayer *playerhandler.Player) *Commander {
    var tmpCommander Commander
	// Initialise map
	tmpCommander.currentMap = initMap

	// Load / Create the database for the map
	tmpCommander.currentDB = initDB

	// Load the player from the database
	tmpCommander.currentPlayer = initPlayer

    return &tmpCommander
}

// Get players inventory -- for displaying
func (c *Commander) GetCurrPlayerInv() []string {
	return c.currentPlayer.Inventory
}

// Get players name -- for messages
func (c *Commander) GetCurrPlayerName() string {
	return c.currentPlayer.Name
}

// When the player gives a command, this handles it and returns a string to be shown to the player in response
func (c *Commander) PlayerCommand(cmd string) string {
    trimmedCmd := strings.TrimSpace(cmd)
    cleanedCmd := strings.Fields(trimmedCmd)
	if len(cleanedCmd) > 2 { // This is so dumb
		return "Hmm..."
	}

    // Command prefix switch
	switch cleanedCmd[0] {
    case "go": // Switch rooms
        if len(cleanedCmd) == 1 { return "Please specify a direction" }
		if !c.currentMap.MoveDirection(cleanedCmd[1], c.currentPlayer.Inventory) {
			return "Could not move there"
		}
		return "Moved to " + c.currentMap.CurrentRoom.Name
	case "look": // Give us the look of the room - will remain the same - Maybe this should be called immediatly on entering...
		return c.currentMap.CurrentRoom.Look
	case "get": // Get the item in the room
        if len(cleanedCmd) == 1 { return "Please specify an item" }
		if !c.currentMap.ItemInRoom(cleanedCmd[1]) {
			return "That doesn't appear to be here."
		}
		if c.currentPlayer.IsInInv(cleanedCmd[1]) {
			return "You already have " + cleanedCmd[1]
		}
        // Save whatever item we get from the room
		c.currentPlayer.AddToInv(cleanedCmd[1])
		c.currentDB.SavePlayerInfo(c.currentPlayer)
		return "Got " + cleanedCmd[1]
	case "whereami": // Duh
		return c.currentMap.CurrentRoom.Name
	default:
		return "Hmm..."
	}
}
