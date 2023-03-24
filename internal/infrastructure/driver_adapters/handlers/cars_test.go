package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type carsDependencies struct {
	carsService *mocks.MockCarsService
}

func NewCarsDependencies(carsSrv *mocks.MockCarsService) *carsDependencies {
	return &carsDependencies{
		carsService: carsSrv,
	}
}

func TestCarsRegister(t *testing.T) {
	car := dtos.Car{
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		car dtos.Car
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 201 when body is appropriate",
			args: args{
				car: car,
			},
			wants: wants{
				statusCode: http.StatusCreated,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Register(gomock.Any(), car.ToDomain()).Return(domain.Car{}, nil)
			},
		},
		{
			name: "returns 400 status code when car type in body is not allowed",
			args: args{
				car: dtos.Car{
					Type: "Seda",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 500 status code when car service fails to register car",
			args: args{
				car: car,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Register(gomock.Any(), car.ToDomain()).Return(domain.Car{}, errors.New("error registering car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			body, _ := json.Marshal(test.args.car)
			URL := baseURL + "cars/"
			req, err := http.NewRequest(http.MethodGet, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.Register(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsGet(t *testing.T) {
	car := dtos.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		requestID string
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 200 when id is properly set and car was found",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(car.ToDomain(), nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-an-not-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 404 status code when the car was not found",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(domain.Car{}, errors.New(services.ErrCarNotFound))
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Get(gomock.Any(), car.ID).Return(domain.Car{}, errors.New("error getting car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "cars/" + test.args.requestID)
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.Get(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsFullUpdate(t *testing.T) {
	car := dtos.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		requestID string
		car       dtos.Car
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 200 when car was successfully updated",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-an-not-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 400 status code when car status in body is not allowed",
			args: args{
				requestID: car.ID.String(),
				car: dtos.Car{
					Status: "available",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 404 status code when car was not found by id",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(errors.New(services.ErrCarNotFound))
			},
		},
		{
			name: "returns 500 status code when car service fails to update the car",
			args: args{
				requestID: car.ID.String(),
				car:       car,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().FullUpdate(gomock.Any(), car.ToDomain()).Return(errors.New("error registering car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "cars/" + test.args.requestID)
			URL := urlObj.String()

			body, _ := json.Marshal(test.args.car)
			req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.FullUpdate(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsDelete(t *testing.T) {
	car := dtos.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		City:           "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		requestID string
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 204 when car register was deleted successfully",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNoContent,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Delete(gomock.Any(), car.ID).Return(nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-an-not-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: car.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().Delete(gomock.Any(), car.ID).Return(errors.New("error getting car"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "cars/" + test.args.requestID)
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodDelete, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.requestID,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.Delete(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestCarsList(t *testing.T) {
	foundCars := []domain.Car{
		{
			ID:             uuid.New(),
			Type:           "Sedan",
			Seats:          5,
			HourlyRentCost: 90,
			City:           "New York",
			Status:         "Available",
		},
		{
			ID:             uuid.New(),
			Type:           "Sedan",
			Seats:          5,
			HourlyRentCost: 100,
			City:           "New York",
			Status:         "Available",
		},
	}

	type args struct {
		city        string
		from_car_id string
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies)
	}{
		{
			name: "returns status code 200 when city and from_car_id are provided and service works with no error",
			args: args{
				city:        "New York",
				from_car_id: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().List(gomock.Any(), "New York", "5ae5d956-5a8d-40dd-9aef-5340fda345e8").Return(foundCars, nil)
			},
		},
		{
			name: "returns status code 400 when city query param is empty",
			args: args{
				city:        "",
				from_car_id: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *carsDependencies) {
			},
		},
		{
			name: "returns status code 200 when city param arrives and service works with no error",
			args: args{
				city:        "New York",
				from_car_id: "",
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().List(gomock.Any(), "New York", "").Return(foundCars, nil)
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				city:        "New York",
				from_car_id: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *carsDependencies) {
				d.carsService.EXPECT().List(gomock.Any(), "New York", "5ae5d956-5a8d-40dd-9aef-5340fda345e8").Return([]domain.Car{}, errors.New("error getting cars list"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockCarsService(mockCtlr)
			d := NewCarsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			values := url.Values{}
			values.Set("city", test.args.city)
			values.Set("from_car_id", test.args.from_car_id)
			urlObj, _ := url.Parse(baseURL + "cars/?" + values.Encode())
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"city":        test.args.city,
				"from_car_id": test.args.from_car_id,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewCars(carsSrv)
			carsHandler.List(rr, req)

			if rr.Result().StatusCode == http.StatusOK {
				body := dtos.ListCarsResponse{}
				err = json.Unmarshal(rr.Body.Bytes(), &body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body.Cars) != 0 {
					assert.Equal(t, foundCars[0].ID, body.Cars[0].ID)
					assert.Equal(t, len(foundCars), len(body.Cars))
				}
			}
			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
