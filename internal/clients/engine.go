package clients

import (
	"github.com/TB-Systems/financial-control-backend-commons/dtos"
	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/google/uuid"
)

type Engine interface {
	GetMonthlyTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	GetAnnualTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	GetInstallmentsTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	PostMonthlyTransaction(dtos.TransactionRequestFromRecurrentTransaction) error
	PostAnnualTransaction(dtos.TransactionRequestFromRecurrentTransaction) error
	PostInstallmentsTransaction(dtos.TransactionRequestFromRecurrentTransaction) error
}

type engine struct {
}

func NewEngineClient() Engine {
	return &engine{}
}

func (e *engine) GetMonthlyTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return commonsmodels.ResponseList[uuid.UUID]{
		Items: []uuid.UUID{
			uuid.MustParse("104baf42-53d4-4264-a565-e6ecfc91f39f"),
			uuid.MustParse("e4696c5c-81a4-4d89-88d7-d9a3639165e6"),
		},
		Total: 2,
	}, nil
}

func (e *engine) GetAnnualTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return commonsmodels.ResponseList[uuid.UUID]{
		Items: []uuid.UUID{
			uuid.MustParse("25271f83-2c07-44db-bf61-7412c811d6f1"),
			uuid.MustParse("6019e1ef-f04d-436f-9e25-651913841a37"),
		},
		Total: 2,
	}, nil
}

func (e *engine) GetInstallmentsTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return commonsmodels.ResponseList[uuid.UUID]{
		Items: []uuid.UUID{
			uuid.MustParse("87382f6f-528b-4cdf-bf2a-f7fb74acde41"),
			uuid.MustParse("2448c9fb-a0a1-4390-bca7-b90bd298f51b"),
		},
		Total: 2,
	}, nil
}

func (e *engine) PostMonthlyTransaction(dtos.TransactionRequestFromRecurrentTransaction) error {
	return nil
}

func (e *engine) PostAnnualTransaction(dtos.TransactionRequestFromRecurrentTransaction) error {
	return nil
}

func (e *engine) PostInstallmentsTransaction(dtos.TransactionRequestFromRecurrentTransaction) error {
	return nil
}
