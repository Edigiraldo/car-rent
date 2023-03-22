package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/handlers/dtos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidID           = "id could not be converted to uuid"
	ErrInternalServerError = "internal server error"
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
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)

		return
	}

	car.FromDomain(newCar)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

func (ch *Cars) Get(w http.ResponseWriter, r *http.Request) {
	var car dtos.Car

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	dc, err := ch.CarsService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrCarNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	car.FromDomain(dc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

func (ch *Cars) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// Get the ID from path param
	car.ID = ID

	if err = ch.CarsService.FullUpdate(r.Context(), car.ToDomain()); err != nil {
		if err.Error() == services.ErrCarNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
			log.Println(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}

func (ch *Cars) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, ErrInvalidID, http.StatusBadRequest)

		return
	}

	err = ch.CarsService.Delete(r.Context(), ID)
	if err != nil {
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		log.Println(err)

		return
	}
	w.WriteHeader(http.StatusNoContent)
}
