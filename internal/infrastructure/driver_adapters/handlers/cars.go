package handlers

import (
	"log"
	"net/http"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/ports"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	"github.com/Edigiraldo/car-rent/pkg/httphandler"
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

func NewCars(cs ports.CarsService) Cars {
	return Cars{
		CarsService: cs,
	}
}

func (ch Cars) Register(w http.ResponseWriter, r *http.Request) {
	var newCar domain.Car
	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if newCar, err = ch.CarsService.Register(r.Context(), car.ToDomain()); err != nil {
		if err.Error() == services.ErrInvalidCityName {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Println(err)
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)

		return
	}

	car.FromDomain(newCar)

	httphandler.WriteSuccessResponse(w, http.StatusCreated, car)
}

func (ch Cars) Get(w http.ResponseWriter, r *http.Request) {
	var car dtos.Car

	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	dc, err := ch.CarsService.Get(r.Context(), ID)
	if err != nil {
		if err.Error() == services.ErrCarNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	car.FromDomain(dc)

	httphandler.WriteSuccessResponse(w, http.StatusOK, car)
}

func (ch Cars) FullUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	car, err := dtos.CarFromBody(r.Body)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the ID from path param
	car.ID = ID

	if err = ch.CarsService.FullUpdate(r.Context(), car.ToDomain()); err != nil {
		if err.Error() == services.ErrCarNotFound {
			httphandler.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == services.ErrInvalidCityName {
			httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
			log.Println(err)
		}

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusOK, car)
}

func (ch Cars) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := uuid.Parse(id)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrInvalidID)
		return
	}

	err = ch.CarsService.Delete(r.Context(), ID)
	if err != nil {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	httphandler.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

// Lists cars from a city in pages of 20 elements. from_car_id parameter
// is taken as the last seen car in a previous page.
func (ch Cars) List(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, ErrCityQueryParamEmpty)
		return
	}
	from_car_id := r.URL.Query().Get("from_car_id")
	if _, err := uuid.Parse(from_car_id); err != nil && from_car_id != "" {
		httphandler.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cars, err := ch.CarsService.List(r.Context(), city, from_car_id)
	if err != nil && err.Error() != services.ErrInvalidCityName {
		httphandler.WriteErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		log.Println(err)

		return
	}

	listCarsResponse := getListCarsResponse(cars)

	httphandler.WriteSuccessResponse(w, http.StatusOK, listCarsResponse)
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
