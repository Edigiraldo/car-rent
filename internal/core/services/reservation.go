package services

import (
	"context"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/google/uuid"
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
