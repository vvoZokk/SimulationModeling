package main

import (
	//"chain"
	"fmt"
	"math/rand"
	"os"
	"sim"
	"time"
	//"transaction"
)

const (
	Points = 9

	PointAw = iota
	PointA
	PointBw
	PointB
	PointCw
	PointC
	PointAC
	PointBC
	Terminate
)

func main() {

	Rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	CLSim := sim.New(Points)
	//Timings := map[string]sim.Paar{}

	if Time, err := sim.Uniform(Rand, sim.Paar{35, 55}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		CLSim.Generate(Time, PointAw)
	}
	CLSim.Generate(43, PointBw)

	/*
		myFirstTr := transaction.New(1, 43, 3)
		myFirstChain := chain.New("FEC")
		myRand := c

		myFirstTr.CorrectTime(10, 4)
		myFirstChain.Insert(myFirstTr)

		if _, err := myFirstChain.Insert(transaction.New(3, 130, 3)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if Time, err := sim.Uniform(myRand, 35, 55); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			myFirstChain.Insert(mySim.Generate(Time, 6))
		}
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
	*/
}
