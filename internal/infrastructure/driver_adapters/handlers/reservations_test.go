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
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, errors.New("car not found"))
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
				d.reservationsService.EXPECT().Book(gomock.Any(), gomock.Any()).Return(domain.Reservation{}, errors.New("user not found"))
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
				d.reservationsService.EXPECT().Get(gomock.Any(), reservation.ID).Return(domain.Reservation{}, errors.New("reservation was not found"))
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
				d.reservationsService.EXPECT().FullUpdate(gomock.Any(), gomock.Any()).Return(errors.New("reservation was not found"))
			},
		},
		{
			name: "returns 400 status code when user was not found",
			args: args{
				requestID:   reservation.ID.String(),
				reservation: reservation,
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().FullUpdate(gomock.Any(), gomock.Any()).Return(errors.New("user not found"))
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
			name: "returns 404 status code when reservation was not found",
			args: args{
				requestID: reservation.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().Delete(gomock.Any(), reservation.ID).Return(errors.New("reservation was not found"))
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

func TestReservationsList(t *testing.T) {
	initConstantsFromHandlers(t)
	datetimeLayout := "2006-01-02T15:04:05Z07:00"
	foundReservations := []domain.Reservation{
		{
			UserID:        uuid.New(),
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
		{
			UserID:        uuid.New(),
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}

	type args struct {
		startDate         string
		endDate           string
		fromReservationId string
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
			name: "returns status code 400 when from_reservation_id query param is not a valid uuid",
			args: args{
				startDate:         time.Now().Format(datetimeLayout),
				endDate:           time.Now().AddDate(0, 0, 8).Format(datetimeLayout),
				fromReservationId: "ttt5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns status code 400 when start_date query param is not in the expected time format",
			args: args{
				startDate:         time.Now().Format("2006-01-02"),
				endDate:           time.Now().AddDate(0, 0, 8).Format(datetimeLayout),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns status code 400 when end_date query param is not in the expected time format",
			args: args{
				startDate:         time.Now().Format(datetimeLayout),
				endDate:           time.Now().AddDate(0, 0, 8).Format("2006-01-02"),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns status code 400 when end_date is not after start_date",
			args: args{
				startDate:         time.Now().Format(datetimeLayout),
				endDate:           time.Now().AddDate(0, 0, -1).Format(datetimeLayout),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns status code 500 when there was an error while searching for data",
			args: args{
				startDate:         time.Now().Format(datetimeLayout),
				endDate:           time.Now().AddDate(0, 0, 8).Format(datetimeLayout),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().List(gomock.Any(), "5ae5d956-5a8d-40dd-9aef-5340fda345e8", gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
			},
		},
		{
			name: "returns status code 200 when the data was successfully found",
			args: args{
				startDate:         time.Now().Format(datetimeLayout),
				endDate:           time.Now().AddDate(0, 0, 8).Format(datetimeLayout),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().List(gomock.Any(), "5ae5d956-5a8d-40dd-9aef-5340fda345e8", gomock.Any(), gomock.Any()).
					Return(foundReservations, nil)
			},
		},
		{
			name: "returns status code 200 when params are empty and the data was successfully found",
			args: args{
				startDate:         "",
				endDate:           "",
				fromReservationId: "",
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().List(gomock.Any(), "", gomock.Any(), gomock.Any()).
					Return(foundReservations, nil)
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
			values := url.Values{}
			values.Set("start_date", test.args.startDate)
			values.Set("end_date", test.args.endDate)
			values.Set("from_reservation_id", test.args.fromReservationId)
			urlObj, _ := url.Parse(baseURL + "reservations/?" + values.Encode())
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"start_date":          test.args.startDate,
				"end_date":            test.args.endDate,
				"from_reservation_id": test.args.fromReservationId,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			reservationsHandler := NewReservations(reservationsSrv)
			reservationsHandler.List(rr, req)

			if rr.Result().StatusCode == http.StatusOK {
				body := dtos.Reservations{}
				err = json.Unmarshal(rr.Body.Bytes(), &body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body.Reservations) != 0 {
					assert.Equal(t, foundReservations[0].ID, body.Reservations[0].ID)
					assert.Equal(t, len(foundReservations), len(body.Reservations))
				}
			}
			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestReservationsListByCarID(t *testing.T) {
	car_id := "1c6bd954-7e8d-73df-8ae9-6905fda236e8"
	c_id, _ := uuid.Parse(car_id)
	foundReservations := []domain.Reservation{
		{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			CarID:         c_id,
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
		{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			CarID:         c_id,
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}

	type args struct {
		car_id string
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
			name: "returns status code 200 when car reservations where found",
			args: args{
				car_id: car_id,
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().GetByCarID(gomock.Any(), c_id).Return(foundReservations, nil)
			},
		},
		{
			name: "returns status code 400 when car id param is empty",
			args: args{
				car_id: "",
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
				car_id: car_id,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().GetByCarID(gomock.Any(), c_id).Return([]domain.Reservation{}, errors.New("error getting reservations list"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			carsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(carsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/cars/"
			values := url.Values{}
			values.Set("id", test.args.car_id)
			urlObj, _ := url.Parse(baseURL + values.Encode() + "/reservations")
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.car_id,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			carsHandler := NewReservations(carsSrv)
			carsHandler.GetByCarID(rr, req)

			if rr.Result().StatusCode == http.StatusOK {
				body := dtos.Reservations{}
				err = json.Unmarshal(rr.Body.Bytes(), &body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body.Reservations) != 0 {
					assert.Equal(t, foundReservations[0].ID, body.Reservations[0].ID)
					assert.Equal(t, len(foundReservations), len(body.Reservations))
				}
			}
			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}

func TestReservationsListByUserID(t *testing.T) {
	user_id := "6ef5d956-8a8d-22dd-0aef-5340fda236e8"
	u_id, _ := uuid.Parse(user_id)
	foundReservations := []domain.Reservation{
		{
			ID:            uuid.New(),
			UserID:        u_id,
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
		{
			ID:            uuid.New(),
			UserID:        u_id,
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}

	type args struct {
		user_id string
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
			name: "returns status code 200 when user reservaations where found",
			args: args{
				user_id: user_id,
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().GetByUserID(gomock.Any(), u_id).Return(foundReservations, nil)
			},
		},
		{
			name: "returns status code 400 when user id param is empty",
			args: args{
				user_id: "",
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
				user_id: user_id,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsService.EXPECT().GetByUserID(gomock.Any(), u_id).Return([]domain.Reservation{}, errors.New("error getting reservations list"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsSrv := mocks.NewMockReservationsService(mockCtlr)
			d := NewReservationsDependencies(reservationsSrv)
			test.setMocks(d)

			baseURL := "/api/v1/users/"
			values := url.Values{}
			values.Set("id", test.args.user_id)
			urlObj, _ := url.Parse(baseURL + values.Encode() + "/reservations")
			URL := urlObj.String()

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Include request vars for gorilla mux to interpret path params
			vars := map[string]string{
				"id": test.args.user_id,
			}
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()

			reservationHandler := NewReservations(reservationsSrv)
			reservationHandler.GetByUserID(rr, req)

			if rr.Result().StatusCode == http.StatusOK {
				body := dtos.Reservations{}
				err = json.Unmarshal(rr.Body.Bytes(), &body)
				if err != nil {
					t.Fatal(err)
				}
				if len(body.Reservations) != 0 {
					assert.Equal(t, foundReservations[0].ID, body.Reservations[0].ID)
					assert.Equal(t, len(foundReservations), len(body.Reservations))
				}
			}
			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
