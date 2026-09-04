package reporting

import (
	"testing"
	"time"
)

func TestNewPeriodRejectsMissingBounds(t *testing.T) {
	if _, err := NewPeriod(time.Time{}, time.Now()); err == nil {
		t.Fatal("expected missing bound error")
	}
}

func TestNewPeriodRejectsReversedBounds(t *testing.T) {
	from := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, -1)
	if _, err := NewPeriod(from, to); err == nil {
		t.Fatal("expected reversed period error")
	}
}

func TestNewPeriodBoundsAndPrevious(t *testing.T) {
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	period, err := NewPeriod(from, to)
	if err != nil {
		t.Fatalf("NewPeriod() error = %v", err)
	}
	if got := period.Days(); got != 30 {
		t.Fatalf("Days() = %d, want 30", got)
	}
	previous := period.Previous()
	if !previous.From.Equal(time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)) || !previous.To.Equal(time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected previous period: %#v", previous)
	}
}

func TestNewPeriodRejectsOversizedWindow(t *testing.T) {
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 365)
	if _, err := NewPeriod(from, to); err == nil {
		t.Fatal("expected oversized period error")
	}
}
