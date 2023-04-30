package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/Edigiraldo/car-rent/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrReservationNotFound         = "reservation was not found"
	ErrInvalidReservationTimeFrame = "reservation time frame is invalid"
	ErrMinimumReservationHours     = "period is shorter than minimun allowed"
)

type Reservations struct {
	reservationsRepository ports.ReservationsRepo
}

func NewReservations(rr ports.ReservationsRepo) Reservations {
	return Reservations{
		reservationsRepository: rr,
	}
}

func (rs Reservations) Book(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error) {
	if err := rs.CheckReservation(ctx, reservation); err != nil {
		return domain.Reservation{}, err
	}

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
	if err := rs.CheckReservation(ctx, reservation); err != nil {
		return err
	}

	return rs.reservationsRepository.FullUpdate(ctx, reservation)
}

func (rs Reservations) Delete(ctx context.Context, id uuid.UUID) error {
	return rs.reservationsRepository.Delete(ctx, id)
}

func (rs Reservations) GetByCarID(ctx context.Context, carID uuid.UUID) ([]domain.Reservation, error) {
	drs, err := rs.reservationsRepository.GetByCarID(ctx, carID)
	if err != nil {
		return nil, err
	}

	return drs, nil
}

func (rs Reservations) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Reservation, error) {
	drs, err := rs.reservationsRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return drs, nil
}

func (rs Reservations) CheckReservation(ctx context.Context, reservation domain.Reservation) error {
	if isValid := utils.IsValidTimeFrame(reservation.StartDate, reservation.EndDate); !isValid {
		return errors.New(ErrInvalidReservationTimeFrame)
	}

	if reservation.StartDate.Before(time.Now()) {
		return errors.New(ErrInvalidReservationTimeFrame)
	}

	if reservation.EndDate.Sub(reservation.StartDate).Hours() < float64(constants.Values.MINIMUM_RESERVATION_HOURS) {
		return fmt.Errorf("%s (%d hours)", ErrMinimumReservationHours, constants.Values.MINIMUM_RESERVATION_HOURS)
	}

	reservations, err := rs.reservationsRepository.GetByCarIDAndTimeFrame(ctx, reservation.CarID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return err
	}

	if len(reservations) > 0 {
		return errors.New(ErrCarNotAvailable)
	}

	return nil

}
