// Miguel Nobre Castro
// https://adventofcode.com/2022/day/10

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Signal struct.
type Signal struct {
	cycle    int
	X        int
	strength int
}

// Device struct.
type Device struct {
	cycle int
	X     int
	list  []Signal
	cmds  chan string
	crt   *CRT
}

// Device constructor.
func NewDevice(filename string, rows int, cols int) (dev *Device) {

	cmds := make(chan string, 0)
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(cmds)
		}

		r := bufio.NewReader(f)
		for true {
			s, err := r.ReadString('\n')
			if !errors.Is(err, io.EOF) {
				s = s[:len(s)-1]
				if len(s) > 0 {
					// Adds a cmd to the generator
					cmds <- s
				} else {
					continue
				}
			} else {
				defer close(cmds)
				break
			}
		}
		f.Close()
	}()

	dev = &Device{
		cycle: 0,
		X:     1,
		cmds:  cmds,
		crt:   NewCRT(rows, cols),
	}
	return
}

// [PRIVATE] Checks the strength of the signal at the current cycle.
func (dev *Device) _CheckStrength() {

	sig := Signal{
		cycle:    dev.cycle,
		X:        dev.X,
		strength: dev.cycle * dev.X,
	}
	dev.list = append(dev.list, sig)
	return
}

// Executes the input list of instructions.
// The strength of the signal will be computed for cycle 'cycle1st' and for each 'interval' cycles.
func (dev *Device) Execute(cycle1st int, interval int) {

	for instruction := range dev.cmds {
		if instruction[0:4] == "addx" {
			V, _ := strconv.Atoi(instruction[5:])
			dev.cycle++ // 1st cycle
			dev.crt.DrawPx(dev.X)
			if dev.cycle == cycle1st || (dev.cycle-cycle1st)%interval == 0 {
				dev._CheckStrength()
			}
			dev.cycle++ // 2nd cycle
			dev.crt.DrawPx(dev.X)
			if dev.cycle == cycle1st || (dev.cycle-cycle1st)%interval == 0 {
				dev._CheckStrength()
			}
			dev.X += V // increment register
		} else if instruction[0:4] == "noop" {
			dev.cycle++
			dev.crt.DrawPx(dev.X)
			if dev.cycle == cycle1st || (dev.cycle-cycle1st)%interval == 0 {
				dev._CheckStrength()
			}
		}
	}
	return
}

// Returns the sum of signal strengths.
func (dev *Device) GetStrengths() (val int) {

	val = 0
	for _, sig := range dev.list {
		val += sig.strength
	}
	return
}

// Flushes the CRT screen.
func (dev *Device) FlushScreen() {

	for _, row := range dev.crt.pixels {
		for _, px := range row {
			fmt.Printf("%c", px)
		}
		fmt.Printf("\n")
	}
}

// CRT screen struct.
type CRT struct {
	pixels [][]rune
	rows   int
	cols   int
	x      int
	y      int
}

// CRT screen constructor.
func NewCRT(rows int, cols int) (crt *CRT) {

	pxs := make([][]rune, rows)
	for i := 0; i < rows; i++ {
		pxs[i] = make([]rune, cols)
	}
	crt = &CRT{
		pixels: pxs,
		rows:   rows,
		cols:   cols,
		x:      0,
		y:      0,
	}
	return
}

// Draws a pixel in the CRT screen given the current register 'X'.
func (crt *CRT) DrawPx(X int) {

	if X-1 <= crt.x && crt.x <= X+1 {
		crt.pixels[crt.y][crt.x] = '#' // lit
	} else {
		crt.pixels[crt.y][crt.x] = '.' // dark
	}
	crt.x++
	if crt.x == crt.cols {
		crt.y++
		crt.x = 0
	}
}

func main() {

	const filename string = "input.txt"
	dev := NewDevice(filename, 6, 40)
	dev.Execute(20, 40)
	fmt.Printf("The sum of the signal strengths is %d.\n", dev.GetStrengths())
	dev.FlushScreen()
}
