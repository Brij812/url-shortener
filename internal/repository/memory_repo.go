package repository

import (
	"net/url"
	"sync"
	"time"
)

type MemoryRepo struct {
	mu           sync.RWMutex
	urlToCode    map[string]string
	codeToURL    map[string]string
	domainCounts map[string]int
	userLinks    map[int][]map[string]string
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		urlToCode:    make(map[string]string),
		codeToURL:    make(map[string]string),
		domainCounts: make(map[string]int),
		userLinks:    make(map[int][]map[string]string),
	}
}

// Save a URL–code pair associated with a user
func (r *MemoryRepo) Save(u, code string, userID int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urlToCode[u] = code
	r.codeToURL[code] = u

	// Track per-user links
	r.userLinks[userID] = append(r.userLinks[userID], map[string]string{
		"short_url":  "http://localhost:8080/" + code,
		"long_url":   u,
		"created_at": time.Now().Format(time.RFC3339),
	})

	// Increment domain count
	d, err := url.Parse(u)
	if err == nil && d.Host != "" {
		r.domainCounts[d.Host]++
	}
}

func (r *MemoryRepo) GetCode(u string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	code, ok := r.urlToCode[u]
	return code, ok
}

func (r *MemoryRepo) GetURL(code string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.codeToURL[code]
	return u, ok
}

func (r *MemoryRepo) GetTopDomains(n int) map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.domainCounts
}

func (r *MemoryRepo) IncrementDomainCount(u string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, err := url.Parse(u)
	if err == nil && d.Host != "" {
		r.domainCounts[d.Host]++
	}
}

// new method added — satisfies Repository interface
func (r *MemoryRepo) GetAllURLsByUser(userID int) []map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.userLinks[userID]
}
