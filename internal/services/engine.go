package services

import (
	"auto-service/internal/clients"
	"time"

	"github.com/TB-Systems/financial-control-backend-commons/dtos"
	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/google/uuid"
)

type Engine interface {
	GetMonthlyTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	GetAnnualTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	GetInstallmentsTransactions() (commonsmodels.ResponseList[uuid.UUID], error)
	PostMonthlyTransaction(uuid.UUID) error
	PostAnnualTransaction(uuid.UUID) error
	PostInstallmentsTransaction(uuid.UUID) error
}

type engine struct {
	client clients.Engine
}

func NewEngineService(client clients.Engine) Engine {
	return &engine{
		client: client,
	}
}

func (e *engine) GetMonthlyTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return e.client.GetMonthlyTransactions()
}

func (e *engine) GetAnnualTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return e.client.GetAnnualTransactions()
}

func (e *engine) GetInstallmentsTransactions() (commonsmodels.ResponseList[uuid.UUID], error) {
	return e.client.GetInstallmentsTransactions()
}

func (e *engine) PostMonthlyTransaction(id uuid.UUID) error {
	param := buildTransactionRequestFromRecurrentTransaction(id)
	return e.client.PostMonthlyTransaction(param)
}

func (e *engine) PostAnnualTransaction(id uuid.UUID) error {
	param := buildTransactionRequestFromRecurrentTransaction(id)
	return e.client.PostAnnualTransaction(param)
}

func (e *engine) PostInstallmentsTransaction(id uuid.UUID) error {
	param := buildTransactionRequestFromRecurrentTransaction(id)
	return e.client.PostInstallmentsTransaction(param)
}

func buildTransactionRequestFromRecurrentTransaction(id uuid.UUID) dtos.TransactionRequestFromRecurrentTransaction {
	now := time.Now()

	return dtos.TransactionRequestFromRecurrentTransaction{
		ID:    id,
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}
}
