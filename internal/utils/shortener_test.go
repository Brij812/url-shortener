package utils

import "testing"

func TestGenerateShortCode(t *testing.T) {
	a := GenerateShortCode("https://example.com")
	b := GenerateShortCode("https://example.com")
	if a != b {
		t.Fatalf("expected deterministic short code, got %s and %s", a, b)
	}
	if len(a) == 0 {
		t.Fatal("expected non-empty short code")
	}
}
