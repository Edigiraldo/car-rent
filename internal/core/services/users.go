package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = "user not found"
)

type Users struct {
	usersRepository ports.UsersRepo
}

func NewUsers(ur ports.UsersRepo) *Users {
	return &Users{
		usersRepository: ur,
	}
}

func (us *Users) Register(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()

	if err := us.usersRepository.Insert(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us *Users) Get(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	du, err := us.usersRepository.Get(ctx, ID)
	if err != nil {
		return domain.User{}, err
	}

	return du, nil
}
