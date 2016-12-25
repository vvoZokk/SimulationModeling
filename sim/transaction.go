package sim

import (
	"fmt"
)

// Indexes of current and next waypoints of a transaction.
type Points struct {
	Current, Next int
}

// Single transaction.
type Transaction struct {
	id, currentPoint, nextPoint int
	time, lifetime              float64
}

// New returns new transaction by id, initial value of timer and index of next waypoint.
func NewTransaction(id int, time float64, nextPoint int) *Transaction {
	return &Transaction{id, 0, nextPoint, time, 0}
}

// CorrectTimer sets new value of time, new points for transaction and makes time shift.
func (tr *Transaction) CorrectTime(newTime float64, newNextPoint int) {
	tr.lifetime += newTime
	tr.time += newTime
	tr.currentPoint, tr.nextPoint = tr.nextPoint, newNextPoint
}

// Wait sets new value of time and makes time shift without change points.
func (tr *Transaction) Wait(waitingTime float64) {
	tr.lifetime += waitingTime
	tr.time += waitingTime
}

// String returns information about transaction.
func (tr Transaction) String() string {
	return fmt.Sprintf("TRANSACTION [%d, %f, %d, %d], LIFETIME: %f", tr.id, tr.time, tr.currentPoint, tr.nextPoint, tr.lifetime)
}

// GetParam returns all parameters of transaction.
func GetParam(tr Transaction) (int, float64, int, int, float64) {
	return tr.id, tr.time, tr.currentPoint, tr.nextPoint, tr.lifetime
}

// GetId returns transaction's ID.
func GetId(tr Transaction) int {
	return tr.id
}

// GetTime returns value of transaction's timer.
func GetTime(tr Transaction) float64 {
	return tr.time
}

// GetPoints returns current and next transaction's waypoints.
func GetPoints(tr Transaction) Points {
	return Points{tr.currentPoint, tr.nextPoint}
}
