package main

import (
	"chain"
	"fmt"
	"os"
	"transaction"
)

func main() {
	myFirstTr := transaction.New(1, 43, 3)
	myFirstChain := chain.New("FEC")

	myFirstTr.CorrectTime(10, 4)
	myFirstChain.Insert(myFirstTr)

	if _, err := myFirstChain.Insert(transaction.New(3, 130, 3)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	myFirstChain.Insert(transaction.New(4, 53, 3))
	fmt.Println(myFirstChain)

	if head, err := myFirstChain.GetHead(); err != nil {
		fmt.Println(err)
		os.Exit(1)

	} else {
		fmt.Println("Transsactions in head:")
		for _, tr := range head {
			fmt.Println(tr)
		}
	}
	fmt.Printf("FEC current length %d", myFirstChain.Len())
}
