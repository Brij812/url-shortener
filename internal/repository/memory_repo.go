package repository

import (
	"net/url"
	"sort"
	"sync"
)

type MemoryRepo struct {
	URLToCode   map[string]string
	CodeToURL   map[string]string
	DomainCount map[string]int
	mu          sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		URLToCode:   make(map[string]string),
		CodeToURL:   make(map[string]string),
		DomainCount: make(map[string]int),
	}
}

func (r *MemoryRepo) GetCode(u string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.URLToCode[u]
	return c, ok
}

func (r *MemoryRepo) GetURL(code string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.CodeToURL[code]
	return u, ok
}

func (r *MemoryRepo) Save(u, code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.URLToCode[u] = code
	r.CodeToURL[code] = u
	d, err := url.Parse(u)
	if err == nil && d.Host != "" {
		r.DomainCount[d.Host] = r.DomainCount[d.Host] + 1
	}
}

func (r *MemoryRepo) GetTopDomains(n int) map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	type kv struct {
		Key   string
		Value int
	}
	var arr []kv
	for k, v := range r.DomainCount {
		arr = append(arr, kv{k, v})
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Value > arr[j].Value
	})
	if len(arr) > n {
		arr = arr[:n]
	}
	out := make(map[string]int)
	for _, kv := range arr {
		out[kv.Key] = kv.Value
	}
	return out
}

func (r *MemoryRepo) IncrementDomainCount(u string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, err := url.Parse(u)
	if err == nil && d.Host != "" {
		r.DomainCount[d.Host]++
	}
}
