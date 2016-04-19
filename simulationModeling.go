package main

import (
	"fmt"
	"transaction"
)

func main() {
	var myFirstTr transaction.Transaction
	myFirstTr.Test(10)
	fmt.Println(myFirstTr)
}
