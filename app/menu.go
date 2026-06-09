package app

import (
	"fmt"
	"os"

	menu "github.com/octarahq/goclic"
	"github.com/octarahq/goclic/components"
)

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func Menu() {
	clearConsole()
	m := menu.NewMenu()
	var action func()

	m.Add(components.NewDisplay("BattleGO"))

	m.Add(components.NewButton("Create", components.WithButtonOnClick(func() {
		action = func() { CreateGame() }
		m.Stop()
	})))

	m.Add(components.NewButton("Connect", components.WithButtonOnClick(func() {
		action = MenuConnect
		m.Stop()
	})))

	m.Add(components.NewButton("Help", components.WithButtonOnClick(func() {
		action = MenuHelp
		m.Stop()
	})))

	m.Add(components.NewButton("Exit", components.WithButtonOnClick(func() {
		action = func() { os.Exit(0) }
		m.Stop()
	})))

	m.Start()

	if action != nil {
		action()
	}
}

func MenuConnect() {
	clearConsole()
	m := menu.NewMenu()
	var ip string
	var action func()

	m.Add(components.NewDisplay("Connect to a Game"))
	m.Add(components.NewInput("Enter URL or IP", &ip))

	m.Add(components.NewButton("Join", components.WithButtonOnClick(func() {
		if ip == "" {
			return
		}
		action = func() {
			clearConsole()
			JoinGame(ip)
		}
		m.Stop()
	})))

	m.Add(components.NewButton("Back", components.WithButtonOnClick(func() {
		action = Menu
		m.Stop()
	})))

	m.Start()

	if action != nil {
		action()
	}
}

func MenuHelp() {
	clearConsole()
	m := menu.NewMenu()
	var action func()

	m.Add(components.NewDisplay("--- HELP & INFOS ---"))
	m.Add(components.NewDisplay("Navigation: Use arrow keys to navigate."))
	m.Add(components.NewDisplay("Validation: Press Enter to confirm a choice or input."))
	m.Add(components.NewDisplay("--------------------"))
	m.Add(components.NewDisplay("BattleGO is a command-line battleship game."))
	m.Add(components.NewDisplay("You can create a local lobby or connect to a friend."))
	m.Add(components.NewDisplay("Using HTTPS proxies like devtunnels is supported."))

	m.Add(components.NewButton("Back", components.WithButtonOnClick(func() {
		action = Menu
		m.Stop()
	})))

	m.Start()

	if action != nil {
		action()
	}
}

