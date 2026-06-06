package game

import (
	utilsboard "battleship/utils/board"
	"encoding/json"
	"fmt"
	"net"
)

type Packet struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ReadyPacket struct {
}

type AttackPacket struct {
	Notation string `json:"notation"`
}

type AttackResultPacket struct {
	Hit      bool `json:"hit"`
	Sunk     bool `json:"sunk"`
	GameOver bool `json:"game_over"`
}

func SendReadyPacket(conn net.Conn) error {
	packet := Packet{
		Type: "ready",
		Data: json.RawMessage(`{}`),
	}

	return json.NewEncoder(conn).Encode(packet)
}

func SendAttackPacket(conn net.Conn, notation string, incoming <-chan Packet) (bool, bool, bool, error) {
	correct, err := utilsboard.IsValidNotation(notation)
	if !correct && err == nil {
		return false, false, false, fmt.Errorf("Invalid notation")
	} else if !correct && err != nil {
		return false, false, false, fmt.Errorf("Invalid notation : %s", err)
	}
	data, err := json.Marshal(AttackPacket{Notation: notation})
	if err != nil {
		return false, false, false, err
	}

	if err := json.NewEncoder(conn).Encode(Packet{
		Type: "attack",
		Data: data,
	}); err != nil {
		return false, false, false, err
	}

	packet, ok := <-incoming
	if !ok {
		return false, false, false, fmt.Errorf("connection closed")
	}

	if packet.Type != "attack_result" {
		return false, false, false, fmt.Errorf("unexpected packet type: got %s, want attack_result", packet.Type)
	}

	var result AttackResultPacket
	if err := json.Unmarshal(packet.Data, &result); err != nil {
		return false, false, false, err
	}

	return result.Hit, result.Sunk, result.GameOver, nil
}

func WaitNextPacket(conn net.Conn) (*Packet, error) {
	var packet Packet

	err := json.NewDecoder(conn).Decode(&packet)
	if err != nil {
		return nil, err
	}

	return &packet, nil
}

func WaitNextPacketType(conn net.Conn, Type string) (*Packet, error) {
	packet, err := WaitNextPacket(conn)
	if err != nil {
		return nil, err
	}

	if packet.Type != Type {
		return nil, fmt.Errorf("unexpected packet type: got %s, want %s", packet.Type, Type)
	}

	return packet, nil
}
