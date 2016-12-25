package sim

import (
	"errors"
	"fmt"
	"sort"
)

// Sorted event chain.
type EventChain struct {
	chain []*Transaction
	name  string
}

// New returns a new sorted event chain by specified name.
// Slice of transactions has length 0 and capacity 20.
func NewChain(name string) *EventChain {
	return &EventChain{make([]*Transaction, 0, 20), name}
}

// Insert adds new transaction in sorted chain.
func (ch *EventChain) Insert(tr *Transaction) error {
	if ch.Len() == 0 {
		ch.chain = append(ch.chain, tr)
	} else {
		result := make([]*Transaction, len(ch.chain)+1)
		position := sort.Search(ch.Len(), func(i int) bool { return GetTime(*ch.chain[i]) >= GetTime(*tr) })

		result = append(ch.chain[:position], append([]*Transaction{tr}, ch.chain[position:]...)...)
		ch.chain = result
	}
	if sort.IsSorted(ch) {
		return nil
	} else {
		return errors.New("chain is not sorted")
	}
}

// Len returns length of chain.
func (ch EventChain) Len() int {
	return len(ch.chain)
}

// Less returns result of comparison two elements of chain.
func (ch EventChain) Less(i, j int) bool {
	return GetTime(*ch.chain[i]) < GetTime(*ch.chain[j])
}

// Swap swaps two elements of chain.
func (ch EventChain) Swap(i, j int) {
	ch.chain[i], ch.chain[j] = ch.chain[j], ch.chain[i]
}

// String returns information about chain.
func (ch EventChain) String() string {
	string := fmt.Sprintf("CHAIN \"%s\", LENGTH: %d, SORTED: %t] \n", ch.name, ch.Len(), sort.IsSorted(ch))
	for i := 0; i < ch.Len(); i++ {
		string += fmt.Sprintf("\t%s\n", ch.chain[i])
	}
	return string
}

// GetHead returns slice of transaction with least value of timer.
func (ch *EventChain) GetHead() ([]*Transaction, error) {
	length := len(ch.chain)
	if length == 0 {
		return nil, errors.New("no transaction in chain")
	}
	tailPosition := 1
	earliestTime := GetTime(*ch.chain[0])

	for i := tailPosition; i < length; i++ {
		if GetTime(*ch.chain[i]) == earliestTime {
			tailPosition++
		} else {
			break
		}
	}
	head := ch.chain[:tailPosition]
	ch.chain = ch.chain[tailPosition:]
	if len(head) < 1 {
		return nil, errors.New("no transaction in chain")
	}
	return head, nil
}
