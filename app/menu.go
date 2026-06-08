package app

import (
	menu "github.com/octarahq/goclic"
	"github.com/octarahq/goclic/components"
)

func Menu() {
	m := menu.NewMenu()
	var ip string
	var action func()

	m.Add(components.NewDisplay("BattleGO"))
	m.Add(components.NewButton("Create a lobby", components.WithButtonOnClick(func() {
		action = func() { CreateGame() }
		m.Stop()
	})))
	m.Add(components.NewInput("Host IP:Port (Required)", &ip))
	m.Add(components.NewButton("Join", components.WithButtonOnClick(func() {
		if ip == "" {
			return
		}
		action = func() { JoinGame(ip) }
		m.Stop()
	})))
	m.Start()

	if action != nil {
		action()
	}
}
