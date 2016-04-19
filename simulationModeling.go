package main

import (
	"chain"
	"fmt"
	"transaction"
)

func main() {
	myFirstTr := transaction.New(1, 23.4, 3)
	myFirstChain := chain.EventChain{myFirstTr, transaction.New(2, 23, 2)}

	myFirstChain[1].CorrectTime(10, 4)
	fmt.Println(myFirstChain[1])
	fmt.Println("Chain length ", myFirstChain.Len())
}
