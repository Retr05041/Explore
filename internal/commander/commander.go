package commander

import (
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"strings"
)

type Commander struct {
	InventoryChangeChannel chan struct{} // Made the channel a struct type cause it can have no data - i.e. serves as just a signal type

	currentMap    *maphandler.MapInfo
	currentDB     *playerhandler.Database
	currentPlayer *playerhandler.Player
}

// Create a new commander
func Init(initMap *maphandler.MapInfo, initDB *playerhandler.Database, initPlayer *playerhandler.Player) *Commander {
	var tmpCommander Commander
    // Channels
    tmpCommander.InventoryChangeChannel = make(chan struct{}, 1) // Buffer capacity is set to 1

	// Game Info
	tmpCommander.currentMap = initMap
	tmpCommander.currentDB = initDB
	tmpCommander.currentPlayer = initPlayer

	return &tmpCommander
}

func (c *Commander) NotifyInvChange() {
    select {
    case c.InventoryChangeChannel <- struct{}{}: // When this function is called, an empty struct is sent into the channel - because we made the buffer 1 it succeeds if it can put the empty struct into it
    default:
        // InventoryChangeChannel is full, ignore
    }
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
		if len(cleanedCmd) == 1 {
			return "Please specify a direction"
		}
		if !c.currentMap.MoveDirection(cleanedCmd[1], c.currentPlayer.Inventory) {
			return "Could not move there"
		}
		return "Moved to " + c.currentMap.CurrentRoom.Name
	case "look": // Give us the look of the room - will remain the same - Maybe this should be called immediatly on entering...
		return c.currentMap.CurrentRoom.Look
	case "get": // Get the item in the room
		if len(cleanedCmd) == 1 {
			return "Please specify an item"
		}
		if !c.currentMap.ItemInRoom(cleanedCmd[1]) {
			return "That doesn't appear to be here."
		}
		if c.currentPlayer.IsInInv(cleanedCmd[1]) {
			return "You already have " + cleanedCmd[1]
		}
		// Save whatever item we get from the room
		c.currentPlayer.AddToInv(cleanedCmd[1])
		c.currentDB.SavePlayerInfo(c.currentPlayer)

        // Notify UI of change - Add a signal to the channel
        c.NotifyInvChange()
		return "Got " + cleanedCmd[1]
	case "whereami": // Duh
		return c.currentMap.CurrentRoom.Name
	default:
		return "Hmm..."
	}
}
