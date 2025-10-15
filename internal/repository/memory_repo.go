package repository

import "sync"

type MemoryRepo struct {
	URLToCode map[string]string
	CodeToURL map[string]string
	mu        sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		URLToCode: make(map[string]string),
		CodeToURL: make(map[string]string),
	}
}

func (r *MemoryRepo) GetCode(url string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	code, exists := r.URLToCode[url]
	return code, exists
}

func (r *MemoryRepo) GetURL(code string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	url, exists := r.CodeToURL[code]
	return url, exists
}

func (r *MemoryRepo) Save(url, code string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.URLToCode[url] = code
	r.CodeToURL[code] = url
}
