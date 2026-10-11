package idgen

import (
	"strings"
	"testing"
)

func TestNewAddsPrefixAndEntropy(t *testing.T) {
	first, err := New("prd")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	second, err := New("prd")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if !strings.HasPrefix(first, "prd_") {
		t.Fatalf("New() = %q, want prd_ prefix", first)
	}
	if len(first) != len("prd_")+32 {
		t.Fatalf("New() length = %d, want %d", len(first), len("prd_")+32)
	}
	if first == second {
		t.Fatalf("New() produced duplicate ids %q", first)
	}
}

func TestNewWithoutPrefixReturnsHexIdentifier(t *testing.T) {
	id, err := New("")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	if len(id) != 32 {
		t.Fatalf("New() length = %d, want 32", len(id))
	}
	for _, char := range id {
		if !strings.ContainsRune("0123456789abcdef", char) {
			t.Fatalf("New() contains non-hex character %q", char)
		}
	}
}
