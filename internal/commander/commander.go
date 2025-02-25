package commander

import (
	"explore/internal/maphandler"
	"explore/internal/playerhandler"
	"strings"
)

type Commander struct {
	InventoryChangeChannel chan struct{} // Made the channel a struct type cause it can have no data - i.e. serves as just a signal type

	Response        string
	ResponseChannel chan struct{}

    QuitChannel chan struct{}

	currentMap    *maphandler.MapInfo
	currentDB     *playerhandler.Database
	currentPlayer *playerhandler.Player
}

// Create a new commander
func Init(initMap *maphandler.MapInfo, initDB *playerhandler.Database, initPlayer *playerhandler.Player) *Commander {
	var tmpCommander Commander
	// Channels
	tmpCommander.InventoryChangeChannel = make(chan struct{}, 1) // Buffer capacity is set to 1
	tmpCommander.ResponseChannel = make(chan struct{}, 1)
    tmpCommander.QuitChannel = make(chan struct{}, 1)

	// Game Info
	tmpCommander.currentMap = initMap
	tmpCommander.currentDB = initDB
	tmpCommander.currentPlayer = initPlayer

    tmpCommander.currentMap.HardSetRoom(tmpCommander.currentPlayer.CurrentRoomIndex)

	return &tmpCommander
}

func (c *Commander) NotifyInvChange() {
	select {
	case c.InventoryChangeChannel <- struct{}{}: // When this function is called, an empty struct is sent into the channel - because we made the buffer 1 it succeeds if it can put the empty struct into it
	default:
		// InventoryChangeChannel is full, ignore
	}
}

func (c *Commander) NotifyResponse() {
	select {
	case c.ResponseChannel <- struct{}{}:
	default:
		// ResponseChannel full, ignore
	}
}

func (c *Commander) NotifyQuit() {
	select {
	case c.QuitChannel <- struct{}{}:
	default:
		// QuitChannel is full, ignore
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
func (c *Commander) PlayerCommand(cmd string) {
	trimmedCmd := strings.TrimSpace(cmd)
	cleanedCmd := strings.Fields(trimmedCmd)
	if  len(cleanedCmd) == 0 { // This is so dumb
		c.Response = "Hmm..."
		c.NotifyResponse()
		return
	}
    givenCommand := cleanedCmd[0]
    givenOptions := strings.Join(cleanedCmd[1:], " ")

	// Command prefix switch
	switch givenCommand {
	case "go": // Switch rooms
		if len(cleanedCmd) == 1 {
			c.Response = "Please specify a direction"
			c.NotifyResponse()
			return
		}
		if !c.currentMap.MoveDirection(givenOptions, c.currentPlayer.Inventory) {
			c.Response = "Could not move there"
			c.NotifyResponse()
			return
		}

        c.currentDB.SaveCurrentRoom(c.currentPlayer, c.currentMap.CurrentRoom.Index)

		c.Response = "Moved to " + c.currentMap.CurrentRoom.Name
        c.NotifyResponse()
	case "look": // Give us the look of the room - will remain the same - Maybe this should be called immediatly on entering...
		c.Response = c.currentMap.CurrentRoom.Look
		c.NotifyResponse()
	case "get": // Get the item in the room
		if len(cleanedCmd) == 1 {
			c.Response = "Please specify an item"
			c.NotifyResponse()
			return
		}
		if !c.currentMap.ItemInRoom(givenOptions) {
			c.Response = "That doesn't appear to be here."
			c.NotifyResponse()
			return
		}
		if c.currentPlayer.IsInInv(givenOptions) {
			c.Response = "You already have " + givenOptions
			c.NotifyResponse()
			return
		}
		// Save whatever item we get from the room
		c.currentPlayer.AddToInv(givenOptions)
		c.currentDB.SavePlayerInv(c.currentPlayer)

		// Notify UI of change - Add a signal to the channel
		c.NotifyInvChange()
		c.Response = "Got " + givenOptions 
		c.NotifyResponse()
	case "whereami": // Duh
		c.Response = c.currentMap.CurrentRoom.Name
		c.NotifyResponse()
    case "quit":
        c.NotifyQuit()
    case "escape":
        if c.currentMap.CurrentRoom.Index == c.currentMap.MetaData.EndRoomIndex {
            c.NotifyQuit()
        } else {
            c.Response = "You cannot escape from this room"
            c.NotifyResponse()
        }
	default:
		c.Response = "Hmm..."
		c.NotifyResponse()
	}
}
