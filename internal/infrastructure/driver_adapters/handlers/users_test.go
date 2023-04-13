package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driver_adapters/dtos"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type usersDependencies struct {
	usersService *mocks.MockUsersService
}

func NewUsersDependencies(usersSrv *mocks.MockUsersService) *usersDependencies {
	return &usersDependencies{
		usersService: usersSrv,
	}
}

func TestUsersSingUp(t *testing.T) {
	initConstantsFromHandlers(t)

	user := dtos.User{
		FirstName: "Richard",
		LastName:  "Feynman",
		Email:     "richard.feynman@caltech.edu.us",
		Type:      "Customer",
		Status:    "Active",
	}

	type args struct {
		user dtos.User
	}
	type wants struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns status code 201 when body is appropriate",
			args: args{
				user: user,
			},
			wants: wants{
				statusCode: http.StatusCreated,
			},
			setMocks: func(d *usersDependencies) {
				d.usersService.EXPECT().Register(gomock.Any(), user.ToDomain()).Return(domain.User{}, nil)
			},
		},
		{
			name: "returns 400 status code when user type in body is not allowed",
			args: args{
				user: dtos.User{
					Type: "Not an allowed type",
				},
			},
			wants: wants{
				statusCode: http.StatusBadRequest,
			},
			setMocks: func(d *usersDependencies) {
			},
		},
		{
			name: "returns 500 status code when user service fails to register user",
			args: args{
				user: user,
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *usersDependencies) {
				d.usersService.EXPECT().Register(gomock.Any(), user.ToDomain()).Return(domain.User{}, errors.New("error registering user"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			usersSrv := mocks.NewMockUsersService(mockCtlr)
			d := NewUsersDependencies(usersSrv)
			test.setMocks(d)

			baseURL := "/api/v1/"
			body, _ := json.Marshal(test.args.user)
			URL := baseURL + "users/"
			req, err := http.NewRequest(http.MethodGet, URL, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			usersHandler := NewUsers(usersSrv)
			usersHandler.SignUp(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
