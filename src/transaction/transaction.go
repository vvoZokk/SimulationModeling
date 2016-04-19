package transaction

import (
	"fmt"
)

type Transaction struct {
	id, currentPoint, nextPoint int
	time, lifetime              float64
}

// Create new transaction with id, time and next point.
func New(Id int, Time float64, NextPoint int) *Transaction {
	return &Transaction{Id, 0, NextPoint, Time, 0}
}

// Correct time, shift and set points for transaction.
func (tr *Transaction) CorrectTime(NewTime float64, NewNextPoint int) {
	tr.lifetime += tr.time
	tr.time += NewTime
	tr.currentPoint, tr.nextPoint = tr.nextPoint, NewNextPoint
}

// Base print transaction info.
func (tr Transaction) String() string {
	return fmt.Sprintf("TRANSACTION [%d, %f, %d, %d], LIFETIME: %f", tr.id, tr.time, tr.currentPoint, tr.nextPoint, tr.lifetime)
}

// Get all transaction's paramerts.
func GetParam(tr Transaction) (int, float64, int, int, float64) {
	return tr.id, tr.time, tr.currentPoint, tr.nextPoint, tr.lifetime
}

// Get transaction's time.
func GetTime(tr Transaction) float64 {
	return tr.time
}
