package services

type RecurringTransactionsService struct {
}

func NewRecurringTransactionsService() *RecurringTransactionsService {
	return &RecurringTransactionsService{}
}

func (s *RecurringTransactionsService) Run() error {
	return nil
}
