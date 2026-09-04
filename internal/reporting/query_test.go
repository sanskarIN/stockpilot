package reporting

import (
	"testing"
	"time"
)

func TestQueryCarriesBoundsAndPeriod(t *testing.T) {
	from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 4, 23, 59, 59, 0, time.UTC)
	q := Query{From: &from, To: &to, Limit: 250, Offset: 500}
	if q.From == nil || !q.From.Equal(from) || q.To == nil || !q.To.Equal(to) {
		t.Fatalf("unexpected period: %#v", q)
	}
	if q.Limit != 250 || q.Offset != 500 {
		t.Fatalf("unexpected bounds: %#v", q)
	}
}

func TestQueryAllowsNonTimeBasedReports(t *testing.T) {
	q := Query{Limit: 1000, Offset: 0}
	if q.From != nil || q.To != nil {
		t.Fatalf("expected optional period to be nil: %#v", q)
	}
}
