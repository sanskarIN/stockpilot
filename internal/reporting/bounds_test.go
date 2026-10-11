package reporting

import "testing"

func TestNewBoundsDefaults(t *testing.T) {
	bounds, err := NewBounds(0, 0, 100, 5000)
	if err != nil {
		t.Fatalf("NewBounds() error = %v", err)
	}
	if bounds.Limit != 100 || bounds.Offset != 0 {
		t.Fatalf("unexpected bounds: %#v", bounds)
	}
}

func TestNewBoundsRejectsInvalidValues(t *testing.T) {
	cases := []struct {
		name   string
		limit  int
		offset int
	}{
		{"negative limit", -1, 0},
		{"oversized limit", 5001, 0},
		{"negative offset", 100, -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewBounds(tc.limit, tc.offset, 100, 5000); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestParsePositiveInt(t *testing.T) {
	if got, err := ParsePositiveInt("42", 10); err != nil || got != 42 {
		t.Fatalf("ParsePositiveInt() = %d, %v; want 42, nil", got, err)
	}
	if got, err := ParsePositiveInt("", 10); err != nil || got != 10 {
		t.Fatalf("fallback = %d, %v; want 10, nil", got, err)
	}
	if _, err := ParsePositiveInt("-2", 10); err == nil {
		t.Fatal("expected negative value error")
	}
	if _, err := ParsePositiveInt("nope", 10); err == nil {
		t.Fatal("expected malformed value error")
	}
}
