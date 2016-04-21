package sim

import (
	"chain"
	"errors"
	"fmt"
	"math/rand"
	"transaction"
)

type Sim struct {
	points     int
	pointState []int
	idCounter  int
	simTime    float64
	fec        *chain.EventChain
}

func (s *Sim) Generate(NextTime float64, TargetPoint int) *transaction.Transaction {
	tr := transaction.New(s.idCounter, s.simTime+NextTime, TargetPoint)
	s.idCounter++
	return tr
}

func Uniform(Left, Right float64) (float64, error) {
	if Left > Right {
		return 0.0, errors.New(fmt.Sprintf("incorrect limits in Uniform: (%f, %f)", Left, Right))
	} else {
		return Left + rand.Float64()*(Right-Left), nil
	}
}
