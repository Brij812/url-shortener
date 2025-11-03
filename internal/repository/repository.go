package repository

import "time"

type Repository interface {
	Save(u, code string, userID int, expiresAt *time.Time)
	GetCode(u string) (string, bool)
	GetURL(code string) (string, bool)
	GetTopDomains(userID, n int) map[string]int
	IncrementDomainCount(u string, userID int)
	GetAllURLsByUser(userID int) []map[string]string
}
