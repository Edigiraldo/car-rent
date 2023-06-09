package docs

import "github.com/google/uuid"

type CarRequest struct {
	Type           string  `json:"type" example:"Luxury"`
	Seats          int16   `json:"seats" example:"4"`
	HourlyRentCost float64 `json:"hourly_rent_cost" example:"99.99"`
	CityName       string  `json:"city_name" example:"New York"`
	Status         string  `json:"status" example:"Available"`
}

type ListCarsResponse struct {
	Cars []CarResponse `json:"cars"`
}

type CarResponse struct {
	ID             uuid.UUID `json:"id,omitempty" example:"bdaf243e-b4d3-49d7-8be4-5ed1fb4dba0e"`
	Type           string    `json:"type" example:"Luxury"`
	Seats          int16     `json:"seats" example:"4"`
	HourlyRentCost float64   `json:"hourly_rent_cost" example:"99.99"`
	CityName       string    `json:"city_name" example:"New York"`
	Status         string    `json:"status" example:"Available"`
}
