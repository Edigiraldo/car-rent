package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrReservationNotFound = "reservation was not found"
)

type Reservations struct {
	reservationsRepository ports.ReservationsRepo
}

func NewReservations(rr ports.ReservationsRepo) *Reservations {
	return &Reservations{
		reservationsRepository: rr,
	}
}

func (rs Reservations) Book(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error) {
	reservation.ID = uuid.New()

	if err := rs.reservationsRepository.Insert(ctx, reservation); err != nil {
		return domain.Reservation{}, err
	}

	return reservation, nil
}

func (rs Reservations) Get(ctx context.Context, ID uuid.UUID) (domain.Reservation, error) {
	dc, err := rs.reservationsRepository.Get(ctx, ID)
	if err != nil {
		return domain.Reservation{}, err
	}

	return dc, nil
}

func (rs Reservations) FullUpdate(ctx context.Context, reservation domain.Reservation) error {
	return rs.reservationsRepository.FullUpdate(ctx, reservation)
}
