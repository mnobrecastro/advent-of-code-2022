// Miguel Nobre Castro
// https://adventofcode.com/2022/day/3

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

// SackScanner class.
type SackScanner struct {
	num    int         // Number of scanned sacks
	sum    int         // Sum of item priorities
	badges int         // Sum of badge priorities
	sack   chan string // Sack generator
}

// Constructor of SackScanner from strings using Channels.
func NewSackScanner(filename string) (scan *SackScanner) {

	// A Generator per rucksack
	sack := make(chan string, 1)
	buffer := ""
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(sack)
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
				defer close(sack)
				break
			}
			sack <- buffer
			buffer = ""
		}
		f.Close()
	}()

	scan = &SackScanner{
		num:  0,
		sum:  0,
		sack: sack,
	}
	return
}

// Calculates the priority of a given carried item.
func Priority(r rune) (i int) {

	if int(r) >= int('a') && int(r) <= int('z') {
		i = int(r-'a') + 1
	}
	if int(r) >= int('A') && int(r) <= int('Z') {
		i = int(r-'A') + 27
	}
	return
}

// Inspects both compartments in each rucksack.
func (scan *SackScanner) InspectAll() {

	comp1 := "" // First compartment
	comp2 := "" // Second compartment
	for sack := range scan.sack {
		comp1 += SortString(sack[:len(sack)/2])
		comp2 += SortString(sack[len(sack)/2:])

		size := len(comp1)
		i, j := 0, 0
		r := '0'
		for i != size && j != size {
			if comp1[i] == comp2[j] && r != rune(comp1[i]) {
				r = rune(comp1[i])
				val := Priority(r)
				scan.sum += val
				fmt.Printf("Found item %c with priority %d\n", r, val)
				break
			}

			if int(comp1[i]-'0') <= int(comp2[j]-'0') {
				if i < size {
					i += 1
				} else {
					break
				}
			} else {
				if j < size {
					j += 1
				} else {
					break
				}
			}
		}

		comp1 = ""
		comp2 = ""
		scan.num += 1
	}
}

// Finds the Badge among each three consecutive rucksacks
func (scan *SackScanner) FindBadges() {

	sacks := [3]string{"", "", ""} // Elfs' items
	trio := 0
	for sack := range scan.sack {
		sacks[trio] = SortString(sack)

		if trio == 2 {
			i, j, k := 0, 0, 0
			r := '0'
			for i != len(sacks[0]) && j != len(sacks[1]) && k != len(sacks[2]) {
				if sacks[0][i] == sacks[1][j] && sacks[0][i] == sacks[2][k] && r != rune(sacks[0][i]) {
					r = rune(sacks[0][i])
					val := Priority(r)
					scan.badges += val
					fmt.Printf("Found badge %c with priority %d\n", r, val)
					break
				}

				if int(sacks[0][i]-'0') <= int(sacks[1][j]-'0') && int(sacks[0][i]-'0') <= int(sacks[2][k]-'0') {
					if i < len(sacks[0]) {
						i += 1
					} else {
						break
					}
				} else if int(sacks[1][j]-'0') <= int(sacks[2][k]-'0') {
					if j < len(sacks[1]) {
						j += 1
					} else {
						break
					}
				} else {
					if k < len(sacks[2]) {
						k += 1
					} else {
						break
					}
				}
			}
			sacks = [3]string{"", "", ""}
		}
		if trio == 2 {
			trio = 0
		} else {
			trio += 1
		}
		scan.num += 1
	}
}

// A slice of Rune type to enable sort
type RuneSlice []rune

func (s RuneSlice) Len() (l int) {

	l = len(s)
	return
}

func (s RuneSlice) Less(i int, j int) (b bool) {

	b = s[i] < s[j]
	return
}

func (s RuneSlice) Swap(i int, j int) {

	s[i], s[j] = s[j], s[i]
	return
}

// Sort a string by converting it to a slice of Rune
func SortString(in string) (out string) {

	r := []rune(in)
	sort.Sort(RuneSlice(r))
	out = string(r)
	return
}

func main() {

	const filename string = "input.txt"
	scan1 := NewSackScanner(filename)
	scan1.InspectAll()
	fmt.Printf("The sum of priorities of these item types is %d.\n", scan1.sum)
	scan2 := NewSackScanner(filename)
	scan2.FindBadges()
	fmt.Printf("The sum of priorities of all bages is %d.\n", scan2.badges)
}
