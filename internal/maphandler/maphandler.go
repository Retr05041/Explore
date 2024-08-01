package maphandler

import (
	"encoding/json"
	"os"
)

type metaData struct {
	StartRoomIndex int `json:"start"`
	EndRoomIndex   int `json:"end"`
}

type room struct {
	Index      int     `json:"index"`
	Name       string  `json:"name"`
	NeededItem *string `json:"needed item"`
	North      *int    `json:"north"`
	East       *int    `json:"east"`
	South      *int    `json:"south"`
	West       *int    `json:"west"`
	Item       *string `json:"item"`
	Look       string  `json:"look"`
}

type MapInfo struct {
	MetaData    metaData `json:"METADATA"`
	Rooms       []room   `json:"GAME"`
	CurrentRoom room
}

// Loads the map into the MapInfo struct
func loadMap(filename string) (*MapInfo, error) {
	tmpMap := new(MapInfo)
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &tmpMap); err != nil {
		return nil, err
	}

	tmpMap.CurrentRoom = tmpMap.Rooms[tmpMap.MetaData.StartRoomIndex]

	return tmpMap, nil
}

func (m *MapInfo) validDirection(roomIndex *int, playerInv []string) bool {
	if roomIndex != nil {
		if m.Rooms[*roomIndex].NeededItem != nil {
			for _, item := range playerInv {
				if *m.Rooms[*roomIndex].NeededItem == item {
					return true
				}
			}
			return false
		}
		return true
	}
	return false
}

func (m *MapInfo) MoveDirection(direction string, playerInv []string) bool {
	switch direction {
	case "north":
		if !m.validDirection(m.CurrentRoom.North, playerInv) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.North]
	case "east":
		if !m.validDirection(m.CurrentRoom.East, playerInv) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.East]
	case "south":
		if !m.validDirection(m.CurrentRoom.South, playerInv) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.South]
	case "west":
		if !m.validDirection(m.CurrentRoom.West, playerInv) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.West]
	default:
		return false
	}
	return true
}

// Checks if item is in the current room
func (m *MapInfo) ItemInRoom(item string) bool {
	if m.CurrentRoom.Item == nil {
		return false
	}
	if item != *m.CurrentRoom.Item {
		return false
	}
	return true
}

// Public load map function
func InitNewMap(filename string) (*MapInfo, error) {
	newMap, err := loadMap(filename)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}
