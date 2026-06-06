# Battleship

A terminal-based multiplayer Battleship game written in Go.

## Requirements

- Go 1.20 or higher (adjust according to your go.mod)

## Build and Run

To build for Linux:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o battleship main.go
```

To build for Windows:

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o battleship.exe main.go
```

To run the game:

```bash
./battleship
```

## Structure

- `app/`: Application lifecycle and game flow.
- `console/`: Terminal UI and input handling.
- `game/`: Core game logic and board state.
- `server/`: Network synchronization and multiplayer handling.
- `utils/`: Shared utilities.
