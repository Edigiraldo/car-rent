package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
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
					StartDate:     time.Now().Add(1 * time.Hour),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				withError: false,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "returns an error when reservation repository fails to book the reservation",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now().Add(1 * time.Hour),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				withError: true,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("error booking reservation"))
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "returns an error when validations fail",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     time.Now().Add(1 * time.Hour),
					EndDate:       time.Now().AddDate(0, 0, 7),
				},
			},
			wants: wants{
				withError: true,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some validation failed"))
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

func TestReservationsGet(t *testing.T) {
	reservation := domain.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now().Add(1 * time.Hour),
		EndDate:       time.Now().AddDate(0, 0, 7),
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	type wants struct {
		reservation domain.Reservation
		withError   bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservation was found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  reservation.ID,
			},
			wants: wants{
				reservation: reservation,
				withError:   false,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Get(gomock.Any(), reservation.ID).Return(reservation, nil)
			},
		},
		{
			name: "returns an error when reservation was not found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  reservation.ID,
			},
			wants: wants{
				reservation: domain.Reservation{},
				withError:   true,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Get(gomock.Any(), reservation.ID).Return(domain.Reservation{}, errors.New(ErrReservationNotFound))
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
			reservation, err := reservationsService.Get(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.reservation, reservation)
			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}

func TestReservationsFullUpdate(t *testing.T) {
	reservation := domain.Reservation{
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now().Add(1 * time.Hour),
		EndDate:       time.Now().AddDate(0, 0, 7),
	}

	type args struct {
		ctx         context.Context
		reservation domain.Reservation
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservation update was successfully done",
			args: args{
				ctx:         context.TODO(),
				reservation: reservation,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().FullUpdate(gomock.Any(), reservation).Return(nil)
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "returns an error when reservation update fails",
			args: args{
				ctx:         context.TODO(),
				reservation: reservation,
			},
			wants: wants{
				err: errors.New("failure while updating reservation"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().FullUpdate(gomock.Any(), reservation).Return(errors.New("failure while updating reservation"))
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "returns an error when validations fail",
			args: args{
				ctx:         context.TODO(),
				reservation: reservation,
			},
			wants: wants{
				err: errors.New("some validation failed"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some validation failed"))
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
			err := reservationsService.FullUpdate(test.args.ctx, test.args.reservation)

			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsDelete(t *testing.T) {
	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservation deletion was successful",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when reservation deletion fails",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: errors.New("failure while deleting reservation"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("failure while deleting reservation"))
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
			err := reservationsService.Delete(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsList(t *testing.T) {
	initConstantsFromServices(t)
	var nilTime time.Time

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
		ctx               context.Context
		startDate         time.Time
		endDate           time.Time
		fromReservationId string
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when params are nil and repository did not fail",
			args: args{
				ctx:               context.TODO(),
				startDate:         nilTime,
				endDate:           nilTime,
				fromReservationId: "",
			},
			wants: wants{
				reservations: foundReservations,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().List(gomock.Any(), "00000000-0000-0000-0000-000000000000", gomock.Any(), gomock.Any(), gomock.Any()).Return(foundReservations, nil)
			},
		},
		{
			name: "returns nil error when repository did not fail",
			args: args{
				ctx:               context.TODO(),
				startDate:         time.Now(),
				endDate:           time.Now().AddDate(0, 0, 7),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				reservations: foundReservations,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().List(gomock.Any(), "5ae5d956-5a8d-40dd-9aef-5340fda345e8", gomock.Any(), gomock.Any(), gomock.Any()).Return(foundReservations, nil)
			},
		},
		{
			name: "returns error when repository fails",
			args: args{
				ctx:               context.TODO(),
				startDate:         time.Now(),
				endDate:           time.Now().AddDate(0, 0, 7),
				fromReservationId: "5ae5d956-5a8d-40dd-9aef-5340fda345e8",
			},
			wants: wants{
				reservations: []domain.Reservation{},
				err:          errors.New("internal server error"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().List(gomock.Any(), "5ae5d956-5a8d-40dd-9aef-5340fda345e8", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal server error"))
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
			reservations, err := reservationsService.List(test.args.ctx, test.args.fromReservationId, test.args.startDate, test.args.endDate)

			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestGetByCarID(t *testing.T) {
	initConstantsFromServices(t)

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
		ctx   context.Context
		CarID uuid.UUID
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservations were found",
			args: args{
				ctx:   context.TODO(),
				CarID: c_id,
			},
			wants: wants{
				reservations: foundReservations,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarID(gomock.Any(), c_id).Return(foundReservations, nil)
			},
		},
		{
			name: "returns an error when repository fails retrieving reservations",
			args: args{
				ctx:   context.TODO(),
				CarID: c_id,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("there was some internal error"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarID(gomock.Any(), c_id).Return([]domain.Reservation{}, errors.New("there was some internal error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewReservationsDependencies(reservationsRepo)
			test.setMocks(d)

			carsService := NewReservations(reservationsRepo)
			reservations, err := carsService.GetByCarID(test.args.ctx, test.args.CarID)

			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestGetByUserID(t *testing.T) {
	initConstantsFromServices(t)

	user_id := "7d4bd954-8a8d-55dd-0aef-6440fda236e8"
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
		ctx    context.Context
		userID uuid.UUID
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns nil error when reservations were found",
			args: args{
				ctx:    context.TODO(),
				userID: u_id,
			},
			wants: wants{
				reservations: foundReservations,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByUserID(gomock.Any(), u_id).Return(foundReservations, nil)
			},
		},
		{
			name: "returns an error when repository fails retrieving reservations",
			args: args{
				ctx:    context.TODO(),
				userID: u_id,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("there was some internal error"),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByUserID(gomock.Any(), u_id).Return([]domain.Reservation{}, errors.New("there was some internal error"))
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
			reservations, err := reservationsService.GetByUserID(test.args.ctx, test.args.userID)

			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCheckReservation(t *testing.T) {
	now := time.Now()
	initConstantsFromServices(t)

	type args struct {
		ctx         context.Context
		reservation domain.Reservation
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies)
	}{
		{
			name: "returns an error when time frame is invalid",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(1 * time.Hour),
					EndDate:       now.Add(1 * time.Hour),
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidReservationTimeFrame),
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns an error when period is shorter than allowed",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(1 * time.Hour),
					EndDate:       now.Add(1*time.Hour + 1*time.Second),
				},
			},
			wants: wants{
				err: fmt.Errorf("%s (%d hours)", ErrMinimumReservationHours, constants.Values.MINIMUM_RESERVATION_HOURS),
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns an error when car is unavailable",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					ID:            uuid.New(),
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(1 * time.Hour),
					EndDate:       now.Add(30 * 24 * time.Hour),
				},
			},
			wants: wants{
				err: errors.New(ErrCarNotAvailable),
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Reservation{
					{
						ID:            uuid.New(),
						UserID:        uuid.New(),
						CarID:         uuid.New(),
						Status:        "Reserved",
						PaymentStatus: "Pending",
						StartDate:     now,
						EndDate:       now.Add(30 * 24 * time.Hour),
					},
				}, nil)
			},
		},
		{
			name: "returns nil error when new reservation time frame intersects the past one",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(1 * time.Hour),
					EndDate:       now.Add(30 * 24 * time.Hour),
				},
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Reservation{
					{
						UserID:        uuid.New(),
						CarID:         uuid.New(),
						Status:        "Reserved",
						PaymentStatus: "Pending",
						StartDate:     now,
						EndDate:       now.Add(30 * 24 * time.Hour),
					},
				}, nil)
			},
		},
		{
			name: "returns an error when car reservation start date is before now",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(-1 * time.Minute),
					EndDate:       now,
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidReservationTimeFrame),
			},
			setMocks: func(d *reservationsDependencies) {
			},
		},
		{
			name: "returns nil error when all validations pass",
			args: args{
				ctx: context.TODO(),
				reservation: domain.Reservation{
					UserID:        uuid.New(),
					CarID:         uuid.New(),
					Status:        "Reserved",
					PaymentStatus: "Pending",
					StartDate:     now.Add(10 * 24 * time.Hour),
					EndDate:       now.Add(30 * 24 * time.Hour),
				},
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) {
				d.reservationsRepository.EXPECT().GetByCarIDAndTimeFrame(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
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
			err := reservationsService.CheckReservation(test.args.ctx, test.args.reservation)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
