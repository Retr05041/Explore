package maphandler

import (
	"encoding/json"
	"os"
)

type metaData struct {
	StartRoomIndex int `json:"start"`
}

type room struct {
	Name       string  `json:"name"`
	NeededItem string  `json:"needed item"`
	North      *int    `json:"north"`
	East       *int    `json:"east"`
	South      *int    `json:"south"`
	West       *int    `json:"west"`
	Item       string `json:"item"`
	Look       string  `json:"look"`
	End        bool    `json:"end"`
}

type MapInfo struct {
	MetaData    metaData `json:"METADATA"`
	Rooms       []room   `json:"GAME"`
	CurrentRoom room
}

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

func validDirection(roomIndex *int) bool {
	if roomIndex != nil {
		return true
	}
	return false
}

func (m *MapInfo) MoveDirection(direction string) bool {
	switch direction {
	case "north":
		if !validDirection(m.CurrentRoom.North) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.North]
	case "east":
		if !validDirection(m.CurrentRoom.East) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.East]
	case "south":
		if !validDirection(m.CurrentRoom.South) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.South]
	case "west":
		if !validDirection(m.CurrentRoom.West) {
			return false
		}
		m.CurrentRoom = m.Rooms[*m.CurrentRoom.West]
	default:
		return false
	}
	return true
}

func (m *MapInfo) ItemInRoom(item string) bool {
	if m.CurrentRoom.Item == "nothing" { return false }
	if item != m.CurrentRoom.Item {
		return false	
	} 
	return true
}

func InitNewMap(filename string) (*MapInfo, error) {
	newMap, err := loadMap(filename)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}
