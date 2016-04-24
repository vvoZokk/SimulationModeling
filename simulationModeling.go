package main

import (
	//"chain"
	"fmt"
	"math/rand"
	"os"
	"sim"
	"time"
	"transaction"
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

type Checks struct {
	cur, next int
	check     bool
}

func GenUniform(S *sim.Sim, R *rand.Rand, Limits sim.Pair, PointList []int) {
	for _, point := range PointList {
		if time, err := sim.Uniform(R, Limits); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			S.Generate(time, point)
		}
	}
}

func TimerCorrectionPhase(S *sim.Sim, CheckTable map[transaction.Points][]int) {
	cec, err := S.Extraction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, tr := range cec {
		check, err := S.Test(CheckTable[transaction.GetPoints(*tr)])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//
		fmt.Println(check)
	}
}

func main() {

	// Init section

	timings := map[string]sim.Pair{
		"Station": sim.Pair{35, 55},
		"AC":      sim.Pair{12, 18},
		"BC":      sim.Pair{17, 23},
		"Timer":   sim.Pair{1440, 1440}}

	checks := map[transaction.Points][]int{
		transaction.Points{PointAw, PointA}:  []int{PointAC, PointCw},
		transaction.Points{PointBw, PointB}:  []int{PointBC, PointCw},
		transaction.Points{PointAC, PointC}:  []int{PointBC},
		transaction.Points{PointCw, PointBC}: []int{PointBC},
		transaction.Points{PointBC, PointC}:  []int{PointAC},
		transaction.Points{PointCw, PointAC}: []int{PointAC},
	}

	transfers := map[Checks]string{
		{PointAw, PointA, false}:      "wait",
		{PointAw, PointA, true}:       "tc0_AC + use_A + gen_A", // >A****C****B
		{PointA, PointAC, true /*_*/}: "tc0_C + use_AC",         // A->***C****B
		{PointAC, PointC, true}:       "tcUnif_BC + use_C",      // A***->C****B
		{PointAC, PointC, false}:      "tcUnif_Cw + use_C",      // A***->W****B
		{PointC, PointCw, true /*_*/}: "tc0_BC + use_Cw",        // A****>W<****B
		{PointCw, PointBC, false}:     "wait",                   // A****>W<****B
		{PointCw, PointBC, true}:      "tcUnif_B + use_BC",      // A****W->***B
		{PointC, PointBC, true /*_*/}: "tcUnif_B + use_BC",      // A****C->***B
		{PointBC, PointB, true /*_*/}: "tc0_ + use_B",           // A****C***->B

		{PointBw, PointB, false}:      "wait",
		{PointBw, PointB, true}:       "tc0_BC + use_B + gen_B", // A****C****B<
		{PointB, PointBC, true /*_*/}: "tc0_C + use_BC",         // A****C***<-B
		{PointBC, PointC, true}:       "tcUnif_BC + use_C",      // A****C<-***B
		{PointBC, PointC, false}:      "tcUnif_Cw + use_C",      // A****W<-***B
		{PointC, PointCw, true /*_*/}: "tc0_AC + use_Cw",        // A****>W<****B
		{PointCw, PointAC, false}:     "wait",                   // A****>W<****B
		{PointCw, PointAC, true}:      "tcUnif_A + use_AC",      // A***<-W****B
		{PointC, PointAC, true /*_*/}: "tcUnif_A + use_AC",      // A***<-C****B
		{PointAC, PointA, true /*_*/}: "tc0_ + use_A",           // A<-***C****B
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	CLSim := sim.New(Points)

	// Begin simulation

	GenUniform(CLSim, rand, timings["Timer"], []int{Terminate})
	GenUniform(CLSim, rand, timings["Station"], []int{PointAw, PointBw})
	fmt.Println(CLSim)

	TimerCorrectionPhase(CLSim, checks)

}
