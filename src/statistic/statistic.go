package statistic

type Unit struct {
	sum   float64
	count int
}

// Add new value.
func (u *Unit) AddValue(Value float64) {
	u.sum += Value
	u.count++
}

// Get mean value.
func (u *Unit) GetMean() float64 {
	if u.count != 0 {
		return u.sum / float64(u.count)
	} else {
		return 0.0
	}
}

// Get summary value.
func (u *Unit) GetSum() float64 {
	return u.sum
}
