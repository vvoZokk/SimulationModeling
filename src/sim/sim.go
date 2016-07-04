// Package sim implements simulator and general operations for simulation modeling.
package sim

import (
	"errors"
	"fmt"
	"math/rand"
	"sim/chain"
	"sim/statistic"
	"sim/transaction"
)

// List of supported states of points.
const (
	NAvailable = iota
	NUsed
	Used
)

// One pair of float64 numbers for generating random value.
type Pair struct {
	Left, Right float64
}

// Simulator.
type Sim struct {
	points         int
	pointState     []int
	idCounter      int
	simTime        float64
	fec            *chain.EventChain
	pointStatistic []statistic.Unit
	waitingList    []*transaction.Transaction
	finish         bool
}

// New returns new simulator by specified number of points.
// Simulator has a future event chain, a slice of statistic unit for each point and a waitlist of transactions.
// Waitlist is slice with length 0 and capacity 10.
func New(Points int) *Sim {
	return &Sim{Points,
		make([]int, Points),
		0,
		0.0,
		chain.New("FEC"),
		make([]statistic.Unit, Points),
		make([]*transaction.Transaction, 0, 10),
		true}
}

// Init makes initiation of simulator.
func (s *Sim) Init() {
	for i, _ := range s.pointState {
		s.pointState[i] = NUsed
	}
	s.finish = false
	fmt.Println("> Simulation initialization")
}

// Generate creates new transaction in simulator by target waypoint.
func (s *Sim) Generate(NextTime float64, TargetPoint int) error {
	s.idCounter++
	return s.fec.Insert(transaction.New(s.idCounter, s.simTime+NextTime, TargetPoint))
}

// Advance moves transaction to next waypoint by specified time.
func (s *Sim) Advance(Tr *transaction.Transaction, NextTime float64, NextPoint int) {
	Tr.CorrectTime(NextTime, NextPoint)
}

// Test returns result of check of state of point.
func (s *Sim) Test(List []int) (bool, error) {
	for _, point := range List {
		if !(point < s.points) {
			return false, errors.New("incorrect point's id in Sim.Test")
		}
		if s.pointState[point] != NUsed {
			return false, nil
		}
	}
	return true, nil
}

// SeizePoint sets "Used" state of point.
func (s *Sim) SeizePoint(Point int) error {
	if Point < s.points {
		if s.pointState[Point] != NAvailable {
			s.pointState[Point] = Used
			return nil
		} else {
			return errors.New("point not available in Sim.Seize")
		}
	} else {
		return errors.New("incorrect point's id in Sim.Seize")
	}
}

// SeizePoint sets "NUsed" state of point.
func (s *Sim) ReleasePoint(Point int) error {
	if Point < s.points {
		if s.pointState[Point] != NAvailable {
			s.pointState[Point] = NUsed
			return nil
		} else {
			return errors.New("point not available in Sim.Seize")
		}
	} else {
		return errors.New("incorrect point's id in Sim.Seize")
	}
}

// Terminate completes simulation.
func (s *Sim) Terminate() {
	s.finish = true
}

// AddToWaitlist adds transaction to waitlist.
func (s *Sim) AddToWaitlist(Tr *transaction.Transaction) int {
	s.waitingList = append(s.waitingList, Tr)
	return len(s.waitingList)
}

// RemoveFromWaitlist removes transaction from waitlist.
func (s *Sim) RemoveFromWaitlist(Tr *transaction.Transaction) int {
	number, check := 0, false
	for i := 0; i < len(s.waitingList); i++ {
		if transaction.GetId(*s.waitingList[i]) == transaction.GetId(*Tr) {
			number = i
			check = true
			break
		}
	}
	if check {
		s.waitingList = append(s.waitingList[:number], s.waitingList[number+1:]...)
	}
	return len(s.waitingList)
}

// UsePoint releases current, seizes next waypoint and sets next waypoint for transaction,
func (s *Sim) UsePoint(Tr *transaction.Transaction, NextTime float64, NextPoint int) error {
	//fmt.Println("GEBUG PRINT IN USE: ", Tr)
	points := transaction.GetPoints(*Tr)
	if err := s.ReleasePoint(points.Current); err != nil {
		return err
	}
	if err := s.SeizePoint(points.Next); err != nil {
		return err
	}
	Tr.CorrectTime(NextTime, NextPoint)
	//fmt.Println("GEBUG PRINT IN USE BEFORE CORRECTION: ", Tr)
	if err := s.fec.Insert(Tr); err != nil {
		return err
	}
	return nil
}

// GetSimTime returns current value of sumulation timer.
func (s *Sim) GetSimTime() float64 {
	return s.simTime
}

// CorrectTime changes simulation timer by specified value.
func (s *Sim) CorrectTime(NewTime float64) error {
	if s.simTime = NewTime; s.simTime == 0 {
		return errors.New("simulation not started")
	}
	return nil
}

// Extraction returns current events chain.
func (s *Sim) Extraction() ([]*transaction.Transaction, error) {
	if cec, err := s.fec.GetHead(); err != nil {
		return nil, err
	} else {
		s.simTime = transaction.GetTime(*cec[0])
		return cec, nil
	}
}

// GetWaitlist returns waitlist.
func (s *Sim) GetWaitlist() []*transaction.Transaction {
	return s.waitingList
}

// AddStatistic adds new statistic value for point.
func (s *Sim) AddStatistic(Point int, Value float64) error {
	if Point < s.points {
		s.pointStatistic[Point].AddValue(Value)
		return nil
	} else {
		return errors.New("incorrect point's id in Sim.Seize")
	}
}

// GetStatistic returns mean and summary values of statistic for point.
func (s *Sim) GetStatistic(Point int) (float64, float64, error) {
	if Point < s.points {
		return s.pointStatistic[Point].GetMean(), s.pointStatistic[Point].GetSum(), nil
	} else {
		return 0.0, 0.0, errors.New("incorrect point's id in Sim.Seize")
	}
}

// IsFinish returns result of check of ending.
func (s *Sim) IsFinish() bool {
	if s.finish {
		fmt.Println("> Simulation end")
	}
	return s.finish
}

// String returns information about current state of simulation.
func (s *Sim) String() string {
	if len(s.waitingList) != 0 {
		return fmt.Sprintf("> Simulation time: %.1f, total transaction: %d, in FEC: %d, in waitlist: %d", s.simTime, s.idCounter, s.fec.Len(), len(s.waitingList))
	} else {
		return fmt.Sprintf("> Simulation time: %.1f, total transaction: %d, in FEC: %d", s.simTime, s.idCounter, s.fec.Len())
	}
}

// DebugString returns detailed information about simulation.
func (s *Sim) DebugString() string {
	log := "\nSIM DEBUG\n"
	log += fmt.Sprintf("SIM TIME: %.1f, TRANSACTION: TOTAL %d, IN FEC %d, IN WAITLIST %d\n",
		s.simTime,
		s.idCounter,
		s.fec.Len(),
		len(s.waitingList))
	log += fmt.Sprintln(s.fec)
	log += fmt.Sprintf("WAITLIST, LENGTH: %d\n", len(s.waitingList))
	for _, tr := range s.waitingList {
		log += fmt.Sprintf("\t%s\n", tr)
	}
	return log
}

// Uniform returns uniformly distributed random number between specified limits.
func Uniform(R *rand.Rand, Limits Pair) (float64, error) {
	if Limits.Left > Limits.Right {
		return 0.0, errors.New(fmt.Sprintf("incorrect limits in Uniform: (%f, %f)", Limits.Left, Limits.Right))
	} else {
		return Limits.Left + R.Float64()*(Limits.Right-Limits.Left), nil
	}
}
