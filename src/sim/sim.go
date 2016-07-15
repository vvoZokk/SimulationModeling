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
func New(points int) *Sim {
	return &Sim{points,
		make([]int, points),
		0,
		0.0,
		chain.New("FEC"),
		make([]statistic.Unit, points),
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
func (s *Sim) Generate(nextTime float64, targetPoint int) error {
	s.idCounter++
	return s.fec.Insert(transaction.New(s.idCounter, s.simTime+nextTime, targetPoint))
}

// Advance moves transaction to next waypoint by specified time.
func (s *Sim) Advance(tr *transaction.Transaction, nextTime float64, nextPoint int) {
	tr.CorrectTime(nextTime, nextPoint)
}

// Test returns result of check of state of point.
func (s *Sim) Test(listOfPoint []int) (bool, error) {
	for _, point := range listOfPoint {
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
func (s *Sim) SeizePoint(p int) error {
	if p < s.points {
		if s.pointState[p] != NAvailable {
			s.pointState[p] = Used
			return nil
		} else {
			return errors.New("point not available in Sim.SeizePoint")
		}
	} else {
		return errors.New("incorrect point's id in Sim.SeizePoint")
	}
}

// SeizePoint sets "NUsed" state of point.
func (s *Sim) ReleasePoint(p int) error {
	if p < s.points {
		if s.pointState[p] != NAvailable {
			s.pointState[p] = NUsed
			return nil
		} else {
			return errors.New("point not available in Sim.ReleasePoint")
		}
	} else {
		return errors.New("incorrect point's id in Sim.ReleasePoint")
	}
}

// Terminate completes simulation.
func (s *Sim) Terminate() {
	s.finish = true
}

// AddToWaitlist adds transaction to waitlist.
func (s *Sim) AddToWaitlist(tr *transaction.Transaction) int {
	s.waitingList = append(s.waitingList, tr)
	return len(s.waitingList)
}

// RemoveFromWaitlist removes transaction from waitlist.
func (s *Sim) RemoveFromWaitlist(tr *transaction.Transaction) int {
	number, check := 0, false
	for i := 0; i < len(s.waitingList); i++ {
		if transaction.GetId(*s.waitingList[i]) == transaction.GetId(*tr) {
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
func (s *Sim) UsePoint(tr *transaction.Transaction, nextTime float64, nextPoint int) error {
	//fmt.Println("GEBUG PRINT IN USE: ", Tr)
	points := transaction.GetPoints(*tr)
	if err := s.ReleasePoint(points.Current); err != nil {
		return err
	}
	if err := s.SeizePoint(points.Next); err != nil {
		return err
	}
	tr.CorrectTime(nextTime, nextPoint)
	//fmt.Println("GEBUG PRINT IN USE BEFORE CORRECTION: ", Tr)
	if err := s.fec.Insert(tr); err != nil {
		return err
	}
	return nil
}

// GetSimTime returns current value of sumulation timer.
func (s *Sim) GetSimTime() float64 {
	return s.simTime
}

// CorrectTime changes simulation timer by specified value.
func (s *Sim) CorrectTime(newTime float64) error {
	if s.simTime = newTime; s.simTime == 0 {
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
func (s *Sim) AddStatistic(point int, value float64) error {
	if point < s.points {
		s.pointStatistic[point].AddValue(value)
		return nil
	} else {
		return errors.New("incorrect point's id in Sim.AddStatistic")
	}
}

// GetStatistic returns mean and summary values of statistic for point.
func (s *Sim) GetStatistic(point int) (float64, float64, error) {
	if point < s.points {
		return s.pointStatistic[point].Mean(), s.pointStatistic[point].Sum(), nil
	} else {
		return 0.0, 0.0, errors.New("incorrect point's id in Sim.GetStatistic")
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
func Uniform(r *rand.Rand, limits Pair) (float64, error) {
	if limits.Left > limits.Right {
		return 0.0, errors.New(fmt.Sprintf("incorrect limits in Uniform: (%f, %f)", limits.Left, limits.Right))
	} else {
		return limits.Left + r.Float64()*(limits.Right-limits.Left), nil
	}
}
