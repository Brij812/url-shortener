package repository

type Repository interface {
	Save(u, code string, userID int)
	GetCode(u string) (string, bool)
	GetURL(code string) (string, bool)
	GetTopDomains(userID, n int) map[string]int
	IncrementDomainCount(u string, userID int)
	GetAllURLsByUser(userID int) []map[string]string
}
