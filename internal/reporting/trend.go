package reporting

import "math"

// Change describes a period-over-period change. Percent is nil when the
// previous value is zero, because an ordinary percentage change is undefined.
type Change struct {
	Current  float64  `json:"current"`
	Previous float64  `json:"previous"`
	Delta    float64  `json:"delta"`
	Percent  *float64 `json:"percent,omitempty"`
}

func Compare(current, previous float64) Change {
	change := Change{Current: current, Previous: previous, Delta: current - previous}
	if previous == 0 {
		return change
	}
	percent := ((current - previous) / math.Abs(previous)) * 100
	change.Percent = &percent
	return change
}
