# Go BattleShip Server
Implements a server for turn based BattleShip clone that allows for multiple concurrent matches. Developed as class project to learn how applications maintain multiple connections with low latency and to learn the Go programming language. 

### Project Learning Goals: 
- Build a server capable of maintining multiple simultaneous connections independently and routing messages between clients.

- Host and manage a game of BattleShip with a complete state model consisting of game boards, ship placement, turn order, and win condition detection.

### Features: 
- Websocket connection with per-client input/output routines and teardown on game end or disconnection.
- Central hub tracking client connections, matchmaking (FIFO), and spawning instances of matches.
- Turn-based match loop handling input from clients and tracking game state.
- Game board & ship hit detection.

### Known Issues: 
- No client application for users to play the game. All of the testing was done via command line with JSON formatted messages.
- Duplicate code for Player branches in match loop.
- Errors logged but not communicated to user.
- No win/lose messaging to players at the conclusion of a match. 
- After a match ends players are disconnected without an option to play again and no messaging sent. 
- Ships are placed randomly without input from players.
- Match can’t be played without two clients connected. 
- Unused fields in structs.
- Code quality. 


## Run
#### Prerequisites 
- websocat to connect and send JSON messages

#### 1. Clone the repository
```bash
git clone https://github.com/sdonlin5/go-battleship-server.git
```

#### 2. Navigate to the repository
```bash
cd go-battleship-server
```

#### 3. Start the server
```bash
# from the project root
go build
./main

```
**Or**

```bash
# from the project root
go run main.go
```

This starts the hub and listens on port `8080` with the websocket endpoint at `/ws`. 

Leave it running in its own terminal — you'll watch its log output as you play.


#### 4. Connect clients
The hub spawns a match when only when **two** clients are queued. You will need to connect to the server from two separate console terminal windows. 

**Terminal A**

```bash
websocat ws://localhost:8080/ws
```
**Terminal B**

```bash
websocat ws://localhost:8080/ws
```

When each client connects to the server, the console will display:

```bash
Client Spawned: CLIENT_ID
Client Registered: ID = CLIENT_ID
Client Queued: ID = CLIENT_ID
```

The first player to connect will be `Player1` and go first, while the second is `Player2` and goes second. 

Once both clients are queued the console will output: 

```bash
Match Spawned: ID = MATCH_ID
Start Match
```

Once the match starts each player will be notified: 

**Player1 terminal**
```bash
{"type":"firstTurn","payload":null}
```

**Player2 terminal**
```
{"type":"secondTurn","payload":null}
```

#### Messaging

```json
{"type": "<Type>", "payload": <Payload>}
```

**Inputs**

Every message in each direction uses the same format: 
`{“type”: “<Type>”, “payload:” <PayLoad>}`

The input pump only recognizes: 
| type         | payload                   | notes                  |
| ------------ | ------------------------- | ---------------------- |
| `Shot`       | `{“x": <0-9>,”y”: <0-9>}` | fires at x,y           |
| `PlayerQuit` | `{}`                      | ends match immediately |

**Outputs**
Response types are defined `game_match.go`


| type                     | payload                                                        | when                                                      |
| ------------------------ | -------------------------------------------------------------- |:--------------------------------------------------------- |
| `firstTurn`/`secondTurn` | `null`                                                         | right after match spawns                                  |
| `yourTurn`               | `null`                                                         | sent to current player after swicth                       |
| `ShotResult`             | `{"shot":{"x":..,"y":..},"valid":true,"hit":bool,"sunk":bool}` | broadcast to both players after a shot is played          |
| `OutOfTurnResult`        | `null`                                                         | sent back to player if shot played out of turn            |
| `InvalidShotResult`      | `null`                                                         | sent back to player if invalid shot coordinates sent      |
| `PlayerWin`              | `null`                                                         | broadcast to both players once a sides fleet is destroyed |
| `PlayerQuitResult`       | `null`                                                         | broadcast to both players if “Quit” is played. Match ends |
| `Disconnect`             | `null`                                                         | sent to remaining player if opponent disconnects          |


#### Play a turn 

Player1 goes first:

```json
{"type":"Shot","payload":{"x":0,"y":0}}
```

Both terminals print the result:

```json
{"type":"ShotResult","payload":{"shot":{"x":0,"y":0},"valid":true,"hit":false,"sunk":false}}
```

Server swaps players and Player2’s terminal outputs: 

```bash
{"type":"yourTurn","payload":null}
```
After each valid shot, the turn alternates and the result is broadcast to each player. 

If you play a a shot from the terminal that *doesn’t* currently hold the turn, you’ll get `OutOfTurnResult` back (only in that terminal — the other player isn't
notified).

Fleets are placed randomly and hidden and not displayed to the players. If a shot hits an opponent's ship, the `ShotResult` `hit` field will be: `hit:true`. If a ship’s health is fully depleted, meaning every tile it occupied was hit, The `sunk` field in `ShotResult` will be `true`. 

#### End the match early
Each terminal can end the match early at any time by playing 

```json
{"type": "PlayerQuit", "payload" {}}
```

#### Notes
- Additional clients can connect via `websocat` but a match will only be spawned once there are two clients queued.
- `x` and `y` shot inputs are `uint8` type, sending a negative value will result in an error
- Closing the `websocat` connection with `ctrl+C` is equivalent to sending `PlayerQuit`
