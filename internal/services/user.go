package services

import (
	"auto-service/internal/clients"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/google/uuid"
)

type User interface {
	GetPageCount() (int, error)
	GetUsersPaginated() (commonsmodels.PaginatedResponse[uuid.UUID], error)
}

type user struct {
	client clients.User
}

func NewUserService(client clients.User) User {
	return &user{
		client: client,
	}
}

func (u *user) GetPageCount() (int, error) {
	return u.client.GetPageCount()
}

func (u *user) GetUsersPaginated() (commonsmodels.PaginatedResponse[uuid.UUID], error) {
	return u.client.GetUsersPaginated()
}
