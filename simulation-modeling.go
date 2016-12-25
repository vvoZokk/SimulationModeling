package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"simulation-modeling/sim"
	"time"
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
			fmt.Println(err, S.DebugString())
			os.Exit(1)
		} else {
			S.Generate(time, point)
			if point != ClockPoint {
				S.AddStatistic(Point0, time)
			}
		}
	}
}

func UseBlock(S *sim.Sim, Tr *sim.Transaction, Time float64, NextPoint int) {
	if err := S.UsePoint(Tr, Time, NextPoint); err != nil {
		fmt.Println(err, S.DebugString())
		os.Exit(1)
	}
}

func Phases(S *sim.Sim, R *rand.Rand, TimeTable map[int]sim.Pair, CheckTable map[sim.Points][]int, RoadMap map[Checks][]Action) {
	cec, err := S.Extraction()
	if err != nil {
		fmt.Println(err, S.DebugString())
		os.Exit(1)
	}
	for _, tr := range cec {
		points := sim.GetPoints(*tr)
		check, err := S.Test(CheckTable[points])
		if err != nil {
			fmt.Println(err, S.DebugString())
			os.Exit(1)
		}
		actions := RoadMap[Checks{points.Current, points.Next, check}]
		for _, action := range actions {
			if action.Type == Wait {
				// GEBUG PRINT
				//fmt.Println("WAIT ACTION", tr)
				S.AddToWaitlist(tr)
			}
			if action.Type == Generate {
				// GEBUG PRINT
				//fmt.Println("GENERATE ACTION")
				GenerateUniform(S, R, TimeTable[action.Arguments[0]], []int{action.Arguments[1]})
			}
			if action.Type == Use {
				// GEBUG PRINT
				//fmt.Println("USE ACTION")
				switch {
				case action.Arguments[0] == 0:
					UseBlock(S, tr, 0.0, action.Arguments[1])
				default:
					if time, err := sim.Uniform(R, TimeTable[action.Arguments[0]]); err != nil {
						fmt.Println(err, S.DebugString())
						os.Exit(1)
					} else {
						UseBlock(S, tr, time, action.Arguments[1])
						S.AddStatistic(points.Next, time)
					}
				}
			}
			if action.Type == Terminate {
				// GEBUG PRINT
				//fmt.Println("TERMINATE ACTION")
				S.Terminate()
			}
		}
	}
	for i := 0; i < len(S.GetWaitlist()); i++ {
		waitList := S.GetWaitlist()
		points := sim.GetPoints(*waitList[i])
		check, err := S.Test(CheckTable[points])
		if err != nil {
			fmt.Println(err, S.DebugString())
			os.Exit(1)
		}
		actions := RoadMap[Checks{points.Current, points.Next, check}]
		for _, action := range actions {
			waitingTime := S.GetSimTime() - sim.GetTime(*waitList[i])
			if action.Type == Use {
				// GEBUG PRINT
				//fmt.Println("USE ACTION FOR WAITING TRANSACTION")
				switch {
				case action.Arguments[0] == 0:
					UseBlock(S, waitList[i], waitingTime, action.Arguments[1])
				default:
					if time, err := sim.Uniform(R, TimeTable[action.Arguments[0]]); err != nil {
						fmt.Println(err, S.DebugString())
						os.Exit(1)
					} else {
						UseBlock(S, waitList[i], waitingTime+time, action.Arguments[1])
						S.AddStatistic(points.Next, time)
					}
				}
				S.RemoveFromWaitlist(waitList[i])
				if waitingTime != 0 {
					if points.Current == Point0 {
						S.AddStatistic(points.Next, waitingTime)
					} else {
						S.AddStatistic(points.Current, waitingTime)
					}
				}
				// GEBUG PRINT
				//fmt.Println("WAITING TIME: ", waitingTime)
			}
		}
	}
}

func WriteData(Writer *bufio.Writer, Data string) {
	if _, err := Writer.WriteString(Data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetMeanTime(S *sim.Sim, Point int) float64 {
	if meanTime, _, err := S.GetStatistic(Point); err != nil {
		fmt.Println(err, S.DebugString())
		os.Exit(1)
	} else {
		return meanTime
	}
	return 0.0
}

func GetSumTime(S *sim.Sim, Point int) float64 {
	if _, sumTime, err := S.GetStatistic(Point); err != nil {
		fmt.Println(err, S.DebugString())
		os.Exit(1)
	} else {
		return sumTime
	}
	return 0.0
}

func main() {
	duration := 24.0
	outFile := os.Stdout
	defer outFile.Close()

	outputFlag := flag.String("o", "", "write output to file")
	durationFlag := flag.Float64("d", 24, "set simulation duration in hours")
	flag.Parse()
	if *outputFlag != "" {
		if file, err := os.Create(*outputFlag); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			outFile = file
		}
	}
	duration = *durationFlag
	/*
		helpString := fmt.Sprint(fmt.Sprintf("usage: %s [-h] [-o FILE] [-d DURATION]\n\n", filepath.Base(os.Args[0])),
					"Crossing Loop Simulation\n\n",
					"optional arguments:\n",
					"  -h, --help\t show this help message and exit\n",
					"  -o FILE\t write output to FILE\n",
					"  -d DURATION\t set simulation duration in hours (default: 24)")
	*/

	// Init section

	timings := map[int]sim.Pair{
		Station: sim.Pair{35, 55},
		AC:      sim.Pair{12, 18},
		BC:      sim.Pair{17, 23},
		Timer:   sim.Pair{duration * 60, duration * 60}}

	checks := map[sim.Points][]int{
		sim.Points{Point0, PointA}:   []int{PointAC, PointCm, PointCr},
		sim.Points{Point0, PointB}:   []int{PointBC, PointCm, PointCr},
		sim.Points{PointA, PointAC}:  []int{PointBC},
		sim.Points{PointCr, PointBC}: []int{PointBC},
		sim.Points{PointCm, PointBC}: []int{PointBC},
		sim.Points{PointB, PointCr}:  []int{PointAC},
		sim.Points{PointCr, PointAC}: []int{PointAC},
		sim.Points{PointCm, PointAC}: []int{PointAC},
	}

	transfers := map[Checks][]Action{
		{Point0, PointA, false}:   []Action{Action{Wait, []int{}}, Action{Generate, []int{Station, PointA}}},          // >A****C****B
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

		{Point0, PointB, false}:   []Action{Action{Wait, []int{}}, Action{Generate, []int{Station, PointB}}},          // A****C****B<
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
	writer := bufio.NewWriter(outFile)

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	CLSim := sim.New(Points)
	CLSim.Init()

	// Begin simulation

	GenerateUniform(CLSim, rand, timings[Timer], []int{ClockPoint})
	GenerateUniform(CLSim, rand, timings[Station], []int{PointA, PointB})

	for !CLSim.IsFinish() {
		Phases(CLSim, rand, timings, checks, transfers)
		//fmt.Println(CLSim)
	}

	// Get statistic
	WriteData(writer, "Crossing loop simulation statistic\n")
	WriteData(writer, fmt.Sprintf("Duration: %.0f minutes\n", duration*60))
	WriteData(writer, fmt.Sprintf("Mean waiting time on station A: %.2f\n", GetMeanTime(CLSim, PointA)))
	WriteData(writer, fmt.Sprintf("Mean waiting time on station B: %.2f\n", GetMeanTime(CLSim, PointB)))
	waitingTime := GetMeanTime(CLSim, PointCm)
	waitingTime += GetMeanTime(CLSim, PointCr)
	WriteData(writer, fmt.Sprintf("Mean waiting time on crossing: %.2f\n", waitingTime/2))
	WriteData(writer, fmt.Sprintf("Utilization ratio for AC track: %.2f\n", GetSumTime(CLSim, PointAC)/(duration*60)))
	WriteData(writer, fmt.Sprintf("Utilization ratio for BC track: %.2f\n", GetSumTime(CLSim, PointBC)/(duration*60)))

	err := writer.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
