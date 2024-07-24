package commander

import (
    "strings"

    "github.com/spf13/viper"
)

var (
    currentMap  map[string]interface{}
    currentRoom map[string]interface{}
    directions = []string { "north", "east", "south", "west" }
)

func Init(currMap *viper.Viper) {
    err := currMap.Unmarshal(&currentMap)
    if err != nil {
        panic(err)
    }
    loadRoom("start")
}

func directionChecker(direction string) bool {
    for _, token := range directions {
        if strings.Contains(token, direction) { return true }
    }
    return false
}
    

func loadRoom(id string) string { 
    currentRoom = currentMap[id].(map[string]interface{}) 
    return "You have entered " + currentRoom["name"].(string)
}
func GetCurrentRoom() string { return currentRoom["name"].(string) }

func GameCommand(cmd string) string {
    splitCmd := strings.Split(cmd, " ")
    for _, token := range splitCmd {
        strings.ToLower(strings.ReplaceAll(token, " ", ""))
    }

    switch splitCmd[0] {
    case "go":
        if ! directionChecker(splitCmd[1]) { return "That is not a valid direction" }
        return loadRoom(currentRoom[splitCmd[1]].(string)) 
    case "look":
        return currentRoom["look"].(string)
    case "whereami":
        return GetCurrentRoom()
    default:
        return "Hmm..."
    }
}
