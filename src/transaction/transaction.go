package transaction

import (
	"fmt"
)

type Transaction struct {
	Id, CurrentPoint, NextPoint int
	Time, lifetime              float64
}

// Create new transaction with id, time and next point.
func New(Id int, Time float64, NextPoint int) *Transaction {
	return &Transaction{Id, 0, NextPoint, Time, 0}
}

// Correct time, shift and set points for transaction.
func (tr *Transaction) CorrectTime(NewTime float64, NextPoint int) {
	tr.lifetime += tr.Time
	tr.Time += NewTime
	tr.CurrentPoint = tr.NextPoint
	tr.NextPoint = NextPoint
}

// Base print transaction info.
func (tr Transaction) String() string {
	return fmt.Sprintf("TRANSACTION [%d, %f, %d, %d], LIFETIME: %f", tr.Id, tr.Time, tr.CurrentPoint, tr.NextPoint, tr.lifetime)
}
