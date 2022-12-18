// Miguel Nobre Castro
// https://adventofcode.com/2022/day/8

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// TreeGrid struct
type TreeGrid struct {
	grid    [][]int
	num     int
	visible int
	score   int
}

// TreeGrid constructor.
func NewTreeGrid(filename string) (g *TreeGrid) {

	grid := make([][]int, 0)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for true {
		row := make([]int, 0) // Empty row
		s, err := r.ReadString('\n')
		if !errors.Is(err, io.EOF) {
			s = s[:len(s)-1]
			if len(s) > 0 {
				for _, digit := range s {
					row = append(row, int(digit-'0'))
				}
				grid = append(grid, row)
			}
		} else {
			break
		}
	}
	f.Close()

	g = &TreeGrid{
		grid:    grid,
		num:     len(grid) * len(grid[0]),
		visible: -1,
		score:   -1,
	}
	return
}

// Inspects the TreeGrid for visible trees and best scenic score.
func (g *TreeGrid) Inspect() {

	visible := g.num // Total visible trees
	score := 0       // Best scenic score
	for i := 1; i < len(g.grid)-1; i++ {
		for j := 1; j < len(g.grid[0])-1; j++ {
			// Look up-down-left-right
			up := g.grid[i][j] > g.grid[i-1][j]
			down := g.grid[i][j] > g.grid[i+1][j]
			left := g.grid[i][j] > g.grid[i][j-1]
			right := g.grid[i][j] > g.grid[i][j+1]
			// Current senic score
			cur_score := 1
			if up && i > 1 {
				// Look all up
				k := i - 2
				for k > -1 {
					if g.grid[i][j] <= g.grid[k][j] {
						up = false
						cur_score *= i - k
						break
					}
					k -= 1
				}
				if up {
					cur_score *= i - k - 1
				}
			}
			if down && i < len(g.grid)-2 {
				// Look all down
				k := i + 2
				for k < len(g.grid) {
					if g.grid[i][j] <= g.grid[k][j] {
						down = false
						cur_score *= k - i
						break
					}
					k += 1
				}
				if down {
					cur_score *= k - i - 1
				}
			}
			if left && j > 1 {
				// Look all left
				k := j - 2
				for k > -1 {
					if g.grid[i][j] <= g.grid[i][k] {
						left = false
						cur_score *= j - k
						break
					}
					k -= 1
				}
				if left {
					cur_score *= j - k - 1
				}
			}
			if right && j < len(g.grid[0])-2 {
				// Look all right
				k := j + 2
				for k < len(g.grid[0]) {
					if g.grid[i][j] <= g.grid[i][k] {
						right = false
						cur_score *= k - j
						break
					}
					k += 1
				}
				if right {
					cur_score *= k - j - 1
				}
			}
			if !up && !down && !left && !right {
				visible -= 1
			}
			if cur_score > score {
				score = cur_score
			}
		}
	}
	// Update struct vars
	g.visible = visible
	g.score = score
	return
}

func main() {
	const filename string = "input.txt"
	g := NewTreeGrid(filename)
	g.Inspect()
	fmt.Printf("Visible: %d, Senic Score: %d\n", g.visible, g.score)
}
