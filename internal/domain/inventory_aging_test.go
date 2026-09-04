package domain

import "testing"

func TestAgingBucket(t *testing.T) {
	tests := []struct {
		days int64
		want InventoryAgingBucket
	}{
		{0, Aging0To30}, {30, Aging0To30}, {31, Aging31To60}, {60, Aging31To60},
		{61, Aging61To90}, {90, Aging61To90}, {91, Aging91To180}, {180, Aging91To180}, {181, Aging181Plus},
	}
	for _, tt := range tests {
		if got := AgingBucket(tt.days); got != tt.want { t.Errorf("AgingBucket(%d)=%q, want %q", tt.days, got, tt.want) }
	}
}
