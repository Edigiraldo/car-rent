package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidID           = "id could not be converted to uuid"
	ErrInternalServerError = "internal server error"
	ErrCityQueryParamEmpty = "city query param can not be empty"
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
		if err.Error() == services.ErrInvalidCityName {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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
		} else if err.Error() == services.ErrInvalidCityName {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

// Lists cars from a city in pages of 20 elements. from_car_id parameter
// is taken as the last seen car in a previous page.
func (ch *Cars) List(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, ErrCityQueryParamEmpty, http.StatusBadRequest)

		return
	}
	from_car_id := r.URL.Query().Get("from_car_id")

	cars, err := ch.CarsService.List(r.Context(), city, from_car_id)
	if err != nil && err.Error() != services.ErrInvalidCityName {
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		log.Println(err)

		return
	}

	listCarsResponse := getListCarsResponse(cars)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listCarsResponse)
}

// Gets a list of cars and builds the user response
func getListCarsResponse(cars []domain.Car) (listCarsResponse dtos.ListCarsResponse) {
	listCarsResponse.Cars = make([]dtos.Car, 0)
	for _, domainCar := range cars {
		car := dtos.Car{}
		car.FromDomain(domainCar)

		listCarsResponse.Cars = append(listCarsResponse.Cars, car)
	}

	return listCarsResponse
}
