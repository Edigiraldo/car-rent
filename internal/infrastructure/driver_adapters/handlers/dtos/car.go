package dtos

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidCarType   = "invalid car type"
	ErrInvalidCarStatus = "invalid car status"
)

type Car struct {
	ID             uuid.UUID `json:"id"`
	Type           CarType   `json:"type"`
	Seats          int16     `json:"seats"`
	HourlyRentCost float64   `json:"hourly_rent_cost"`
	City           string    `json:"city"`
	Status         CarStatus `json:"status"`
}

type CarType string

const (
	Sedan     CarType = "Sedan"
	Luxury    CarType = "Luxury"
	SportsCar CarType = "Sports Car"
	Limousine CarType = "Limousine"
)

type CarStatus string

const (
	Available   CarStatus = "Available"
	Unavailable CarStatus = "Unavailable"
)

func (c Car) ToDomain() domain.Car {
	return domain.Car{
		ID:             c.ID,
		Type:           domain.CarType(c.Type),
		Seats:          c.Seats,
		HourlyRentCost: c.HourlyRentCost,
		City:           c.City,
		Status:         domain.CarStatus(c.Status),
	}
}

func CarFromBody(body io.Reader) (Car, error) {
	var car Car
	err := json.NewDecoder(body).Decode(&car)
	if err != nil {
		return Car{}, err
	}

	if !isValidCarType(car.Type) {
		return Car{}, errors.New(ErrInvalidCarType)
	}

	if !isValidCarStatus(car.Status) {
		return Car{}, errors.New(ErrInvalidCarStatus)
	}

	return car, nil
}

func isValidCarType(carType CarType) bool {
	switch carType {
	case Sedan, Luxury, SportsCar, Limousine:
		return true
	}
	return false
}

func isValidCarStatus(carStatus CarStatus) bool {
	switch carStatus {
	case Available, Unavailable:
		return true
	}
	return false
}
