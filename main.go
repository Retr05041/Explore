package main

import (
	"fmt"
    "explore/internal/commander"

	"github.com/spf13/viper"
)


func main() {
    // Menu
    // Get menu option

    // Load specific map / player info
    // Viper / sqlite?
    currentMap := viper.New()

    // This is just a placeholder for testing... menu needs to be implemented
    currentMap.SetConfigName("prologue")
    currentMap.AddConfigPath("./maps/")
    err := currentMap.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error map file: %w", err))
    }

    // Init commander - using specified map and player info --- Maps should never be changed, their state is represented by the players data
    commander.Init(currentMap)

    // Run tui once commander is setup
    // tui.Start()
}
