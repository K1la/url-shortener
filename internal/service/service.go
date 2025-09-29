package service

type Service struct {
	repo  RepositoryI
	cache CacheI
}

func New(r RepositoryI, c CacheI) *Service {
	return &Service{r, c}
}
