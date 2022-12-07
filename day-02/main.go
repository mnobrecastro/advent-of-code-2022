// Miguel Nobre Castro
// https://adventofcode.com/2022/day/2

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Game of Rock-Paper-Scissors
type Game struct {
	game    chan string
	rounds  int
	player1 rune
	player2 rune
	score1  int
	score2  int
}

// Initiate a NewGame given an input .txt file
//
// Input file example:
// "
// A Y
// B X
// C Z
// "
func NewGame(filename string) (g *Game) {

	game := make(chan string, 1)
	buffer := ""
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(game)
		}

		r := bufio.NewReader(f)
		for true {
			s, err := r.ReadString('\n')
			if !errors.Is(err, io.EOF) {
				s = s[:len(s)-1]
				if len(s) > 0 {
					buffer += s
				} else {
					continue
				}
			} else {
				close(game)
				break
			}
			game <- buffer
			buffer = ""
		}
		f.Close()
	}()

	g = &Game{
		game:    game,
		rounds:  0,
		player1: '_',
		player2: '_',
		score1:  0,
		score2:  0,
	}
	return
}

// Play a Game of Rock-Paper-Scissors
func (g *Game) Play(strategy func(rune, rune) (int, int)) {

	for s := range g.game {
		signs := []rune(s)
		g.player1 = signs[0]
		g.player2 = signs[2]

		p1, p2 := strategy(g.player1, g.player2)

		if Abs(p1) >= p2 {
			fmt.Println("Player1 takes the round...")
		} else if Abs(p1) == p2 {
			fmt.Println("It's a draw...")
		} else {
			fmt.Println("Player2 takes the round...")
		}

		g.score1 += p1
		g.score2 += p2
		g.rounds += 1
	}
	fmt.Printf("!!!GAME OVER!!!\nFinal scores after %d rounds:\n", g.rounds)
	fmt.Printf("Player1 - %d vs %d - Player 2\n", Abs(g.score1), g.score2)
	if Abs(g.score1) >= g.score2 {
		fmt.Println("Player1 WINS!")
	} else if Abs(g.score1) == g.score2 {
		fmt.Println("DRAW! Play again!")
	} else {
		fmt.Println("Player2 WINS!")
	}
	return
}

// Hand sign mapping
func pts(r rune) (p int) {

	keys := map[rune]int{
		'A': -1, // Player1 Rock
		'B': -2, // Player1 Paper
		'C': -3, // Player1 Scissors
		'X': 1,  // Player2 Rock
		'Y': 2,  // Player2 Paper
		'Z': 3,  // Player2 Scissors
	}
	p = keys[r]
	return
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// First strategy guide
func strat1(player1 rune, player2 rune) (p1, p2 int) {

	if Abs(pts(player1)) == pts(player2) {
		p1 = pts(player1) - 3
		p2 = pts(player2) + 3
		return
	}
	if player1 == 'A' {
		if player2 == 'Y' {
			// Win
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'Z' {
			// Lose
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	if player1 == 'B' {
		if player2 == 'Z' {
			// Win
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'X' {
			// Lose
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	if player1 == 'C' {
		if player2 == 'X' {
			// Win
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'Y' {
			// Lose
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	return
}

// Second strategy guide
func strat2(player1 rune, player2 rune) (p1, p2 int) {

	if player2 == 'Y' {
		// Draw
		p1 = pts(player1) - 3
		p2 = Abs(p1)
		return
	}
	if player1 == 'A' {
		// Rock
		if player2 == 'Z' {
			// Win Y
			player2 = 'Y' // Paper
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'X' {
			// Lose Z
			player2 = 'Z' // Scissors
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	if player1 == 'B' {
		// Paper
		if player2 == 'Z' {
			// Win
			player2 = 'Z' // Scissors
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'X' {
			// Lose
			player2 = 'X' // Rock
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	if player1 == 'C' {
		// Scissors
		if player2 == 'Z' {
			// Win
			player2 = 'X' // Rock
			p1 = pts(player1)
			p2 = pts(player2) + 6
		} else if player2 == 'X' {
			// Lose
			player2 = 'Y' // Paper
			p1 = pts(player1) - 6
			p2 = pts(player2)
		}
	}
	return
}

func main() {

	const filename string = "input.txt"
	g1 := NewGame(filename)
	g1.Play(strat1)
	g2 := NewGame(filename)
	g2.Play(strat2)
}
