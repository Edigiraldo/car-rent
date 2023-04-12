package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
)

type Users struct {
	usersRepository ports.UsersRepo
}

func NewUsers(ur ports.UsersRepo) *Users {
	return &Users{
		usersRepository: ur,
	}
}

func (u *Users) Register(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()

	if err := u.usersRepository.Insert(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
