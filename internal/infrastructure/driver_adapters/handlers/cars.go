package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers/dtos"
)

type Cars struct {
	CarsService ports.CarsService
}

func NewCars(cs ports.CarsService) *Cars {
	return &Cars{
		CarsService: cs,
	}
}

func (ch *Cars) Register(w http.ResponseWriter, r *http.Request) {
	var newCar domain.Car
	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if newCar, err = ch.CarsService.Register(r.Context(), car.ToDomain()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	car.FromDomain(newCar)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
