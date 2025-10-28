package repository

type Repository interface {
	GetCode(u string) (string, bool)
	GetURL(code string) (string, bool)
	Save(u, code string, userID int)
	GetTopDomains(n int) map[string]int
	IncrementDomainCount(u string)
	GetAllURLsByUser(userID int) []map[string]string
}
