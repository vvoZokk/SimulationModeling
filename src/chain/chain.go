package chain

import (
	"errors"
	"fmt"
	"sort"
	"transaction"
)

type EventChain struct {
	chain  []*transaction.Transaction
	sorted bool
	name   string
}

func New(Name string) *EventChain {
	return &EventChain{make([]*transaction.Transaction, 0, 10), true, Name}
}

// Insert new transaction in sorted chain.
func (ch *EventChain) Insert(tr *transaction.Transaction) (*EventChain, error) {
	if ch.Len() == 0 {
		ch.chain = append(ch.chain, tr)
	} else {

		result := make([]*transaction.Transaction, len(ch.chain)+1)
		position := sort.Search(ch.Len(), func(i int) bool { return transaction.GetTime(*ch.chain[i]) >= transaction.GetTime(*tr) })

		result = append(ch.chain[:position], append([]*transaction.Transaction{tr}, ch.chain[position:]...)...)
		ch.chain = result
	}
	if sort.IsSorted(ch) {
		return ch, nil
	} else {
		return nil, errors.New("chain is not sorted")
	}
}

// Get chain length for sort.Interface.
func (ch EventChain) Len() int {
	return len(ch.chain)
}

// Less function for sort.Interface.
func (ch EventChain) Less(i, j int) bool {
	return transaction.GetTime(*ch.chain[i]) < transaction.GetTime(*ch.chain[j])
}

// Swap function for sort.Interface.
func (ch EventChain) Swap(i, j int) {
	ch.chain[i], ch.chain[j] = ch.chain[j], ch.chain[i]
}

// Print debug chain info.
func (ch EventChain) String() string {
	string := fmt.Sprintf("CHAIN \"%s\", LENGTH: %d, SORTED: %t] \n", ch.name, ch.Len(), sort.IsSorted(ch))
	for i := 0; i < ch.Len(); i++ {
		string += fmt.Sprintf("\t%s\n", ch.chain[i])
	}
	return string
}

// Get slice of transaction with earliest time parameter.
func (ch *EventChain) GetHead() ([]*transaction.Transaction, error) {
	length := len(ch.chain)
	if length == 0 {
		return nil, errors.New("no transaction in chain")
	}
	tailPosition := 1
	earliestTime := transaction.GetTime(*ch.chain[0])

	for i := 1; i < length; i++ {
		if transaction.GetTime(*ch.chain[i]) == earliestTime {
			tailPosition++
		} else {
			break
		}
	}
	head := ch.chain[:tailPosition]
	ch.chain = ch.chain[tailPosition:]
	return head, nil
}
