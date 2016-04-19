package chain

import (
	//"errors"
	"transaction"
)

type EventChain []*transaction.Transaction

func (slice EventChain) Len() int {
	return len(slice)
}

func (slice EventChain) Less(i, j int) bool {
	return transaction.GetTime(*slice[i]) < transaction.GetTime(*slice[j])
}

func (slice EventChain) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
