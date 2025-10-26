package repository

type Repository interface {
	GetCode(u string) (string, bool)
	GetURL(code string) (string, bool)
	Save(u, code string)
	GetTopDomains(n int) map[string]int
	IncrementDomainCount(u string)
}
