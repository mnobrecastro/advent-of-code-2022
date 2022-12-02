// Miguel Nobre Castro
// https://adventofcode.com/2022/day/1

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Elf struct.
type Elf struct {
	idx   int
	items []int
	n     int
	sum   int
	prev  *Elf
	next  *Elf
}

// Double linked list of Elves.
type List struct {
	num  int
	head *Elf
	tail *Elf
}

// Appends an Elf to the list given its 'idx' and 'items'.
func (l *List) append(idx int, items []int) {

	// If an Elf doesn't carry anything, he is lazy, thus discarded
	if len(items) == 0 {
		return
	}

	// Sum of cals carried by the Elf
	sum := 0
	for _, item := range items {
		sum += item
	}

	// Create an Elf
	Elf := &Elf{
		idx:   idx,
		items: items,
		n:     len(items),
		sum:   sum,
	}

	// Push the Elf down the list
	if l.head == nil && l.tail == nil {
		l.head = Elf
		l.tail = Elf
		l.num += 1
	} else {
		Elf.prev = l.tail
		l.tail.next = Elf
		l.tail = Elf
		l.num += 1
	}
	return
}

// Sorts the Elves by their decreasing 'sum' of cals.
func (l *List) sort() {

	quicksort(l.head, l.tail)
	return
}

// Checks whether the Elves behave and keep their positions.
func check(first *Elf, last *Elf) (behave bool) {

	behave = false
	ptr := first
	for ptr != nil && !behave {
		if ptr == last {
			behave = true
		}
		ptr = ptr.next
	}
	return
}

// Quicksorting them all in descending order.
func quicksort(first *Elf, last *Elf) {

	if !check(first, last) {
		return
	}
	pivot := partition(first, last)
	quicksort(first, pivot.prev)
	quicksort(pivot.next, last)
	return
}

// Partitioning the Elves by decreasing 'sum' of cals.
func partition(first *Elf, last *Elf) (i *Elf) {

	i = nil
	j := first
	pivot := last
	for j != pivot {
		if j.sum >= pivot.sum {
			// Increment i
			if i == nil {
				i = first
			} else {
				i = i.next
			}
			// Swap i with j
			if i != j {
				i.idx, j.idx = j.idx, i.idx
				i.items, j.items = j.items, i.items
				i.n, j.n = j.n, i.n
				i.sum, j.sum = j.sum, i.sum
			}
		}
		j = j.next
	}
	// Increment i once
	if i == nil {
		i = first
	} else {
		i = i.next
	}
	// Swap i with pivot
	if i != pivot {
		i.idx, pivot.idx = pivot.idx, i.idx
		i.items, pivot.items = pivot.items, i.items
		i.n, pivot.n = pivot.n, i.n
		i.sum, pivot.sum = pivot.sum, i.sum
	}
	return
}

// Prints the list of Elves.
func (l *List) print() {
	ptr := l.head
	if ptr == nil {
		fmt.Printf("Empty list of Elves.\n")
		return
	}
	for ptr != nil {
		fmt.Printf("Elf no.%d carries %d items (%d cals).\n", ptr.idx, ptr.n, ptr.sum)
		ptr = ptr.next
	}
	return
}

/* Reads an input file .txt consisting of a list of items.
 *
 * Input file example:
 * "
 * 1000
 * 2000
 * 3000
 *
 * 4000
 *
 * 5000
 * 6000
 *
 * 7000
 * 8000
 * 9000
 *
 * 10000
 * "
 */
func read_input(filename string) (elves List) {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	// Declare an 'idx' and an empty slice of int ('items')
	idx := 0
	items := make([]int, 0)

	r := bufio.NewReader(f)
	for true {
		s, err := r.ReadString('\n')
		if !errors.Is(err, io.EOF) {
			s = s[:len(s)-1]
			if len(s) > 0 {
				// Read a 'val' and append to 'items'
				val, _ := strconv.Atoi(s)
				items = append(items, val)
			} else if len(items) > 0 {
				// Add the elf and its items to the list of 'elves'
				elves.append(idx, items)
				items = make([]int, 0)
				idx += 1
			}
		} else {
			break
		}
	}
	f.Close()

	return
}

func main() {

	const filename string = "input.txt"
	elves := read_input(filename)
	elves.print()
	fmt.Println("********************************")

	elves.sort()
	elves.print()
	fmt.Println("********************************")

	total := 0
	k := 0
	ptr := elves.head
	for k < 3 {
		total += ptr.sum
		ptr = ptr.next
		k += 1
	}
	fmt.Printf("The top three Elves carry a total of %d Calories.\n", total)
	return
}
