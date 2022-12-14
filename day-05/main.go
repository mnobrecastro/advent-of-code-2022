// Miguel Nobre Castro
// https://adventofcode.com/2022/day/5

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

// Crate class.
type Crate struct {
	val  rune
	prev *Crate
}

// Crate class constructor.
func NewCrate(val rune) (crate *Crate) {
	crate = &Crate{
		val:  val,
		prev: nil,
	}
	return
}

// Stack class.
type Stack struct {
	num    int
	bottom *Crate // Useful for 'Prepend'
	top    *Crate
}

// Stack class constructor.
func NewStack() (stack *Stack) {

	stack = &Stack{
		num:    0,
		bottom: nil,
		top:    nil,
	}
	return
}

// Pushes a Crate to the Stack
func (stack *Stack) Push(crate *Crate) {

	crate.prev = stack.top
	stack.top = crate
	stack.num += 1
	return
}

// Pops a Crate from the Stack and returns it.
func (stack *Stack) Pop() (crate *Crate) {

	if stack.top == nil {
		crate = nil
		return
	} else {
		crate = stack.top
		stack.top = stack.top.prev
		crate.prev = nil
		stack.num -= 1
	}
	return
}

// Prepends a Crate to the Stack (top-down approach)
func (stack *Stack) Prepend(crate *Crate) {

	if stack.bottom == nil && stack.top == nil {
		stack.bottom = crate
		stack.top = crate
	} else {
		stack.bottom.prev = crate
		stack.bottom = crate
	}
	return
}

// Cargo class.
type Cargo struct {
	num_stacks   int
	num_crates   int
	stacks       []*Stack
	instructions chan string
}

// Cargo class constructor.
func NewCargo(filename string) (cargo *Cargo) {

	num_stacks := 0
	num_crates := 0
	var stacks []*Stack

	// Read initial Cargo config
	firstline := true
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for true {
		s, err := r.ReadString('\n')
		if !errors.Is(err, io.EOF) {
			s = s[:len(s)-1]
			if len(s) > 0 {
				matched, _ := regexp.MatchString(`[A-Z]`, s)
				if !matched {
					// Reached the instructions
					break
				}
				// Allocate the 'stacks'
				if firstline {
					i := 0
					for i < (len(s)+1)/4 {
						fmt.Printf("Added a new stack %d\n", i)
						stacks = append(stacks, NewStack())
						num_stacks += 1
						i += 1
					}
					firstline = false
				}
				// Prepend the 'crates' to each 'stack'
				i := 0
				for i < (len(s)+1)/4 {
					//re := regexp.MustCompile(`\[[A-Z]\]*`)
					//sr := re.FindAllString(s[i*4:i*4+2], -1)[0]
					r := rune(s[i*4+1])
					if r != ' ' {
						stacks[i].Prepend(NewCrate(r))
						fmt.Printf("Prepended crate %c to stack %d\n", r, i)
					}
					num_crates += 1
					i += 1
				}
			}
		} else {
			break
		}
	}
	f.Close()

	// Generate channel of string for the crane instructions
	instructions := make(chan string, 1)
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(instructions)
		}

		r := bufio.NewReader(f)
		for true {
			s, err := r.ReadString('\n')
			if !errors.Is(err, io.EOF) {
				s = s[:len(s)-1]
				if len(s) > 0 && s[0] == 'm' {
					// Add each instruction to the generator
					instructions <- s
				} else {
					continue
				}
			} else {
				defer close(instructions)
				break
			}
		}
		f.Close()
	}()

	cargo = &Cargo{
		num_stacks:   num_stacks,
		num_crates:   num_crates,
		stacks:       stacks,
		instructions: instructions,
	}
	return
}

// Prints the Cargo object.
func (cargo *Cargo) Print() {

	fmt.Printf("Cargo has %d crates arranged in %d stacks > ", cargo.num_crates, cargo.num_stacks)
	i := 0
	for i < len(cargo.stacks) {
		fmt.Printf("%c", cargo.stacks[i].top.val)
		i += 1
	}
	fmt.Printf("\n")
	return
}

// Moves the Crates from one Stack to another based on a set of instructions.
func (cargo *Cargo) MoveCrates() {

	for instruction := range cargo.instructions {
		re := regexp.MustCompile(`[0-9]+`)
		move := re.FindAllString(instruction, -1)
		num, _ := strconv.Atoi(move[0]) // Number of crates to move
		sc, _ := strconv.Atoi(move[1])
		sc -= 1 // Source crate
		tc, _ := strconv.Atoi(move[2])
		tc -= 1 // Target crate
		i := 0
		for i < num {
			cargo.stacks[tc].Push(cargo.stacks[sc].Pop())
			fmt.Printf("Moving [%c] from %d to %d\n", cargo.stacks[tc].top.val, sc+1, tc+1)
			i += 1
		}
	}
	return
}

// Moves groups of Crates from one Stack to another (CraneMover9001).
func (cargo *Cargo) Move9001() {

	for instruction := range cargo.instructions {
		re := regexp.MustCompile(`[0-9]+`)
		move := re.FindAllString(instruction, -1)
		num, _ := strconv.Atoi(move[0]) // Number of crates to move
		sc, _ := strconv.Atoi(move[1])
		sc -= 1 // Source crate
		tc, _ := strconv.Atoi(move[2])
		tc -= 1 // Target crate
		i := 0
		var crates []*Crate
		for i < num {
			crates = append(crates, cargo.stacks[sc].Pop())
			i += 1
		}
		i = num - 1
		for i >= 0 {
			cargo.stacks[tc].Push(crates[i])
			i -= 1
		}
	}
	return
}

func main() {

	const filename string = "input.txt"
	cargo := NewCargo(filename)
	cargo.Print()
	cargo.MoveCrates()
	cargo.Print()
	cargo9k1 := NewCargo(filename)
	cargo9k1.Move9001()
	cargo9k1.Print()
}
