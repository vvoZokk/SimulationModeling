package statistic

import "testing"

type testPair struct {
	unit   *Unit
	result float64
}

func TestMean(t *testing.T) {
	tests := []testPair{
		{new(3.0, 2), 1.5},
		{new(-3.0, 3), -1.0},
		{new(3.0, 0), 0.0},
		{new(0.0, 5), 0.0},
	}

	for _, test := range tests {
		if r := test.unit.Mean(); r != test.result {
			t.Errorf("Expected %.3f, got %.3f", test.result, r)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []testPair{
		{new(3.0, 2), 3.0},
		{new(-3.0, 3), -3.0},
		{new(3.0, 0), 3.0},
		{new(0.0, 5), 0.0},
	}

	for _, test := range tests {
		if r := test.unit.Sum(); r != test.result {
			t.Errorf("Expected %.3f, got %.3f", test.result, r)
		}
	}
}
