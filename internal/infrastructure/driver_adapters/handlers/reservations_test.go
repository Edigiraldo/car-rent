package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
