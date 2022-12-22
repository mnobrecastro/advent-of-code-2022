// Miguel Nobre Castro
// https://adventofcode.com/2022/day/11

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Item struct.
type Item struct {
	wlevel big.Int
	next   *Item
}

// Queue struct.
type Queue struct {
	size int
	head *Item
	tail *Item
}

// Queue (FIFO) constructor.
func NewQueue() *Queue {

	q := &Queue{
		size: 0,
		head: nil,
		tail: nil,
	}
	return q
}

// Inserts an element at the end of the Queue.
func (q *Queue) Enqueue(wlevel big.Int) {

	item := &Item{
		wlevel: wlevel,
		next:   nil,
	}
	if q.head == nil && q.tail == nil && q.size == 0 {
		q.head = item
		q.tail = item
	} else {
		q.tail.next = item
		q.tail = item
	}
	q.size++
	return
}

// Removes the first element in the Queue.
func (q *Queue) Dequeue() {

	if !(q.head == nil && q.size == 0) {
		// Queue underflow is ignored
		if q.head == q.tail && q.size == 1 {
			q.tail = nil
		}
		q.head = q.head.next
		q.size--
	}
	return
}

// Monkey struct.
type Monkey struct {
	idx       int
	num       int
	q         *Queue
	operation rune
	O         *big.Int // operation const
	T         *big.Int // test const
	throwT    int
	throwF    int
}

// Monkey constructor.
func NewMonkey(idx int, op rune, O int, T int, throwT int, throwF int) *Monkey {

	monkey := &Monkey{
		idx:       idx,
		num:       0,
		q:         NewQueue(),
		operation: op,
		O:         big.NewInt(int64(O)),
		T:         big.NewInt(int64(T)),
		throwT:    throwT,
		throwF:    throwF,
	}
	return monkey
}

// Returns the number of items inspected by the Monkey.
func (monkey *Monkey) GetNumInspected() int {
	return monkey.num
}

// Returns the number of items that Monkey has.
func (monkey *Monkey) GetNumItems() int {
	return monkey.q.size
}

// The Monkey receives an item with worry level 'wlevel'.
func (monkey *Monkey) CatchItem(wlevel big.Int) {

	monkey.q.Enqueue(wlevel)
	return
}

// The Monkey inspects the first item in is queue given my worry factor 'wfactor'.
func (monkey *Monkey) InspectItem(wfactor *big.Int) (big.Int, int) {

	wlevel := &monkey.q.head.wlevel
	monkey.q.Dequeue()
	// "Please be careful..."
	switch monkey.operation {
	case '*':
		if monkey.O.Int64() > 0 {
			wlevel.Mul(wlevel, monkey.O)
		} else {
			wlevel.Mul(wlevel, wlevel)
		}
	case '+':
		if monkey.O.Int64() > 0 {
			wlevel.Add(wlevel, monkey.O)
		} else {
			wlevel.Add(wlevel, wlevel)
		}
	}
	// "Not that stressed..."
	wlevel.Div(wlevel, wfactor)
	monkey.num++
	// "Are you testing the item or my patience?!"
	target := 0
	var modulus big.Int
	/*var mod_6 big.Int
	mod_6.Mod(wlevel, big.NewInt(int64(6)))
	if mod_6.Int64() == -1 || mod_6.Int64() == 1 {
		target = monkey.throwF
	} else {*/
	modulus.Mod(wlevel, monkey.T)
	if modulus.Int64() == 0 {
		target = monkey.throwT
	} else {
		target = monkey.throwF
	}
	//}
	/*var mod_2, mod_3, mod_5, modulus big.Int
	mod_2.Mod(wlevel, big.NewInt(int64(2)))
	mod_3.Mod(wlevel, big.NewInt(int64(3)))
	mod_5.Mod(wlevel, big.NewInt(int64(5)))
	if mod_2.Int64()+mod_3.Int64()+mod_5.Int64() == 0 {
		target = monkey.throwF
	} else {
		modulus.Mod(wlevel, monkey.T)
		if modulus.Int64() == 0 {
			target = monkey.throwT
		} else {
			target = monkey.throwF
		}
	}*/
	return *wlevel, target
}

// Prints an instance of Monkey.
func (monkey *Monkey) Print() {

	fmt.Printf("Monkey %d:\n", monkey.idx)
	fmt.Printf("  Starting items: ")
	ptr := monkey.q.head
	for ptr != nil {
		fmt.Printf("%d ", ptr.wlevel.Int64())
		ptr = ptr.next
	}
	fmt.Printf("\n")
	fmt.Printf("  Operation: new = old %c ", monkey.operation)
	if monkey.O.Int64() > 0 {
		fmt.Printf("%d\n", monkey.O.Int64())
	} else {
		fmt.Printf("old\n")
	}
	fmt.Printf("  Test: divisible by %d\n", monkey.T)
	fmt.Printf("    If true: throw to monkey %d\n", monkey.throwT)
	fmt.Printf("    If false: throw to monkey %d\n", monkey.throwF)
	return
}

// Troop of monkeys struct.
type Troop struct {
	size    int
	monkeys []*Monkey
	leaders [2]*Monkey
	rounds  int
}

// Troop of monkeys constructor.
func NewTroop() *Troop {

	troop := &Troop{
		size:    0,
		monkeys: make([]*Monkey, 0),
		leaders: [2]*Monkey{nil, nil},
		rounds:  0,
	}
	return troop
}

// Adds a Monkey 'm' to the Troop.
func (troop *Troop) AddMonkey(m *Monkey) {

	troop.monkeys = append(troop.monkeys, m)
	troop.size++
	return
}

// Adds multiple Monkeys to the Troop from a 'filename'.
func (troop *Troop) FromFile(filename string, bPrint bool) {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var s_idx string
	var s_items []string
	var s_operation string
	var s_O string
	var s_T string
	var s_throwT string
	var s_throwF string

	r := bufio.NewReader(f)
	for true {
		s, err := r.ReadString('\n')
		if !errors.Is(err, io.EOF) {
			s = s[:len(s)-1]
			if len(s) > 0 {
				split := strings.Split(s, ":")
				if regexp.MustCompile(`Monkey`).FindAllString(split[0], -1) != nil {
					s_idx = regexp.MustCompile(`[0-9]+`).FindAllString(split[0], -1)[0]
				} else if regexp.MustCompile(`Starting items`).FindAllString(split[0], -1) != nil {
					for _, val := range regexp.MustCompile(`[0-9]+`).FindAllString(split[1], -1) {
						s_items = append(s_items, val)
					}
				} else if regexp.MustCompile(`Operation`).FindAllString(split[0], -1) != nil {
					s_operation = regexp.MustCompile(`\+|\*`).FindAllString(split[1], -1)[0]
					if len(regexp.MustCompile(`old`).FindAllString(split[1], -1)) < 2 {
						s_O = regexp.MustCompile(`[0-9]+`).FindAllString(split[1], -1)[0]
					} else {
						s_O = "-1"
					}
				} else if regexp.MustCompile(`Test`).FindAllString(split[0], -1) != nil {
					s_T = regexp.MustCompile(`[0-9]+`).FindAllString(split[1], -1)[0]
				} else if regexp.MustCompile(`If true`).FindAllString(split[0], -1) != nil {
					s_throwT = regexp.MustCompile(`[0-9]+`).FindAllString(split[1], -1)[0]
				} else if regexp.MustCompile(`If false`).FindAllString(split[0], -1) != nil {
					s_throwF = regexp.MustCompile(`[0-9]+`).FindAllString(split[1], -1)[0]

					idx, _ := strconv.Atoi(s_idx)
					op := rune(s_operation[0])
					O, _ := strconv.Atoi(s_O)
					T, _ := strconv.Atoi(s_T)
					throwT, _ := strconv.Atoi(s_throwT)
					throwF, _ := strconv.Atoi(s_throwF)
					monkey := NewMonkey(idx, op, O, T, throwT, throwF)
					for _, item := range s_items {
						wlevel, _ := strconv.Atoi(item)
						monkey.CatchItem(*big.NewInt(int64(wlevel)))
					}
					if bPrint {
						monkey.Print()
					}
					troop.AddMonkey(monkey)
					s_items = make([]string, 0)
				}
			}
		} else {
			break
		}
	}
	f.Close()

	fmt.Printf("Size of troop: %d\n", len(troop.monkeys))
	return
}

// All monkeys in the Troop inspect their items give my worry factor 'wfactor'. A round takes place.
func (troop *Troop) InspectionRound(wfactor *big.Int) {

	for i := 0; i < troop.size; i++ {
		monkey := troop.monkeys[i]
		for j := monkey.GetNumItems(); j > 0; j-- {
			wlevel, target := monkey.InspectItem(wfactor)
			troop.monkeys[target].CatchItem(wlevel)
		}
		// Monkey leaders
		if troop.leaders[0] == nil && troop.leaders[1] == nil {
			troop.leaders[0] = monkey
			troop.leaders[1] = monkey
		} else {
			if monkey != troop.leaders[0] && monkey != troop.leaders[1] {
				if monkey.GetNumInspected() > troop.leaders[0].GetNumInspected() || monkey.GetNumInspected() > troop.leaders[1].GetNumInspected() {
					if troop.leaders[0].GetNumInspected() >= troop.leaders[1].GetNumInspected() {
						troop.leaders[1] = monkey
					} else {
						troop.leaders[0] = monkey
					}
				}
			}
		}
	}
	troop.rounds++
	return
}

// Prints all monkeys in the Troop.
func (troop *Troop) Print() {

	for _, monkey := range troop.monkeys {
		monkey.Print()
	}
	return
}

// Prints the number of items inspected by each monkey in the Troop.
func (troop *Troop) PrintInspected() {

	for _, monkey := range troop.monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", monkey.idx, monkey.GetNumInspected())
	}
	return
}

// Returns the size of the Troop, i.e. the number of monkeys.
func (troop *Troop) GetSize() int {

	return troop.size
}

// Returns the level of "monkey business".
func (troop *Troop) GetBusiness() int {

	return troop.leaders[0].GetNumInspected() * troop.leaders[1].GetNumInspected()
}

func main() {

	const filename = "input.txt"
	troop := NewTroop()
	troop.FromFile(filename, true)
	N_ROUNDS := 20
	for i := 0; i < N_ROUNDS; i++ {
		fmt.Printf("Starting round %d...\n", i+1)
		troop.InspectionRound(big.NewInt(3))
	}
	troop.PrintInspected()
	fmt.Printf("Monkey business: %d.\n\n", troop.GetBusiness())

	troop2 := NewTroop()
	troop2.FromFile(filename, false)
	N_ROUNDS2 := 10000
	for i := 0; i < N_ROUNDS2; i++ {
		fmt.Printf("Starting round %d...\n", i+1)
		troop2.InspectionRound(big.NewInt(1))
		if i == 0 {
			troop2.PrintInspected()
		}
		if i == 19 {
			troop2.PrintInspected()
		}
		if (i+1)%1000 == 0 {
			troop2.PrintInspected()
		}
	}
	fmt.Printf("Monkey business: %d.\n", troop2.GetBusiness())
}
