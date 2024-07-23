package commander

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
    currentMap *viper.Viper
)

func Init(currMap *viper.Viper) {
    currentMap := currMap
    fmt.Println(currentMap.GetString("0.name"))
}

func isGameCommand(cmd string) bool {
    return false
}

func GameCommand(cmd string) string {
    if ! isGameCommand(cmd) { return "Sorry, I didn't understand that..." }
    return ""
}
