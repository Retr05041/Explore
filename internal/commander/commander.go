package commander

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

var (
    currentMap *viper.Viper
    commands = []string { "get", "look", "go" }
)

func Init(currMap *viper.Viper) {
    currentMap := currMap
    fmt.Println(currentMap.GetString("0.name"))
}

func isGameCommand(cmd string) bool {
    splitCmd := strings.Split(cmd, " ")
    if len(splitCmd) == 1 { return false }
    return slices.Contains(commands, splitCmd[0])
}

func GameCommand(cmd string) string {
    if ! isGameCommand(cmd) { return "Invalid command." }
    return "Valid command!"
}
