package repository

import (
	"testing"
	"time"
)

func TestSaveAndGet(t *testing.T) {
	r := NewMemoryRepo()
	userID := 1

	// save one URL (no expiry)
	r.Save("https://a.com", "abc", userID, nil)

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
		r.Save("https://a.com", "a", userID, nil)
	}
	for i := 0; i < 2; i++ {
		r.Save("https://b.com", "b", userID, nil)
	}

	top := r.GetTopDomains(userID, 3)
	if top["a.com"] != 3 || top["b.com"] != 2 {
		t.Fatalf("unexpected domain counts %v", top)
	}
}

// ✅ Bonus test for expiry behavior (optional)
func TestSaveWithExpiry(t *testing.T) {
	r := NewMemoryRepo()
	userID := 99

	exp := time.Now().Add(2 * time.Hour)
	r.Save("https://temp.com", "t123", userID, &exp)

	u, ok := r.GetURL("t123")
	if !ok || u != "https://temp.com" {
		t.Fatalf("expected active link https://temp.com, got %s", u)
	}

	// simulate expired link
	past := time.Now().Add(-2 * time.Hour)
	r.Save("https://expired.com", "e123", userID, &past)

	if _, ok := r.GetURL("e123"); ok {
		t.Fatalf("expected expired link to be inaccessible")
	}
}
