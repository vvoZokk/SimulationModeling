package transaction

import "testing"

type testPair struct {
	transaction *Transaction
	timer       float64
	point       int
}

func TestCorrectTime(t *testing.T) {
	tests := []testPair{
		{New(0, 2.0, 1), 3.0, 2},
		{New(0, 3.0, 3), 1.5, 2},
		{New(1, 0.0, 0), 0.0, 0},
	}

	for _, test := range tests {
		point := test.transaction.nextPoint
		time := test.transaction.time

		test.transaction.CorrectTime(test.timer, test.point)

		if newTime := test.transaction.time; newTime != time+test.timer {
			t.Errorf("Expected transaction's time %.2f, got %.2f", time+test.timer, newTime)
		}
		if lifetime := test.transaction.lifetime; lifetime != test.timer {
			t.Errorf("Expected transaction's lifetime %.2f, got %.2f", test.timer, lifetime)
		}
		if currentPoint := test.transaction.currentPoint; currentPoint != point {
			t.Errorf("Expected %d index of current point in transaction, got %d", point, currentPoint)
		}
		if nextPoint := test.transaction.nextPoint; nextPoint != test.point {
			t.Errorf("Expected %d index of next point in transaction, got %d", test.point, nextPoint)
		}
	}
}

func TestWait(t *testing.T) {
	tests := []testPair{
		{New(0, 2.0, 1), 3.0, 2},
		{New(0, 3.0, 3), 1.5, 2},
		{New(1, 0.0, 0), 0.0, 0},
	}

	for _, test := range tests {
		currentPoint := test.transaction.currentPoint
		nextPoint := test.transaction.nextPoint
		time := test.transaction.time

		test.transaction.Wait(test.timer)

		if newTime := test.transaction.time; newTime != time+test.timer {
			t.Errorf("Expected transaction's time %.2f, got %.2f", time+test.timer, newTime)
		}
		if lifetime := test.transaction.lifetime; lifetime != test.timer {
			t.Errorf("Expected transaction's lifetime %.2f, got %.2f", test.timer, lifetime)
		}
		if newCurrentPoint := test.transaction.currentPoint; newCurrentPoint != currentPoint {
			t.Errorf("Expected %d index of current point in transaction, got %d", newCurrentPoint, currentPoint)
		}
		if newNextPoint := test.transaction.nextPoint; newNextPoint != nextPoint {
			t.Errorf("Expected %d index of next point in transaction, got %d", newNextPoint, nextPoint)
		}
	}
}
