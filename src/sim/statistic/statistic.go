// Package statistic implements small unit for gathering mean and summary value.
package statistic

// Statistic unit.
type Unit struct {
	sum   float64
	count int
}

// AddValue adds new value of unit.
func (u *Unit) AddValue(Value float64) {
	u.sum += Value
	u.count++
}

// GetMean returns mean value.
func (u *Unit) GetMean() float64 {
	if u.count != 0 {
		return u.sum / float64(u.count)
	} else {
		return 0.0
	}
}

// GetSum returns summary value.
func (u *Unit) GetSum() float64 {
	return u.sum
}
