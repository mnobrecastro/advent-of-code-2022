// Miguel Nobre Castro
// https://adventofcode.com/2022/day/6

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Data struct.
type Data struct {
	val  rune
	next *Data
}

// Data struct constructor.
func NewData(r rune) (data *Data) {

	data = &Data{
		val:  r,
		next: nil,
	}
	return
}

// Buffer (FIFO) struct.
type Buffer struct {
	num    int
	length int
	MAXLEN int
	head   *Data
	tail   *Data
	signal chan rune
}

// Buffer struct constructor.
func NewBuffer(MAXLEN int, filename string) (buff *Buffer) {
	signal := make(chan rune, 1)
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(signal)
		}

		r := bufio.NewReader(f)
		for true {
			char, _, err := r.ReadRune()
			if !errors.Is(err, io.EOF) && char != '\n' {
				signal <- char
			} else {
				defer close(signal)
				break
			}
		}
		f.Close()
	}()

	buff = &Buffer{
		num:    0,
		length: 0,
		MAXLEN: MAXLEN,
		head:   nil,
		tail:   nil,
		signal: signal,
	}
	return
}

// Inserts a Data element into the Buffer.
func (buff *Buffer) Enqueue(data *Data) {

	if buff.tail == nil && buff.length == 0 {
		buff.head = data
	} else {
		if buff.length == buff.MAXLEN {
			// No overflow error message required
			buff.Dequeue()
		}
		buff.tail.next = data
	}
	buff.tail = data
	buff.length += 1
	buff.num += 1
	return
}

// Removes the head Data element from the Buffer.
func (buff *Buffer) Dequeue() error {

	if buff.head == nil && buff.length == 0 {
		return errors.New("Underflow: the queue is empty.")
	}
	temp := buff.head
	buff.head = buff.head.next
	temp.next = nil // Deleted by the garbage collector
	buff.length -= 1
	return nil
}

// Reads the signal to detect a start-of-packet marker.
func (buff *Buffer) Read() (num int) {

	num = 0
	for char := range buff.signal {
		buff.Enqueue(NewData(char))
		if buff.IsMarker() {
			num = buff.num
			break
		}
	}
	return
}

// Detects a start-of-packet marker.
func (buff *Buffer) IsMarker() (b bool) {

	b = true
	if buff.length < buff.MAXLEN {
		b = false
	} else {
		ptr1 := buff.head
		for ptr1 != buff.tail && b {
			ptr2 := ptr1.next
			for ptr2 != nil && b {
				if ptr1.val == ptr2.val {
					b = false
				}
				ptr2 = ptr2.next
			}
			ptr1 = ptr1.next
		}
	}
	return
}

func main() {

	const filename string = "input.txt"
	signal := NewBuffer(4, filename)
	val := signal.Read()
	if val > 0 {
		fmt.Printf("First marker detected after character %d.\n", val)
	} else {
		fmt.Print("No marker was detected in the signal.\n")
	}
	sign14 := NewBuffer(14, filename)
	val = sign14.Read()
	if val > 0 {
		fmt.Printf("First marker detected after character %d.\n", val)
	} else {
		fmt.Print("No marker was detected in the signal.\n")
	}
}
