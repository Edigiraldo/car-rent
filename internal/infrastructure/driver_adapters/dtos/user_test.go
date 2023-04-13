package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFromBody(t *testing.T) {
	initConstantsFromDtos(t)

	type args struct {
		user User
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "returns user structure when body structure is as expected",
			args: args{
				user: User{
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Customer",
					Status:    "Active",
				},
			},
			wants: wants{
				err: nil,
			},
		},
		{
			name: "returns empty first name when received one is null",
			args: args{
				user: User{
					FirstName: "",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Customer",
					Status:    "Active",
				},
			},
			wants: wants{
				err: errors.New(ErrEmptyFirstName),
			},
		},
		{
			name: "returns empty last name when received one is null",
			args: args{
				user: User{
					FirstName: "Richard",
					LastName:  "",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Customer",
					Status:    "Active",
				},
			},
			wants: wants{
				err: errors.New(ErrEmptyLastName),
			},
		},
		{
			name: "returns invalid email when it is not in an appropriate format",
			args: args{
				user: User{
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman-caltech.edu.us",
					Type:      "Customer",
					Status:    "Active",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidEmail),
			},
		},
		{
			name: "returns invalid user type when received one is not one of the expected values",
			args: args{
				user: User{
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Some invalid type",
					Status:    "Active",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidUserType),
			},
		},
		{
			name: "returns invalid user status when received one is not one of the expected values",
			args: args{
				user: User{
					FirstName: "Richard",
					LastName:  "Feynman",
					Email:     "richard.feynman@caltech.edu.us",
					Type:      "Customer",
					Status:    "Not a valid status",
				},
			},
			wants: wants{
				err: errors.New(ErrInvalidUserStatus),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userJSON, err := json.Marshal(test.args.user)
			if err != nil {
				t.Fatal(err)
			}

			userReader := bytes.NewBuffer(userJSON)
			_, err = UserFromBody(userReader)

			assert.Equal(t, test.wants.err, err)
		})
	}
}
