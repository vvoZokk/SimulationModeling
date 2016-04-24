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
	Points = 8

	//List of points
	Point0 = iota
	PointA
	PointB
	PointCw
	PointC
	PointAC
	PointBC
	ClockPoint

	// List of actions
	Generate = iota
	Wait
	Use
	Terminate
)

type Checks struct {
	cur, next int
	check     bool
}

func GenerateUniform(S *sim.Sim, R *rand.Rand, Limits sim.Pair, PointList []int) {
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
		fmt.Println("in transaction: ", transaction.GetPoints(*tr))
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
		transaction.Points{Point0, PointA}:   []int{PointAC, PointCw},
		transaction.Points{Point0, PointB}:   []int{PointBC, PointCw},
		transaction.Points{PointAC, PointC}:  []int{PointBC},
		transaction.Points{PointCw, PointBC}: []int{PointBC},
		transaction.Points{PointBC, PointC}:  []int{PointAC},
		transaction.Points{PointCw, PointAC}: []int{PointAC},
	}

	transfers := map[Checks]string{
		{Point0, PointA, false}:       "wait",
		{Point0, PointA, true}:        "tc0_AC + use_A + gen_A", // >A****C****B
		{PointA, PointAC, true /*_*/}: "tc0_C + use_AC",         // A->***C****B
		{PointAC, PointC, true}:       "tcUnif_BC + use_C",      // A***->C****B
		{PointAC, PointC, false}:      "tcUnif_Cw + use_C",      // A***->W****B
		{PointC, PointCw, true /*_*/}: "tc0_BC + use_Cw",        // A****>W<****B
		{PointCw, PointBC, false}:     "wait",                   // A****>W<****B
		{PointCw, PointBC, true}:      "tcUnif_B + use_BC",      // A****W->***B
		{PointC, PointBC, true /*_*/}: "tcUnif_B + use_BC",      // A****C->***B
		{PointBC, PointB, true /*_*/}: "tc0_ + use_B",           // A****C***->B

		{Point0, PointB, false}:       "wait",
		{Point0, PointB, true}:        "tc0_BC + use_B + gen_B", // A****C****B<
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

	GenerateUniform(CLSim, rand, timings["Timer"], []int{Terminate})
	GenerateUniform(CLSim, rand, timings["Station"], []int{PointA, PointB})
	fmt.Println(CLSim)

	fmt.Println(checks[transaction.Points{Point0, PointA}])

	fmt.Println(transfers[Checks{PointAC, PointA, true}])

	TimerCorrectionPhase(CLSim, checks)

}
