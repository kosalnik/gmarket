package service

type AccountService struct {
}

func NewAccountService() (*AccountService, error) {
	return &AccountService{}, nil
}
