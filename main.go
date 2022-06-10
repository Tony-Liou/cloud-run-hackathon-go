package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"

	"github.com/liyue201/gostl/ds/queue"
)

var myName = "https://cloud-run-hackathon-go-n5xjuxfciq-uc.a.run.app"

type Position struct {
	x, y int
}

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	d := json.NewDecoder(req.Body)
	defer req.Body.Close()
	d.DisallowUnknownFields()

	var v ArenaUpdate
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if myName != v.Links.Self.Href {
		myName = v.Links.Self.Href
	}

	resp := play(&v)
	fmt.Fprint(w, resp)
}

func play(input *ArenaUpdate) string {
	log.Printf("IN: %#v", input)

	board, myPos := getBoard(input.Arena.State, input.Arena.Dimensions[0], input.Arena.Dimensions[1])
	targetName := findAttackableEnemy(board, myPos)
	// Not found
	if targetName == "" {
		enemyPos := findNearestEnemy(board, myPos)
		targetName = board[enemyPos.y][enemyPos.x]
	}
	action := takeAction(myName, targetName, input.Arena.State)

	return action.String()
}

// getBoard creates a two dimension string slice. Return the current board and my position.
func getBoard(playerInfo map[string]PlayerState, x, y int) (board [][]string, myself Position) {
	board = make([][]string, y)
	for i := range board {
		board[i] = make([]string, x)
	}

	for name, state := range playerInfo {
		if name != myName {
			board[state.Y][state.X] = name
		} else {
			myself.x = state.X
			myself.y = state.Y
		}
	}

	return
}

// findAttackableEnemy finds if an enemy is in the same line with me.
func findAttackableEnemy(board [][]string, myself Position) string {
	for i := 1; i <= 3; i++ {
		y := myself.y - i
		x := myself.x
		if isInside(x, y, board) && board[y][x] != "" {
			return board[y][x]
		}

		y = myself.y + i
		if isInside(x, y, board) && board[y][x] != "" {
			return board[y][x]
		}

		y = myself.y
		x = myself.x - i
		if isInside(x, y, board) && board[y][x] != "" {
			return board[y][x]
		}

		x = myself.x + i
		if isInside(x, y, board) && board[y][x] != "" {
			return board[y][x]
		}
	}

	return ""
}

// isInside checks the position is in the board or not.
func isInside(x, y int, board [][]string) bool {
	row := len(board)
	col := len(board[0])

	return x >= 0 && x < col && y >= 0 && y < row
}

// takeAction uses the facing direction and location to decide how to move.
func takeAction(attackerName, targetName string, playerInfo map[string]PlayerState) Action {
	attacker := playerInfo[attackerName]
	target := playerInfo[targetName]

	// Same column
	if attacker.X-target.X == 0 {
		// Below target
		if attacker.Y-target.Y > 0 {
			switch attacker.Direction {
			case "N":
				return Throw
			case "E":
				return TurnLeft
			case "S":
				return TurnRight
			case "W":
				return TurnRight
			}
		} else { // Beyond target
			switch attacker.Direction {
			case "N":
				return TurnRight
			case "E":
				return TurnRight
			case "S":
				return Throw
			case "W":
				return TurnLeft
			}
		}
	} else if attacker.Y-target.Y == 0 { // Same row
		if attacker.X-target.X > 0 { // Right
			switch attacker.Direction {
			case "N":
				return TurnLeft
			case "E":
				return TurnLeft
			case "S":
				return TurnRight
			case "W":
				return Throw
			}
		} else { // Left
			switch attacker.Direction {
			case "N":
				return TurnRight
			case "E":
				return Throw
			case "S":
				return TurnLeft
			case "W":
				return TurnLeft
			}
		}
	} else { // Not in line
		if attacker.X-target.X > 0 {
			if attacker.Y-target.Y > 0 { // Bottom right
				switch attacker.Direction {
				case "N":
					return MoveForward
				case "E":
					return TurnLeft
				case "S":
					return TurnRight
				case "W":
					return MoveForward
				}
			} else { // Top right
				switch attacker.Direction {
				case "N":
					return TurnLeft
				case "E":
					return TurnRight
				case "S":
					return MoveForward
				case "W":
					return MoveForward
				}
			}
		} else {
			if attacker.Y-target.Y > 0 { // Bottom left
				switch attacker.Direction {
				case "N":
					return MoveForward
				case "E":
					return MoveForward
				case "S":
					return TurnLeft
				case "W":
					return TurnRight
				}
			} else { // Top left
				switch attacker.Direction {
				case "N":
					return TurnRight
				case "E":
					return MoveForward
				case "S":
					return MoveForward
				case "W":
					return TurnLeft
				}
			}
		}
	}

	rand := rand2.Intn(4)
	return Action(rand)
}

// findNearestEnemy uses BFS to find an enemy's position.
func findNearestEnemy(board [][]string, myself Position) Position {
	dir := [4][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	visited := make([][]bool, len(board))
	for i := range visited {
		visited[i] = make([]bool, len(board[0]))
	}

	q := queue.New()
	q.Push(myself)
	visited[myself.y][myself.x] = true

	for !q.Empty() {
		cur := q.Pop().(Position)

		if board[cur.y][cur.x] != "" {
			return cur
		}

		for i := 0; i < 4; i++ {
			x := cur.x + dir[i][0]
			y := cur.y + dir[i][1]

			if isInside(x, y, board) && !visited[y][x] {
				q.Push(Position{x, y})
				visited[y][x] = true
			}
		}
	}

	return Position{}
}
