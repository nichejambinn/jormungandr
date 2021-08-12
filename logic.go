package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
  "fmt"
	"log"
	"math"
)

var lastMove string

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "nichejambinn",
		Color:      "#BEE7B8",
		Head:       "evil",
		Tail:       "fat-rattle",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
  boardstate := makeBoard(state.Board.Height + 2, state.Board.Width + 2)
  
  loadBoardIntoArray(state, boardstate)
  avoidOrEatSnakes(state, boardstate)
  eatWhenHungry(state, boardstate)
  //steerToCenter(state, boardstate)

  // Get coords of head where [0,0] -> [1,1]
  myHead := state.You.Body[0]
  myHead.Y += 1
  myHead.X += 1

  // TODO: Stay as spread out as possible
  // TODO: Favour a path toward nearby weak snakes?
  // TODO: Try and pin snakes into corners

	// Finally, choose a move from the available safe moves.
  displayBoard(boardstate)

	// Select a move to make based on strategy
	possibleMoves := map[string]int{
		"up":    boardstate[myHead.Y+1][myHead.X],
		"down":  boardstate[myHead.Y-1][myHead.X],
		"left":  boardstate[myHead.Y][myHead.X-1],
		"right": boardstate[myHead.Y][myHead.X+1],
	}

	var nextMove string

  //choose the move with the highest weight
  max := -5000
  for k, v := range possibleMoves {
    if v >= max {
      max = v
      nextMove = k
    }
  }

  lastMove = nextMove

  log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)

	return BattlesnakeMoveResponse{
		Move: nextMove,
  }
}


// creates a matrix of the board with two extra rows and columns to represent the outside edges
// [0, 0] on the gamestate board is [1, 1] in the board array
func loadBoardIntoArray(state GameState, boardstate [][]int) {
  height := len(boardstate)
  width := len(boardstate[0])

  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      centerXSq := math.Pow(float64(x - (width / 2)), 2)
      centerYSq := math.Pow(float64(y - (height / 2)), 2)

      if x == 0 || y == 0 || x == width - 1 || y == height - 1 {
        // rule out going off the board
        boardstate[y][x] = -1000
      } else if centerXSq + centerYSq <= math.Pow(float64(state.Board.Width / 2 - 1), 2) {
        // steer within the circle
        if centerXSq + centerYSq < math.Pow(float64(state.Board.Width / 4), 2) {
          // but away from the centre
          boardstate[y][x] = 100
        } else {        
          boardstate[y][x] = 80
        }
      } else {
        // avoid the corners
        boardstate[y][x] = -50
      }
    }
  }
}



func avoidOrEatSnakes(state GameState, boardstate [][]int) {
  
  for _, snake := range state.Board.Snakes {
    for i, coord := range snake.Body {
      if i == 0 && snake.Length < state.You.Length {
        // if our snake is longer, favour the head of the enemy snake
        bloom(Coord{coord.X+1, coord.Y+1}, 100, 3, boardstate)
        //boardstate[coord.Y+1][coord.X+1] += 100
      } else {
        if i != 0 || snake.Name != state.You.Name {
          // avoid every other part of any snake
          boardstate[coord.Y+1][coord.X+1] -= 1000
          bloom(Coord{coord.X+1, coord.Y+1}, -150 / (i + 1), 2, boardstate)
        }
      }
    }
  }
}




func steerToCenter(state GameState, boardstate [][]int) {
  myHead := state.You.Body[0]
  myHead.Y += 1
  myHead.X += 1

	// Select a move to make based on strategy
	possibleMoves := map[string]Coord{
		"up":    Coord{myHead.X, myHead.Y+1},
		"down":  Coord{myHead.X, myHead.Y-1},
		"left":  Coord{myHead.X-1, myHead.Y},
		"right": Coord{myHead.X+1, myHead.Y},
  }

  var dir string
  min := 1000
  for k, coord := range possibleMoves {
    if sqDistance(coord, myHead) < min {
      dir = k
    }
  }

  c := possibleMoves[dir]
  boardstate[c.Y][c.X] += 30
}



// weigh the area spread around a coord with diminishing power
// a spread of 1 captures the area within radius 1 of the center
func bloom(center Coord, power int, spread int, boardstate [][]int) {
  height := len(boardstate)
  width := len(boardstate[0])
  for y := center.Y - spread; y < center.Y + spread; y++ {
    for x := center.X - spread; x < center.X + spread; x++ {
      if x >= 0 && y >= 0 && x < width && y < height {
        centerXSq := math.Pow(float64(x - center.X), 2)
        centerYSq := math.Pow(float64(y - center.Y), 2)

        for r := 0; r <= spread; r++ {
          if centerXSq + centerYSq <= math.Pow(float64(r), 2) {
            boardstate[y][x] += power / (r+1)
          }
        }
      }
    }
  }
}



func eatWhenHungry(state GameState, boardstate [][]int) {
  var maxLength int32
  for _, snake := range state.Board.Snakes {
    if snake.Length > maxLength {
      maxLength = snake.Length
    }
  }

  isHungry := (state.You.Length < maxLength + int32(len(state.Board.Snakes))) || (state.You.Health < 50)

  if isHungry {
    for _, food := range state.Board.Food {
      bloom(Coord{food.X+1, food.Y+1}, 80, 3, boardstate)
    }
  }
}



func sqDistance(coord1 Coord, coord2 Coord) int {
  return int(math.Pow(float64(coord2.X - coord1.X), 2) + math.Pow(float64(coord2.Y - coord1.Y), 2))
}



func displayBoard(board [][]int) string {
  view := "Displaying board \n"

  height := len(board)
  width := len(board[0])

  for y := height - 1; y >= 0; y-- {
    for x := 0; x < width; x++ {
      view += fmt.Sprintf("%5d|", board[y][x])
    }
    view += "\n"
  }
  view += "\n"

  log.Printf(view)

  return view
}




func makeBoard(dy int, dx int) [][]int {
    a := make([][]int, dy)
    for i := range a {
        a[i] = make([]int, dx)
    }

    return a
}