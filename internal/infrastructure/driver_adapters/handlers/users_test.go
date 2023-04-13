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
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Not an allowed type",
					Status:    "Active",
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

func TestUsersGet(t *testing.T) {
	user := dtos.User{
		ID:        uuid.New(),
		FirstName: "Richard",
		LastName:  "Feynman",
		Email:     "richard.feynman@caltech.edu.us",
		Type:      "Customer",
		Status:    "Active",
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
		setMocks func(*usersDependencies)
	}{
		{
			name: "returns status code 200 when id is properly set and user was found",
			args: args{
				requestID: user.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusOK,
			},
			setMocks: func(d *usersDependencies) {
				d.usersService.EXPECT().Get(gomock.Any(), user.ID).Return(user.ToDomain(), nil)
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
			setMocks: func(d *usersDependencies) {
			},
		},
		{
			name: "returns 404 status code when the user was not found",
			args: args{
				requestID: user.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusNotFound,
			},
			setMocks: func(d *usersDependencies) {
				d.usersService.EXPECT().Get(gomock.Any(), user.ID).Return(domain.User{}, errors.New(services.ErrUserNotFound))
			},
		},
		{
			name: "returns 500 status code when there is a server error",
			args: args{
				requestID: user.ID.String(),
			},
			wants: wants{
				statusCode: http.StatusInternalServerError,
			},
			setMocks: func(d *usersDependencies) {
				d.usersService.EXPECT().Get(gomock.Any(), user.ID).Return(domain.User{}, errors.New("error getting user"))
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
			urlObj, _ := url.Parse(baseURL + "users/" + test.args.requestID)
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

			usersHandler := NewUsers(usersSrv)
			usersHandler.Get(rr, req)

			assert.Equal(t, test.wants.statusCode, rr.Code)
		})
	}
}
