package repository

import "testing"

func TestSaveAndGet(t *testing.T) {
	r := NewMemoryRepo()
	r.Save("https://a.com", "abc")
	if c, ok := r.GetCode("https://a.com"); !ok || c != "abc" {
		t.Fatalf("expected code abc, got %s", c)
	}
	if u, ok := r.GetURL("abc"); !ok || u != "https://a.com" {
		t.Fatalf("expected url https://a.com, got %s", u)
	}
}

func TestDomainCount(t *testing.T) {
	r := NewMemoryRepo()
	for i := 0; i < 3; i++ {
		r.Save("https://a.com", "a")
	}
	for i := 0; i < 2; i++ {
		r.Save("https://b.com", "b")
	}
	top := r.GetTopDomains(3)
	if top["a.com"] != 3 || top["b.com"] != 2 {
		t.Fatalf("unexpected domain counts %v", top)
	}
}
