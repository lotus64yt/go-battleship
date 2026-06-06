package console

import (
	"github.com/turret-io/go-menu/menu"
)

func Menu() {
	commandOptions := []menu.CommandOption{
		menu.CommandOption{"create", "Create a game", CreateGame},
		// menu.CommandOption{"join", "Join a game", JoinGame},
		menu.CommandOption{"stop", "Exit battleship", Exit},
	}

	menuOptions := menu.NewMenuOptions("'menu' for help > ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
