package sim

import (
	"chain"
	"errors"
	"fmt"
	"math/rand"
	"transaction"
)

const (
	NAvailable = iota
	NUsed
	Used
)

type Pair struct {
	Left, Right float64
}

type Sim struct {
	points      int
	pointState  []int
	idCounter   int
	simTime     float64
	fec         *chain.EventChain
	waitingList []*transaction.Transaction
	finish      bool
}

func New(Points int) *Sim {
	return &Sim{Points, make([]int, Points), 0, 0.0, chain.New("FEC"), make([]*transaction.Transaction, 0, 10), true}
}

func (s *Sim) Init() {
	for i, _ := range s.pointState {
		s.pointState[i] = NUsed
	}
	s.finish = false
	fmt.Println(">> Simulator initialization")
}

// GENERATE block.
func (s *Sim) Generate(NextTime float64, TargetPoint int) error {
	s.idCounter++
	return s.fec.Insert(transaction.New(s.idCounter, s.simTime+NextTime, TargetPoint))
}

// ADVANCE block.
func (s *Sim) Advance(Tr *transaction.Transaction, NextTime float64, NextPoint int) {
	Tr.CorrectTime(NextTime, NextPoint)
}

// GATE and TEST block (check point's state).
func (s *Sim) Test(List []int) (bool, error) {
	fmt.Println("GEBUG PRINT FOR TEST: ", List)
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

// SEIZE block.
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

// RELEASE block.
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

// TERMINATE block.
func (s *Sim) Terminate() {
	s.finish = true
}

// Set next point for transaction, release and seize point.
func (s *Sim) UsePoint(Tr *transaction.Transaction, NextTime float64, NextPoint int) error {
	fmt.Println("GEBUG PRINT FOR USE: ", Tr)
	points := transaction.GetPoints(*Tr)
	if err := s.ReleasePoint(points.Current); err != nil {
		return err
	}
	if err := s.SeizePoint(points.Next); err != nil {
		return err
	}
	Tr.CorrectTime(NextTime, NextPoint)
	if err := s.fec.Insert(Tr); err != nil {
		return err
	}
	//
	fmt.Println(s.pointState)
	//
	return nil
}

// Get current events chain.
func (s *Sim) Extraction() ([]*transaction.Transaction, error) {
	if cec, err := s.fec.GetHead(); err != nil {
		return nil, err
	} else {
		return cec, nil
	}
}

// Check end of simulation.
func (s *Sim) IsFinish() bool {
	return s.finish
}

// Print simulation info.
func (s *Sim) String() string {
	return fmt.Sprintf(">>> Simulation time: %f, total transaction: %d, trancsaction in FEC: %d", s.simTime, s.idCounter, s.fec.Len())
}

// Get uniformly distributed random number.
func Uniform(R *rand.Rand, Limits Pair) (float64, error) {
	if Limits.Left > Limits.Right {
		return 0.0, errors.New(fmt.Sprintf("incorrect limits in Uniform: (%f, %f)", Limits.Left, Limits.Right))
	} else {
		return Limits.Left + R.Float64()*(Limits.Right-Limits.Left), nil
	}
}
