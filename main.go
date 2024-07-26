package main

import (
	//"explore/internal/commander"
	//"explore/internal/tui"
    "log"
    "fmt"
    "explore/internal/playerhandler"
)


func main() {
    // Menu
    // Get menu option
    test := playerhandler.LoadDatabase("prologue")
    test.CreatePlayer("testplayer")
    currentPlayer, err := test.LoadPlayer("testplayer")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(currentPlayer.Name)

    //commander.Init("./maps/prologue.json") // Run tui once commander is setup

    //tui.Start()
}
