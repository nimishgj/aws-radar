package collector

import "testing"

func TestSplitKey(t *testing.T) {
	parts := splitKey("a|b|c", 3)
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
	if parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Fatalf("unexpected parts: %v", parts)
	}
}
