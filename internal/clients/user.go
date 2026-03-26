package clients

import (
	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/google/uuid"
)

type User interface {
	GetPageCount() (int, error)
	GetUsersPaginated() (commonsmodels.PaginatedResponse[uuid.UUID], error)
}

type user struct {
}

func NewUserClient() User {
	return &user{}
}

func (u *user) GetPageCount() (int, error) {
	return 1, nil
}

func (u *user) GetUsersPaginated() (commonsmodels.PaginatedResponse[uuid.UUID], error) {
	return commonsmodels.PaginatedResponse[uuid.UUID]{
		Items: []uuid.UUID{
			uuid.MustParse("c6d00ff1-44d8-489b-b30b-6ad5f466d459"),
			uuid.MustParse("89defd09-f83e-4957-b182-779957b4bb9d"),
			uuid.MustParse("5b3d16da-8dec-4761-82ab-02939f4f2f9d"),
			uuid.MustParse("b475ee33-6801-41f9-802e-e2ed7bdc3309"),
		},
		PageCount: 1,
		Page:      1,
	}, nil
}
