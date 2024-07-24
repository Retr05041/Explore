package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Room struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	NeededItem string  `json:"needed item"`
	North      *string `json:"north"`
	East       *string `json:"east"`
	South      *string `json:"south"`
	West       *string `json:"west"`
	Item       *string `json:"item"`
	Look       string  `json:"look"`
	End        bool    `json:"end"`
}

func parseRoom(filename string, roomID string) (*Room, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

    // read open bracket
	t, err := decoder.Token()
	if err != nil {
        return nil, err
	}
	fmt.Printf("%T: %v\n", t, t)

	// Loop until we find the desired room or reach the end of the JSON array
	for decoder.More() {
		var room Room

		// Decode JSON object from the stream
		if err := decoder.Decode(&room); err != nil {
			return nil, err
		}

		// Check if this is the desired room
		if room.ID == roomID {
			return &room, nil
		}
	}

	// If roomID not found, return an error
	return nil, fmt.Errorf("room '%s' not found in JSON", roomID)
}

func main() {
	filename := "prologue.json"
	roomID := "start"

	room, err := parseRoom(filename, roomID)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Room Name:", room.Name)
	fmt.Println("Needed Item:", room.NeededItem)
	fmt.Println("North:", room.North)
	fmt.Println("East:", room.East)
	fmt.Println("South:", room.South)
	fmt.Println("West:", room.West)
	fmt.Println("Item:", room.Item)
	fmt.Println("Look:", room.Look)
	fmt.Println("End:", room.End)
}
