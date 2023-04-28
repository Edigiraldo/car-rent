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
	usersRepository        ports.UsersRepo
	reservationsRepository ports.ReservationsRepo
}

func NewUsers(ur ports.UsersRepo, rr ports.ReservationsRepo) *Users {
	return &Users{
		usersRepository:        ur,
		reservationsRepository: rr,
	}
}

func (us Users) Register(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()

	if err := us.usersRepository.Insert(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us Users) Get(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	du, err := us.usersRepository.Get(ctx, ID)
	if err != nil {
		return domain.User{}, err
	}

	return du, nil
}

func (us Users) FullUpdate(ctx context.Context, user domain.User) error {
	return us.usersRepository.FullUpdate(ctx, user)
}

func (us Users) Delete(ctx context.Context, id uuid.UUID) error {
	return us.usersRepository.Delete(ctx, id)
}

func (us Users) GetReservations(ctx context.Context, userID uuid.UUID) ([]domain.Reservation, error) {
	drs, err := us.reservationsRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return drs, nil
}
