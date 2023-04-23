package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type reservationsDependencies struct {
	reservationsRepository *mocks.MockReservationsRepo
}

func NewReservationsDependencies(reservationsRepo *mocks.MockReservationsRepo) *reservationsDependencies {
	return &reservationsDependencies{
		reservationsRepository: reservationsRepo,
	}
}

func TestReservationsRegister(t *testing.T) {
	type args struct {
		ctx         context.Context
		reservation domain.Reservation
	}
	type wants struct {
		withError bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservation struct is populated appropriately",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				withError: false,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when reservation repository fails to booke the reservation",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				withError: true,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("error booking reservation"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewReservationsDependencies(reservationsRepo)
			test.setMocks(d)

			reservationsService := NewReservations(reservationsRepo)
			_, err := reservationsService.Book(test.args.ctx, test.args.reservation)

			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}
