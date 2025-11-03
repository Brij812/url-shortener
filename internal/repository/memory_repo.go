package repository

import (
	"sync"
	"time"
)

type MemoryRepo struct {
	mu           sync.RWMutex
	urlToCode    map[string]string
	codeToURL    map[string]string
	domainCounts map[int]map[string]int // per-user domain counts
	userLinks    map[int][]map[string]string
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		urlToCode:    make(map[string]string),
		codeToURL:    make(map[string]string),
		domainCounts: make(map[int]map[string]int),
		userLinks:    make(map[int][]map[string]string),
	}
}

// Save a URL–code pair associated with a user
func (r *MemoryRepo) Save(u, code string, userID int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urlToCode[u] = code
	r.codeToURL[code] = u

	if _, ok := r.domainCounts[userID]; !ok {
		r.domainCounts[userID] = make(map[string]int)
	}

	// Track per-user links
	r.userLinks[userID] = append(r.userLinks[userID], map[string]string{
		"short_url":  "http://localhost:8080/" + code,
		"long_url":   u,
		"created_at": time.Now().Format(time.RFC3339),
	})

	// Increment domain count per user
	domain := extractDomain(u)
	if domain != "" {
		r.domainCounts[userID][domain]++
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

// GetTopDomains — per-user
func (r *MemoryRepo) GetTopDomains(userID, n int) map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if domains, ok := r.domainCounts[userID]; ok {
		return domains
	}
	return map[string]int{}
}

// IncrementDomainCount — per-user
func (r *MemoryRepo) IncrementDomainCount(u string, userID int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.domainCounts[userID]; !ok {
		r.domainCounts[userID] = make(map[string]int)
	}

	domain := extractDomain(u)
	if domain != "" {
		r.domainCounts[userID][domain]++
	}
}

// GetAllURLsByUser — unchanged
func (r *MemoryRepo) GetAllURLsByUser(userID int) []map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.userLinks[userID]
}
