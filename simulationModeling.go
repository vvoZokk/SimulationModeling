package main

import (
	"fmt"
	"math/rand"
	"os"
	"sim"
	"time"
	"transaction"
)

const Points = 8
const ( // List of points
	Point0  = iota
	PointA  // 1
	PointB  // 2
	PointCm // 3, main
	PointCr // 4, reserve
	PointAC // 5
	PointBC // 6
	ClockPoint
)
const ( // List of actions
	Generate = iota
	Wait
	Use
	Terminate
)
const ( // List of limits
	Station = iota
	AC
	BC
	Timer
)

type Checks struct {
	cur, next int
	check     bool
}

type Action struct {
	Type      int
	Arguments []int
}

func GenerateUniform(S *sim.Sim, R *rand.Rand, Limits sim.Pair, PointList []int) {
	for _, point := range PointList {
		if time, err := sim.Uniform(R, Limits); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			S.Generate(S.GetSimTime()+time, point)
		}
	}
}

func UseBlock(S *sim.Sim, Tr *transaction.Transaction, Time float64, NextPoint int) {
	//if NextPoint != Point0
	if err := S.UsePoint(Tr, Time, NextPoint); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Phase(S *sim.Sim, R *rand.Rand, TimeTable map[int]sim.Pair, CheckTable map[transaction.Points][]int, RoadMap map[Checks][]Action) {
	cec, err := S.Extraction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, tr := range cec {
		points := transaction.GetPoints(*tr)
		check, err := S.Test(CheckTable[points])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		actions := RoadMap[Checks{points.Current, points.Next, check}]
		for _, action := range actions {
			if action.Type == Wait {
				// GEBUG PRINT
				fmt.Println("WAIT ACTION", tr)
				S.AddToWaitlist(tr)
			}
			if action.Type == Generate {
				// GEBUG PRINT
				fmt.Println("GENERATE ACTION")
				GenerateUniform(S, R, TimeTable[action.Arguments[0]], []int{action.Arguments[1]})
			}
			if action.Type == Use {
				// GEBUG PRINT
				fmt.Println("USE ACTION")
				switch {
				case action.Arguments[0] == 0:
					UseBlock(S, tr, 0.0, action.Arguments[1])
				default:
					if time, err := sim.Uniform(R, TimeTable[action.Arguments[0]]); err != nil {
						fmt.Println(err)
						os.Exit(1)
					} else {
						UseBlock(S, tr, time, action.Arguments[1])
					}
				}
			}
			if action.Type == Terminate {
				// GEBUG PRINT
				fmt.Println("TERMINATE ACTION")
				S.Terminate()
			}
		}
	}
	for _, tr := range S.GetWaitlist() {
		points := transaction.GetPoints(*tr)
		check, err := S.Test(CheckTable[points])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		actions := RoadMap[Checks{points.Current, points.Next, check}]
		for _, action := range actions {
			waitingTime := S.GetSimTime() - transaction.GetTime(*tr)
			if action.Type == Use {
				// GEBUG PRINT
				fmt.Println("USE ACTION FOR WAITING TRANSACTION")
				switch {
				case action.Arguments[0] == 0:
					UseBlock(S, tr, waitingTime, action.Arguments[1])
				default:
					if time, err := sim.Uniform(R, TimeTable[action.Arguments[0]]); err != nil {
						fmt.Println(err)
						os.Exit(1)
					} else {
						UseBlock(S, tr, waitingTime+time, action.Arguments[1])
					}
				}
				S.RemoveFromWaitlist(tr)
				// GEBUG PRINT
				fmt.Println("WAITING TIME: ", waitingTime)
			}
		}
	}

}

func main() {

	// Init section

	timings := map[int]sim.Pair{
		Station: sim.Pair{35, 55},
		AC:      sim.Pair{12, 18},
		BC:      sim.Pair{17, 23},
		Timer:   sim.Pair{3 * 60, 3 * 60}}

	checks := map[transaction.Points][]int{
		transaction.Points{Point0, PointA}:   []int{PointAC, PointCm, PointCr},
		transaction.Points{Point0, PointB}:   []int{PointBC, PointCm, PointCr},
		transaction.Points{PointA, PointAC}:  []int{PointBC},
		transaction.Points{PointCr, PointBC}: []int{PointBC},
		transaction.Points{PointCm, PointBC}: []int{PointBC},
		transaction.Points{PointB, PointCr}:  []int{PointAC},
		transaction.Points{PointCr, PointAC}: []int{PointAC},
		transaction.Points{PointCm, PointAC}: []int{PointAC},
	}

	transfers := map[Checks][]Action{
		{Point0, PointA, false}:   []Action{Action{Wait, []int{}}},                                                    // >A****C****B
		{Point0, PointA, true}:    []Action{Action{Use, []int{0, PointAC}}, Action{Generate, []int{Station, PointA}}}, // >A****C****B
		{PointA, PointAC, false}:  []Action{Action{Use, []int{AC, PointCr}}},                                          // A>***Cr****B
		{PointA, PointAC, true}:   []Action{Action{Use, []int{AC, PointCm}}},                                          // A>***Cm****B
		{PointAC, PointCm, true}:  []Action{Action{Use, []int{0, PointBC}}},                                           //
		{PointAC, PointCr, true}:  []Action{Action{Use, []int{0, PointBC}}},                                           //
		{PointCm, PointBC, false}: []Action{Action{Wait, []int{}}},                                                    // A***>Cm<***B
		{PointCm, PointBC, true}:  []Action{Action{Use, []int{BC, PointB}}},                                           // A****Cm>***B
		{PointCr, PointBC, false}: []Action{Action{Wait, []int{}}},                                                    // A***>Cr<***B
		{PointCr, PointBC, true}:  []Action{Action{Use, []int{BC, PointB}}},                                           // A****Cr>***B
		{PointBC, PointB, true}:   []Action{Action{Use, []int{0, Point0}}},                                            // A****C***->B

		{Point0, PointB, false}:   []Action{Action{Wait, []int{}}},                                                    // A****C****B<
		{Point0, PointB, true}:    []Action{Action{Use, []int{0, PointBC}}, Action{Generate, []int{Station, PointB}}}, // A****C****B<
		{PointB, PointBC, false}:  []Action{Action{Use, []int{BC, PointCr}}},                                          // A****Cr***<B
		{PointB, PointBC, true}:   []Action{Action{Use, []int{BC, PointCm}}},                                          // A****Cm***<B
		{PointBC, PointCm, true}:  []Action{Action{Use, []int{0, PointAC}}},                                           //
		{PointBC, PointCr, true}:  []Action{Action{Use, []int{0, PointAC}}},                                           //
		{PointCm, PointAC, false}: []Action{Action{Wait, []int{}}},                                                    // A***>Cm<***B
		{PointCm, PointAC, true}:  []Action{Action{Use, []int{AC, PointA}}},                                           // A****Cm>***B
		{PointCr, PointAC, false}: []Action{Action{Wait, []int{}}},                                                    // A***>Cr<***B
		{PointCr, PointAC, true}:  []Action{Action{Use, []int{AC, PointA}}},                                           // A****Cr>***B
		{PointAC, PointA, true}:   []Action{Action{Use, []int{0, Point0}}},                                            // A<-***C****B

		{Point0, ClockPoint, true}: []Action{Action{Terminate, []int{}}}, // Clock
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	CLSim := sim.New(Points)
	CLSim.Init()

	// Begin simulation

	GenerateUniform(CLSim, rand, timings[Timer], []int{ClockPoint})
	GenerateUniform(CLSim, rand, timings[Station], []int{PointA, PointB})
	fmt.Println(CLSim)
	for !CLSim.IsFinish() {
		Phase(CLSim, rand, timings, checks, transfers)
		fmt.Println(CLSim)
	}

}
