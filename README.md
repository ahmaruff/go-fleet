```
 $$$$$$\                   $$$$$$$$\ $$\                     $$\     
$$  __$$\                  $$  _____|$$ |                    $$ |    
$$ /  \__| $$$$$$\         $$ |      $$ | $$$$$$\   $$$$$$\$$$$$$\   
$$ |$$$$\ $$  __$$\$$$$$$\ $$$$$\    $$ |$$  __$$\ $$  __$$\_$$  _|  
$$ |\_$$ |$$ /  $$ \______|$$  __|   $$ |$$$$$$$$ |$$$$$$$$ |$$ |    
$$ |  $$ |$$ |  $$ |       $$ |      $$ |$$   ____|$$   ____|$$ |$$\ 
\$$$$$$  |\$$$$$$  |       $$ |      $$ |\$$$$$$$\ \$$$$$$$\ \$$$$  |
 \______/  \______/        \__|      \__| \_______| \_______| \____/ `
```

A terminal-based multiplayer Battleship game built with Go, featuring real-time ASCII graphics and TCP networking.

## Features

- **Multiplayer**: Real-time 1v1 gameplay over TCP
- **Visual**: Beautiful ASCII game boards with live updates
- **Simple Commands**: Easy-to-use command interface
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **No Dependencies**: Uses only Go standard library

## Game Rules

- **Grid**: 10x10 battlefield (A1 to J10)
- **Ships**: Each player places 5 single-cell ships
- **Turns**: Players take turns firing at coordinates
- **Win Condition**: Destroy all enemy ships to win

## Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/ahmaruff/go-fleet
cd go-fleet
go mod tidy
go build -o server cmd/server/main.go
go build -o client cmd/client/main.go
```

### 2. Start Server
```bash
./server --port 8080
```

### 3. Connect Players
**Terminal 1:**
```bash
./client --host localhost --port 8080
```

**Terminal 2:**
```bash
./client --host localhost --port 8080
```

### 4. Play the Game
1. Enter your name when prompted
2. Type `/ready` to join matchmaking
3. Place ships: `/set A1`, `/set B2`, etc.
4. Fire at opponent: `/fire C3`, `/fire D4`, etc.

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/name <name>` | Set your player name | `/name Alice` |
| `/ready` | Join matchmaking queue | `/ready` |
| `/set <coord>` | Place ship at coordinate | `/set A1` |
| `/fire <coord>` | Fire at enemy coordinate | `/fire B3` |
| `/quit` | Exit the game | `/quit` |

## Game Flow

```
1. Players connect and set names
2. Both players send /ready → matched automatically
3. Ship Placement Phase:
   - Each player places 5 ships using /set
   - Real-time board updates
4. Combat Phase:
   - Players fire using /fire commands
   - See hit/miss results instantly
5. Victory:
   - First to destroy all enemy ships wins
```

## Project Structure

```
go-fleet/
├── cmd/
│   ├── server/
│   │   └── main.go         # Game server handler
│   ├── client/
│   │   └── main.go         # Game client handler
│   └── test/
│       └── main.go         # Simple e2e test
├── internal/
│   ├── game/               # Core game logic
│   │   ├── board.go        # Game board and ship management
│   │   ├── game.go         # Game state and flow control
│   │   ├── player.go       # Player data structure
│   │   └── coordinate.go   # Coordinate conversion
│   ├── display/
│   │   └── display.go      # Game UI rendering
│   └── effects/
│       └── effects.go      # ASCII Art Effect
├── .gitignore
├── go.mod
├── LICENSE
└── README.md
```

## Architecture

- **Server**: Manages multiple games, handles matchmaking, coordinates turns, sends effect game state to client
- **Client**: Connects to server, sends commands, displays game state
- **Game Logic**: Pure game rules independent of networking
- **Display System**: Game ASCII rendering with real-time updates
- **Effects**: ASCII Art effect for each game state

## Example Gameplay

```
========================================== GO-FLEET ==========================================
Player: Alice vs Bob | Phase: PLAYING | Current Turn: Player 1
===========================================================================================
Your Board:                               Opponent's Board:
  A B C D E F G H I J                       A B C D E F G H I J
 1 S ~ ~ ~ ~ ~ ~ ~ ~ ~                      1 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 2 ~ S ~ ~ ~ ~ ~ ~ ~ ~                      2 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 3 ~ ~ S ~ ~ ~ ~ ~ ~ ~                      3 ~ ~ X ~ ~ ~ ~ ~ ~ ~
 4 ~ ~ ~ S ~ ~ ~ ~ ~ ~                      4 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 5 ~ ~ ~ ~ S ~ ~ ~ ~ ~                      5 O ~ ~ ~ ~ ~ ~ ~ ~ ~
 6 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~                      6 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 7 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~                      7 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 8 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~                      8 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
 9 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~                      9 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~
10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~                     10 ~ ~ ~ ~ ~ ~ ~ ~ ~ ~

Commands: /ready, /set A1, /fire B2
```

## Symbol Legend

| Symbol | Meaning |
|--------|---------|
| `~` | Water / Unknown |
| `S` | Your ship (not hit) |
| `X` | Hit (ship destroyed) |
| `O` | Miss (water hit) |


## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

