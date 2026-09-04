package reporting

import (
	"testing"
	"time"
)

func TestNewMetadata(t *testing.T) {
	from := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	period, err := NewPeriod(from, to)
	if err != nil {
		t.Fatal(err)
	}
	bounds, err := NewBounds(25, 10, 50, 5000)
	if err != nil {
		t.Fatal(err)
	}
	generated := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	metadata := NewMetadata(period, bounds, false, generated)
	if !metadata.GeneratedAt.Equal(generated) {
		t.Fatalf("unexpected generatedAt: %v", metadata.GeneratedAt)
	}
	if !metadata.From.Equal(from) || !metadata.To.Equal(to) {
		t.Fatalf("unexpected period: %#v", metadata)
	}
	if metadata.Limit != 25 || metadata.Offset != 10 || metadata.Complete {
		t.Fatalf("unexpected bounds/completeness: %#v", metadata)
	}
}
