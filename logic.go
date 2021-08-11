package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math"
)

var boardstart [][]int

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your Battlesnake username
		Color:      "#ACACE6", // TODO: Personalize
		Head:       "beluga", // TODO: Personalize
		Tail:       "default", // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)

  // map initial board layout to array
  boardstart = loadBoardIntoArray(state)
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
  var boardstate = boardstart // arrays are value types in go

  boardstate = avoidOrEatSnakes(state, boardstate)

  myHead := state.You.Body[0] // Coordinates of your head
	possibleMoves := map[string]int{
		"up":    boardstate[myHead.Y+1][myHead.X],
		"down":  boardstate[myHead.Y-1][myHead.X],
		"left":  boardstate[myHead.Y][myHead.X-1],
		"right": boardstate[myHead.Y][myHead.X+1],
	}
	
	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy
	var nextMove string

  // choose the move with the highest weight
  max := -5000
  for k, v := range possibleMoves {
    if v >= max {
      max = v
      nextMove = k
    }
  }

  log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)

	return BattlesnakeMoveResponse{
		Move: nextMove,
  }
}


// creates a matrix of the board with two extra rows and columns to represent the outside edges
// [0, 0] on the gamestate board is [1, 1] in the board array
func loadBoardIntoArray(state GameState) [][]int {
  var boardstate [][]int

  for y := 0; y < state.Board.Height + 1; y++ {
    for x := 0; x < state.Board.Width + 1; x++ {
      centreXSq := math.Pow(float64(x - state.Board.Width / 2), 2)
      centreYSq := math.Pow(float64(y - state.Board.Height / 2), 2)

      if x == 0 || y == 0 || x == state.Board.Width || y == state.Board.Height {
        // rule out going off the board
        boardstate[y][x] = -1000
      } else if centreXSq + centreYSq <= math.Pow(float64(state.Board.Width - 2), 2) {
        // steer within the circle
        if centreXSq + centreYSq < math.Pow(float64((state.Board.Width -2) / 2), 2) {
          // but away from the centre
          boardstate[y][x] = 5
        } else {        
          boardstate[y][x] = 10
        }
      } else {
        // avoid the corners
        boardstate[y][x] = -5
      }
    }
  }

  return boardstate
}



func avoidOrEatSnakes(state GameState, board [][]int) [][]int {
  
  for _, snake := range state.Board.Snakes {
    for i, coord := range snake.Body {
      if i == 0 && snake.Length < state.You.Length {
        // if our snake is longer, favour the head of the enemy snake
        board[coord.Y+1][coord.X+1] += 100
      } else {
        // avoid every other part of any snake
        board[coord.Y+1][coord.X+1] -= 1000
      }
    }
  }

  return board
}
