// Miguel Nobre Castro
// https://adventofcode.com/2022/day/4

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

// Reader class.
type Reader struct {
	num       int
	contained int
	overlaps  int
	pair      chan string
}

// Reader constructor.
func NewReader(filename string) (reader *Reader) {

	pair := make(chan string, 1)
	buffer := ""
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(pair)
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
				defer close(pair)
				break
			}
			pair <- buffer
			buffer = ""
		}
		f.Close()
	}()

	reader = &Reader{
		num:       0,
		contained: 0,
		overlaps:  0,
		pair:      pair,
	}
	return
}

// Finds fully overlapped (Part1) and total overlapped (Part2) assignment pairs.
func (reader *Reader) FindOverlaps() {

	for pair := range reader.pair {
		split := strings.Split(pair, ",")
		elf1 := strings.Split(split[0], "-")
		elf2 := strings.Split(split[1], "-")
		elf1_a, _ := strconv.Atoi(elf1[0])
		elf1_b, _ := strconv.Atoi(elf1[1])
		elf2_a, _ := strconv.Atoi(elf2[0])
		elf2_b, _ := strconv.Atoi(elf2[1])

		if (elf1_a <= elf2_a && elf2_a <= elf1_b) && (elf1_a <= elf2_b && elf2_b <= elf1_b) {
			reader.contained += 1
			reader.overlaps += 1
		} else if (elf2_a <= elf1_a && elf1_a <= elf2_b) && (elf2_a <= elf1_b && elf1_b <= elf2_b) {
			reader.contained += 1
			reader.overlaps += 1
		} else if (elf1_a <= elf2_a && elf2_a <= elf1_b) || (elf1_a <= elf2_b && elf2_b <= elf1_b) || (elf2_a <= elf1_a && elf1_a <= elf2_b) || (elf2_a <= elf1_b && elf1_b <= elf2_b) {
			reader.overlaps += 1
		}
		reader.num += 1
	}
}

func main() {

	const filename string = "input.txt"
	r := NewReader(filename)
	r.FindOverlaps()
	fmt.Printf("There are %d fully contained assignment pairs.\n", r.contained)
	fmt.Printf("There are %d overlapped assignment pairs.\n", r.overlaps)
}
