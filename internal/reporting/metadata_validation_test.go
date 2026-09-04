package reporting

import (
	"testing"
	"time"
)

func TestNewMetadataPreservesBounds(t *testing.T) {
	period, err := NewPeriod(
		time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	bounds, err := NewBounds(1, 0, 30, 5000)
	if err != nil {
		t.Fatal(err)
	}
	m := NewMetadata(period, bounds, true, time.Time{})
	if m.Limit != 1 || m.Offset != 0 || !m.Complete {
		t.Fatalf("unexpected metadata: %#v", m)
	}
	if m.GeneratedAt.IsZero() {
		t.Fatal("generatedAt must be populated when omitted")
	}
}
