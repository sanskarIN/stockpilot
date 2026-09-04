package reporting

import (
	"fmt"
	"time"
)

const (
	DefaultPeriodDays = 30
	MinPeriodDays     = 1
	MaxPeriodDays     = 365
)

// Period describes a bounded inclusive reporting window.
type Period struct {
	From time.Time
	To   time.Time
}

func NewPeriod(from, to time.Time) (Period, error) {
	if from.IsZero() || to.IsZero() {
		return Period{}, fmt.Errorf("reporting period requires both from and to")
	}
	if to.Before(from) {
		return Period{}, fmt.Errorf("reporting period end must not precede start")
	}
	days := int(to.Sub(from).Hours()/24) + 1
	if days < MinPeriodDays || days > MaxPeriodDays {
		return Period{}, fmt.Errorf("reporting period must contain %d-%d days", MinPeriodDays, MaxPeriodDays)
	}
	return Period{From: from, To: to}, nil
}

func (p Period) Days() int {
	if p.From.IsZero() || p.To.IsZero() || p.To.Before(p.From) {
		return 0
	}
	return int(p.To.Sub(p.From).Hours()/24) + 1
}

func (p Period) Previous() Period {
	days := p.Days()
	if days == 0 {
		return Period{}
	}
	end := p.From.AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -(days - 1))
	return Period{From: start, To: end}
}
