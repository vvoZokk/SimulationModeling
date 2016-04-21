package sim

import (
	"chain"
	"errors"
	"fmt"
	"math/rand"
	"transaction"
)

const (
	NUsed = iota
	Used
	NAvailable
)

type Paar struct {
	Left, Right float64
}

type Sim struct {
	points     int
	pointState []int
	idCounter  int
	simTime    float64
	fec        *chain.EventChain
}

func New(Points int) *Sim {
	fmt.Println(">> Simulator initialization")
	return &Sim{Points, make([]int, Points), 0, 0.0, chain.New("FEC")}
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
func (s *Sim) Test(Point int) (bool, error) {
	if !(Point < s.points) {
		return false, errors.New("incorrect point's id in Sim.Test")
	}
	if s.pointState[Point] == NUsed {
		return true, nil
	} else {
		return false, nil
	}
}

// Get uniformly distributed random number.
func Uniform(R *rand.Rand, Paar Paar) (float64, error) {
	if Paar.Left > Paar.Right {
		return 0.0, errors.New(fmt.Sprintf("incorrect limits in Uniform: (%f, %f)", Paar.Left, Paar.Right))
	} else {
		return Paar.Left + R.Float64()*(Paar.Right-Paar.Left), nil
	}
}
