package domain

import "testing"

func TestAgingBucket(t *testing.T) {
	tests := []struct {
		name string
		days int64
		want InventoryAgingBucket
	}{
		{"0 days", 0, Aging0To30},
		{"30 days", 30, Aging0To30},
		{"31 days", 31, Aging31To60},
		{"60 days", 60, Aging31To60},
		{"61 days", 61, Aging61To90},
		{"90 days", 90, Aging61To90},
		{"91 days", 91, Aging91To180},
		{"180 days", 180, Aging91To180},
		{"181 days", 181, Aging181Plus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AgingBucket(tt.days); got != tt.want {
				t.Errorf("AgingBucket(%d)=%q, want %q", tt.days, got, tt.want)
			}
		})
	}
}
