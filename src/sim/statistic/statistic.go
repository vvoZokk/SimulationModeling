// Package statistic implements small unit for gathering mean and summary value.
package statistic

// Statistic unit.
type Unit struct {
	sum   float64
	count int
}

func new(sum float64, count int) *Unit {
	return &Unit{sum, count}
}

// AddValue adds new value of unit.
func (u *Unit) AddValue(v float64) {
	u.sum += v
	u.count++
}

// Mean returns mean value.
func (u *Unit) Mean() float64 {
	if u.count != 0 {
		return u.sum / float64(u.count)
	} else {
		return 0.0
	}
}

// Sum returns summary value.
func (u *Unit) Sum() float64 {
	return u.sum
}
