package repository

import "testing"

func TestSaveAndGet(t *testing.T) {
	r := NewMemoryRepo()
	userID := 1

	// save one URL
	r.Save("https://a.com", "abc", userID)

	// verify GetCode
	if c, ok := r.GetCode("https://a.com"); !ok || c != "abc" {
		t.Fatalf("expected code abc, got %s", c)
	}

	// verify GetURL
	if u, ok := r.GetURL("abc"); !ok || u != "https://a.com" {
		t.Fatalf("expected url https://a.com, got %s", u)
	}

	// verify GetAllURLsByUser
	urls := r.GetAllURLsByUser(userID)
	if len(urls) != 1 {
		t.Fatalf("expected 1 url for user, got %d", len(urls))
	}
	if urls[0]["long_url"] != "https://a.com" {
		t.Fatalf("expected long_url=https://a.com, got %v", urls[0])
	}
}

func TestDomainCount(t *testing.T) {
	r := NewMemoryRepo()
	userID := 42

	for i := 0; i < 3; i++ {
		r.Save("https://a.com", "a", userID)
	}
	for i := 0; i < 2; i++ {
		r.Save("https://b.com", "b", userID)
	}

	top := r.GetTopDomains(3)
	if top["a.com"] != 3 || top["b.com"] != 2 {
		t.Fatalf("unexpected domain counts %v", top)
	}
}
