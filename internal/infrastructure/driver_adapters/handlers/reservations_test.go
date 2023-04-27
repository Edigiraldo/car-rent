package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type reservationsDependencies struct {
	reservationsService *mocks.MockReservationsService
}

func NewReservationsDependencies(reservationsSrv *mocks.MockReservationsService) *reservationsDependencies {
	return &reservationsDependencies{
		reservationsService: reservationsSrv,
	}
}

func TestReservationsBook(t *testing.T) {
	initConstantsFromHandlers(t)

	reservation := dtos.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, 7),
	}

	type args struct {
		reservation dtos.Reservation
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns status code 201 when body is appropriate",
			args: args{
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusCreated,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, nil)
			},
		},
		{
			name: "returns 400 status code when reservation body is empty",
			args: args{
				reservation: dtos.Reservation{},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns 400 status code when reservation car id was not found",
			args: args{
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, errors.New(services.ErrCarNotFound))
			},
		},
		{
			name: "returns 400 status code when reservation user id was not found",
			args: args{
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, errors.New(services.ErrUserNotFound))
			},
		},
		{
			name: "returns 500 status code when reservation service fails to book the reservation",
			args: args{
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, errors.New("error booking reservation"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(reservationsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			body, _ := json.Marshal(test.args.reservation)
			URL := baseURL + "reservations/"
			req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			reservationsHandler := NewReservations(reservationsSrv)
			reservationsHandler.Book(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestReservationsGet(t *testing.T) {
	reservation := dtos.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, 7),
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
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns status code 200 when id is properly set and reservation was found",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Get(gomock.Any(), reservation.ID).Return(reservation.ToDomain(), nil)
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
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns 404 status code when the reservation was not found",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Get(gomock.Any(), reservation.ID).Return(domain.Reservation{}, errors.New(services.ErrReservationNotFound))
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Get(gomock.Any(), reservation.ID).Return(domain.Reservation{}, errors.New("error getting reservation"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(reservationsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "reservations/" + test.args.requestID)
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

			reservationsHandler := NewReservations(reservationsSrv)
			reservationsHandler.Get(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestReservationsFullUpdate(t *testing.T) {
	initConstantsFromHandlers(t)

	reservation := dtos.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, 7),
	}

	type args struct {
		requestID   string
		reservation dtos.Reservation
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns status code 200 when reservation was successfully updated",
			args: args{
				requestID:   reservation.ID.String(),
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().FullUpdate(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-not-a-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns 400 status code when reservation status in body is not allowed",
			args: args{
				requestID: reservation.ID.String(),
				reservation: dtos.Reservation{
					Status: "available",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns 404 status code when reservation was not found by id",
			args: args{
				requestID:   reservation.ID.String(),
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().FullUpdate(gomock.Any(), gomock.Any()).Return(errors.New(services.ErrReservationNotFound))
			},
		},
		{
			name: "returns 500 status code when reservation service fails to update the reservation",
			args: args{
				requestID:   reservation.ID.String(),
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().FullUpdate(gomock.Any(), gomock.Any()).Return(errors.New("error registering reservation"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(reservationsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "reservations/" + test.args.requestID)
			URL := urlObj.String()

			body, _ := json.Marshal(test.args.reservation)
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

			reservationsHandler := NewReservations(reservationsSrv)
			reservationsHandler.FullUpdate(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestReservationsDelete(t *testing.T) {
	reservation := dtos.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, 7),
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
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns status code 204 when reservation register was deleted successfully",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNoContent,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Delete(gomock.Any(), reservation.ID).Return(nil)
			},
		},
		{
			name: "returns 400 status code when path param id is not an uuid",
			args: args{
				requestID: "this-is-not-a-uuid",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Delete(gomock.Any(), reservation.ID).Return(errors.New("error deleting reservation"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(reservationsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			urlObj, _ := url.Parse(baseURL + "reservations/" + test.args.requestID)
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

			reservationsHandler := NewReservations(reservationsSrv)
			reservationsHandler.Delete(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
