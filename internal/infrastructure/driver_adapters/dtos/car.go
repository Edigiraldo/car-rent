package dtos

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	"github.com/Edigiraldo/car-rent/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidSeatsNumber    = "seats number must be greater than 0"
	ErrInvalidHourlyRentCost = "hourly rent cost must be greater than 0"
	ErrEmptyCity             = "city name cannot be empty"
	ErrInvalidCarType        = "invalid car type"
	ErrInvalidCarStatus      = "invalid car status"
)

type ListCarsResponse struct {
	Cars []Car `json:"cars"`
}

type Car struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Type           string    `json:"type"`
	Seats          int16     `json:"seats"`
	HourlyRentCost float64   `json:"hourly_rent_cost"`
	CityName       string    `json:"city_name"`
	Status         string    `json:"status"`
}

func (c Car) ToDomain() domain.Car {
	return domain.Car{
		ID:             c.ID,
		Type:           c.Type,
		Seats:          c.Seats,
		HourlyRentCost: c.HourlyRentCost,
		CityName:       c.CityName,
		Status:         c.Status,
	}
}

func (c *Car) FromDomain(dc domain.Car) {
	c.ID = dc.ID
	c.Type = dc.Type
	c.Seats = dc.Seats
	c.HourlyRentCost = dc.HourlyRentCost
	c.CityName = dc.CityName
	c.Status = dc.Status
}

func CarFromBody(body io.Reader) (Car, error) {
	var car Car
	err := json.NewDecoder(body).Decode(&car)
	if err != nil {
		return Car{}, err
	}

	if car.Seats <= 0 {
		return Car{}, errors.New(ErrInvalidSeatsNumber)
	}

	if car.HourlyRentCost <= 0 {
		return Car{}, errors.New(ErrInvalidHourlyRentCost)
	}

	if car.CityName == "" {
		return Car{}, errors.New(ErrEmptyCity)
	}

	if !isValidCarType(car.Type) {
		return Car{}, errors.New(ErrInvalidCarType)
	}

	if !isValidCarStatus(car.Status) {
		return Car{}, errors.New(ErrInvalidCarStatus)
	}

	return car, nil
}

func isValidCarType(carType string) bool {
	carTypes := constants.Values.CAR_TYPES.Values()

	return utils.IsInSlice(carTypes, carType)
}

func isValidCarStatus(carStatus string) bool {
	carStatuses := constants.Values.CAR_STATUSES.Values()

	return utils.IsInSlice(carStatuses, carStatus)
}
