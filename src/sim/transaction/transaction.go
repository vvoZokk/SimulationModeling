// Package transaction implements single transaction for simulation modeling.
package transaction

import (
	"fmt"
)

// Indexes of current and next waypoints of a transaction.
type Points struct {
	Current, Next int
}

// Single transacctioin.
type Transaction struct {
	id, currentPoint, nextPoint int
	time, lifetime              float64
}

// New returns new transaction by id, initial value of timer and index of next waypoint.
func New(Id int, Time float64, NextPoint int) *Transaction {
	return &Transaction{Id, 0, NextPoint, Time, 0}
}

// CorrectTimer sets new value of time, new points for transaction and makes time shift.
func (tr *Transaction) CorrectTime(NewTime float64, NewNextPoint int) {
	tr.lifetime += NewTime
	tr.time += NewTime
	tr.currentPoint, tr.nextPoint = tr.nextPoint, NewNextPoint
}

// Wait sets new value of time and makes time shift without change points.
func (tr *Transaction) Wait(WaitingTime float64) {
	tr.lifetime += WaitingTime
	tr.time += WaitingTime
}

// String returns information about transaction.
func (tr Transaction) String() string {
	return fmt.Sprintf("TRANSACTION [%d, %f, %d, %d], LIFETIME: %f", tr.id, tr.time, tr.currentPoint, tr.nextPoint, tr.lifetime)
}

// GetParam returns all parameters of transaction.
func GetParam(Tr Transaction) (int, float64, int, int, float64) {
	return Tr.id, Tr.time, Tr.currentPoint, Tr.nextPoint, Tr.lifetime
}

// GetId returns transaction's ID.
func GetId(Tr Transaction) int {
	return Tr.id
}

// GetTime returns value of transaction's timer.
func GetTime(Tr Transaction) float64 {
	return Tr.time
}

// GetPoints returns current and next transaction's waypoints.
func GetPoints(Tr Transaction) Points {
	return Points{Tr.currentPoint, Tr.nextPoint}
}
