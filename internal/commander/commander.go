package commander

import (
	"strings"
    "explore/internal/maphandler"
)

var (
    gameMap maphandler.MapInfo
)

func Init(mapLocation string) error {
    holdMap, err := maphandler.InitNewMap(mapLocation) 
    if err != nil { return err }
    gameMap = *holdMap
    return nil
}
    
func GameCommand(cmd string) string {
    splitCmd := strings.Split(cmd, " ")
    if len(splitCmd) > 2 { return "Hmm..." }

    for _, token := range splitCmd {
        strings.ToLower(strings.ReplaceAll(token, " ", ""))
        if token == "" { splitCmd[len(splitCmd)-1] = "UNKNOWN" } // Bogus check to see if there is a whitespace element - needs to be after the length check to work :(
    }

    switch splitCmd[0] {
    case "go":
        if ! gameMap.MoveDirection(splitCmd[1]) { return "Could not move there" }
        return "Moved to " + gameMap.CurrentRoom.Name
    case "look":
        return gameMap.CurrentRoom.Look
    case "whereami":
        return gameMap.CurrentRoom.Name
    default:
        return "Hmm..."
    }
}
