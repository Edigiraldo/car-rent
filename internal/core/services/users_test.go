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

type usersDependencies struct {
	usersRepository        *mocks.MockUsersRepo
	reservationsRepository *mocks.MockReservationsRepo
}

func NewUsersDependencies(usersRepo *mocks.MockUsersRepo, reservationsRepo *mocks.MockReservationsRepo) *usersDependencies {
	return &usersDependencies{
		usersRepository:        usersRepo,
		reservationsRepository: reservationsRepo,
	}
}

func TestUsersRegister(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.User
	}
	type wants struct {
		withError bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns nil error when user struct is populated appropriately",
			args: args{
				ctx: context.TODO(),
				user: domain.User{
					ID:        uuid.UUID([]byte("754f2869-9f1e-4691-bbe1-2e75412e2912")),
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Customer",
					Status:    "Active",
				},
			},
			wants: wants{
				withError: false,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when user type is not an allowed value",
			args: args{
				ctx: context.TODO(),
				user: domain.User{
					ID:        uuid.UUID([]byte("754f2869-9f1e-4691-bbe1-2e75412e2912")),
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Not an allowed type",
					Status:    "Active",
				},
			},
			wants: wants{
				withError: true,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("type is not allowed"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersRepo := mocks.NewMockUsersRepo(mockCtlr)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewUsersDependencies(usersRepo, reservationsRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo, reservationsRepo)
			_, err := usersService.Register(test.args.ctx, test.args.user)

			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}

func TestUsersGet(t *testing.T) {
	user := domain.User{
		FirstName: "Richard",
		LastName:  "Feynman",
		Email:     "richard.feynman@caltech.edu.us",
		Type:      "Customer",
		Status:    "Active",
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	type wants struct {
		user      domain.User
		withError bool
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns nil error when user was found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  user.ID,
			},
			wants: wants{
				user:      user,
				withError: false,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Get(gomock.Any(), user.ID).Return(user, nil)
			},
		},
		{
			name: "returns an error when user was not found by the given id",
			args: args{
				ctx: context.TODO(),
				ID:  user.ID,
			},
			wants: wants{
				user:      domain.User{},
				withError: true,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Get(gomock.Any(), user.ID).Return(domain.User{}, errors.New(ErrUserNotFound))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersRepo := mocks.NewMockUsersRepo(mockCtlr)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewUsersDependencies(usersRepo, reservationsRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo, reservationsRepo)
			user, err := usersService.Get(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.user, user)
			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}

func TestUsersFullUpdate(t *testing.T) {
	user := domain.User{
		ID:        uuid.New(),
		FirstName: "Richard",
		LastName:  "Feynman",
		Email:     "richard.feynman@caltech.edu.us",
		Type:      "Customer",
		Status:    "Active",
	}

	type args struct {
		ctx  context.Context
		user domain.User
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns nil error when user update was successfully updated",
			args: args{
				ctx:  context.TODO(),
				user: user,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().FullUpdate(gomock.Any(), user).Return(nil)
			},
		},
		{
			name: "returns an error when user update fails",
			args: args{
				ctx:  context.TODO(),
				user: user,
			},
			wants: wants{
				err: errors.New("failure while updating user"),
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().FullUpdate(gomock.Any(), user).Return(errors.New("failure while updating user"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersRepo := mocks.NewMockUsersRepo(mockCtlr)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewUsersDependencies(usersRepo, reservationsRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo, reservationsRepo)
			err := usersService.FullUpdate(test.args.ctx, test.args.user)

			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestUsersDelete(t *testing.T) {
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
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns nil error when user deletion was successful",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "returns an error when user deletion fails",
			args: args{
				ctx: context.TODO(),
				ID:  uuid.New(),
			},
			wants: wants{
				err: errors.New("failure while updating user"),
			},
			setMocks: func(d *usersDependencies) {
				d.usersRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("failure while updating user"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersRepo := mocks.NewMockUsersRepo(mockCtlr)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewUsersDependencies(usersRepo, reservationsRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo, reservationsRepo)
			err := usersService.Delete(test.args.ctx, test.args.ID)

			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestGetReservations(t *testing.T) {
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
		setMocks func(*usersDependencies)
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
			setMocks: func(d *usersDependencies) {
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
			setMocks: func(d *usersDependencies) {
				d.reservationsRepository.EXPECT().GetByUserID(gomock.Any(), u_id).Return([]domain.Reservation{}, errors.New("there was some internal error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersRepo := mocks.NewMockUsersRepo(mockCtlr)
			reservationsRepo := mocks.NewMockReservationsRepo(mockCtlr)
			d := NewUsersDependencies(usersRepo, reservationsRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo, reservationsRepo)
			reservations, err := usersService.GetReservations(test.args.ctx, test.args.userID)

			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}
