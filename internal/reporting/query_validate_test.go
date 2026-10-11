package reporting

import (
	"testing"
	"time"
)

func TestQueryValidateRejectsInvalidBounds(t *testing.T) {
	q := Query{Limit: -1}
	if err := q.Validate(); err == nil {
		t.Fatal("expected negative limit to fail")
	}
	q = Query{Offset: -1}
	if err := q.Validate(); err == nil {
		t.Fatal("expected negative offset to fail")
	}
}

func TestQueryValidateRejectsReversedPeriod(t *testing.T) {
	from := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)
	if err := (Query{From: &from, To: &to}).Validate(); err == nil {
		t.Fatal("expected reversed period to fail")
	}
}
