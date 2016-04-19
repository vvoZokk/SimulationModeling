package main

import (
	"fmt"
	"transaction"
)

func main() {
	myFirstTr := transaction.New(1, 23.4, 3)
	myFirstTr.CorrectTime(10, 4)
	fmt.Println(myFirstTr)
}
