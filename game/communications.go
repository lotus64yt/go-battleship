package game

import "encoding/json"

type Packet struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ReadyPacket struct {
}

type AttackPacket struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type AttackResultPacket struct {
	Hit      bool `json:"hit"`
	Sunk     bool `json:"sunk"`
	GameOver bool `json:"game_over"`
}
