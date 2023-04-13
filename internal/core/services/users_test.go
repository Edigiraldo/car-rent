package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type usersDependencies struct {
	usersRepository *mocks.MockUsersRepo
}

func NewUsersDependencies(usersRepo *mocks.MockUsersRepo) *usersDependencies {
	return &usersDependencies{
		usersRepository: usersRepo,
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
			d := NewUsersDependencies(usersRepo)
			test.setMocks(d)

			usersService := NewUsers(usersRepo)
			_, err := usersService.Register(test.args.ctx, test.args.user)

			assert.Equal(t, test.wants.withError, err != nil)
		})
	}
}
