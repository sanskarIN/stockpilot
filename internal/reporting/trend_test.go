package reporting

import "testing"

func TestCompareIncrease(t *testing.T) {
	change := Compare(120, 100)
	if change.Delta != 20 {
		t.Fatalf("Delta = %v, want 20", change.Delta)
	}
	if change.Percent == nil || *change.Percent != 20 {
		t.Fatalf("Percent = %v, want 20", change.Percent)
	}
}

func TestCompareDecreaseUsesAbsoluteBaseline(t *testing.T) {
	change := Compare(80, 100)
	if change.Percent == nil || *change.Percent != -20 {
		t.Fatalf("Percent = %v, want -20", change.Percent)
	}
}

func TestCompareZeroBaselineLeavesPercentUndefined(t *testing.T) {
	change := Compare(10, 0)
	if change.Percent != nil {
		t.Fatalf("Percent = %v, want nil", change.Percent)
	}
}

func TestCompareNegativeBaselineIsStable(t *testing.T) {
	change := Compare(-5, -10)
	if change.Percent == nil || *change.Percent != 50 {
		t.Fatalf("Percent = %v, want 50", change.Percent)
	}
}
