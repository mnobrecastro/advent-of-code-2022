// Miguel Nobre Castro
// https://adventofcode.com/2022/day/9

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Point struct.
type Point struct {
	val  [2]int
	next *Point
}

// List struct (to track Rope's tail visited Points).
type List struct {
	num   int
	idx   int
	first *Point
}

// List constructor.
func NewList(idx int) (l *List) {

	l = &List{
		num:   0,
		idx:   idx,
		first: nil,
	}
	return
}

// Inserts visited 'point':=(x,y) in sorting order, iff it has not been visited before.
func (l *List) Insert(point [2]int) {

	if l.first == nil {
		p := &Point{
			val:  point,
			next: nil,
		}
		l.first = p
		l.num += 1
	} else {
		ptr := l.first
		for ptr != nil {
			if point[1] < ptr.val[1] {
				// Prepend
				p := &Point{
					val:  point,
					next: ptr,
				}
				l.first = p
				l.num += 1
				break
			} else if point[1] > ptr.val[1] {
				if ptr.next != nil {
					if point[1] < ptr.next.val[1] {
						p := &Point{
							val:  point,
							next: ptr.next,
						}
						ptr.next = p
						l.num += 1
						break
					} else {
						ptr = ptr.next
					}
				} else {
					// Append
					p := &Point{
						val:  point,
						next: nil,
					}
					ptr.next = p
					l.num += 1
					break
				}
			} else {
				// If equal, the point is not inserted
				break
			}
		}
	}
	return
}

// Rope struct.
type Rope struct {
	pos_knots [][2]int
	moves     chan string
	data      []*List
}

// Rope constructor.
// Takes a minimum of two knots (head & tail).
func NewRope(knots int, filename string) (rope *Rope) {

	if knots < 2 {
		rope = nil
		panic(errors.New("ERROR: The rope is too short! Please use 2 or more knots."))
	}

	moves := make(chan string, 0)
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(moves)
		}

		r := bufio.NewReader(f)
		for true {
			s, err := r.ReadString('\n')
			if !errors.Is(err, io.EOF) {
				s = s[:len(s)-1]
				if len(s) > 0 {
					// Adds a move to the generator
					moves <- s
				} else {
					continue
				}
			} else {
				defer close(moves)
				break
			}
		}
		f.Close()
	}()

	// Knots positions
	pos_knots := make([][2]int, 0) // slice of [2]int
	for i := 0; i < knots; i++ {
		pos_knots = append(pos_knots, [2]int{0, 0})
	}

	// Tail's tracking list of visited positions
	list := NewList(0)
	list.Insert([2]int{0, 0})
	data := make([]*List, 0) // slice of *List
	data = append(data, list)

	rope = &Rope{
		pos_knots: pos_knots,
		moves:     moves,
		data:      data,
	}
	return
}

// Absolute value of int.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// [Private] Checks if a 'pos' can be inserted in the visited ones.
func (rope *Rope) _CheckInsert(pos [2]int) {

	for _, list := range rope.data {
		if list.idx == pos[0] {
			list.Insert(pos)
			return
		}
	}
	l := NewList(pos[0])
	l.Insert(pos)
	rope.data = append(rope.data, l)
	return
}

// Moves the head knot of the rope according to the provided moves.
func (rope *Rope) MoveHead() {

	for move := range rope.moves {
		split := strings.Split(move, " ")
		direction := split[0]
		steps, _ := strconv.Atoi(split[1])
		for i := 0; i < steps; i++ {
			// Move 'head' once
			switch direction {
			case "U":
				rope.pos_knots[0][0] += 1
			case "D":
				rope.pos_knots[0][0] -= 1
			case "L":
				rope.pos_knots[0][1] -= 1
			case "R":
				rope.pos_knots[0][1] += 1
			default:
				continue
			}
			// Check and move the remaining knots
			for j := 1; j < len(rope.pos_knots); j++ {
				diff := [2]int{rope.pos_knots[j-1][0] - rope.pos_knots[j][0], rope.pos_knots[j-1][1] - rope.pos_knots[j][1]}
				if Abs(diff[0])+Abs(diff[1]) == 2 && (diff[0] == 0 || diff[1] == 0) {
					if diff[0] > 0 { // U
						rope.pos_knots[j][0] += 1
					}
					if diff[0] < 0 { // D
						rope.pos_knots[j][0] -= 1
					}
					if diff[1] < 0 { // L
						rope.pos_knots[j][1] -= 1
					}
					if diff[1] > 0 { // R
						rope.pos_knots[j][1] += 1
					}
				} else if Abs(diff[0])+Abs(diff[1]) > 2 {
					if diff[0] > 0 && diff[1] > 0 { // Q1
						rope.pos_knots[j][0] += 1
						rope.pos_knots[j][1] += 1
					}
					if diff[0] < 0 && diff[1] > 0 { // Q2
						rope.pos_knots[j][0] -= 1
						rope.pos_knots[j][1] += 1
					}
					if diff[0] < 0 && diff[1] < 0 { // Q3
						rope.pos_knots[j][0] -= 1
						rope.pos_knots[j][1] -= 1
					}
					if diff[0] > 0 && diff[1] < 0 { // Q4
						rope.pos_knots[j][0] += 1
						rope.pos_knots[j][1] -= 1
					}
				}
			}
			rope._CheckInsert(rope.pos_knots[len(rope.pos_knots)-1])
		}
	}
	return
}

// Counts the number positions visited, at least once, by the tail of the Rope.
// (a boolean can be passed as argument to print those nodes)
func (rope *Rope) CountTailPos(print bool) (k int) {

	k = 0
	for _, list := range rope.data {
		k += list.num
		ptr := list.first
		for ptr != nil {
			if print {
				fmt.Printf("(%d,%d) ", ptr.val[0], ptr.val[1])
			}
			ptr = ptr.next
		}
	}
	if print {
		fmt.Printf("\n")
	}
	return
}

func main() {

	const filename string = "input.txt"
	rope := NewRope(2, filename)
	rope.MoveHead()
	fmt.Printf("The tail visited %d positions at least once.\n", rope.CountTailPos(false))
	rope10 := NewRope(10, filename)
	rope10.MoveHead()
	fmt.Printf("The tail visited %d positions at least once.\n", rope10.CountTailPos(false))
}
