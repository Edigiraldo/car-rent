package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type usersDependencies struct {
	db *mocks.MockDatabase
}

func NewUsersDependencies(db *mocks.MockDatabase) *usersDependencies {
	return &usersDependencies{
		db: db,
	}
}

func TestUsersInsert(t *testing.T) {
	initConstantsFromRepository(t)

	du := domain.User{
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
		setMocks func(*usersDependencies) *sql.DB
	}{
		{
			name: "returns error when exec context fails",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: errors.New("exec context"),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(du.ID, du.FirstName, du.LastName, du.Email, du.Type, du.Status).
					WillReturnError(errors.New("exec context"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when user was successfully inserted",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("INSERT INTO users").
					WithArgs(du.ID, du.FirstName, du.LastName, du.Email, du.Type, du.Status).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewUsersDependencies(db)
			dbHandle := test.setMocks(d)

			usersRepo := NewUsersRepository(db)
			err := usersRepo.Insert(test.args.ctx, test.args.user)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestUsersGet(t *testing.T) {
	initConstantsFromRepository(t)

	du := domain.User{
		ID:        uuid.New(),
		FirstName: "Richard",
		LastName:  "Feynman",
		Email:     "richard.feynman@caltech.edu.us",
		Type:      "Customer",
		Status:    "Active",
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type wants struct {
		user domain.User
		err  error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies) *sql.DB
	}{
		{
			name: "returns error when no rows were found",
			args: args{
				ctx: context.TODO(),
				id:  du.ID,
			},
			wants: wants{
				user: domain.User{},
				err:  errors.New(services.ErrUserNotFound),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM users WHERE ID = \$1`).
					WithArgs(du.ID).
					WillReturnError(sql.ErrNoRows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query to user repo fails",
			args: args{
				ctx: context.TODO(),
				id:  du.ID,
			},
			wants: wants{
				user: domain.User{},
				err:  errors.New("there was some error"),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM users WHERE ID = \$1`).
					WithArgs(du.ID).
					WillReturnError(errors.New("there was some error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns user when was successfully found",
			args: args{
				ctx: context.TODO(),
				id:  du.ID,
			},
			wants: wants{
				user: du,
				err:  nil,
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := du.ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "type", "status"}).
					AddRow(userIdByte, du.FirstName, du.LastName, du.Email, du.Type, du.Status)
				mock.ExpectQuery(`SELECT \* FROM users WHERE ID = \$1`).
					WithArgs(du.ID).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewUsersDependencies(db)
			dbHandle := test.setMocks(d)

			usersRepo := NewUsersRepository(db)
			user, err := usersRepo.Get(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.user, user)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestUsersFullUpdate(t *testing.T) {
	initConstantsFromRepository(t)

	du := domain.User{
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
		setMocks func(*usersDependencies) *sql.DB
	}{
		{
			name: "returns error when exec context fails",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: errors.New("exec context"),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("UPDATE users SET").
					WithArgs(du.FirstName, du.LastName, du.Email, du.Type, du.Status, du.ID).
					WillReturnError(errors.New("exec context"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows affected method fails",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: errors.New("rows affected error"),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewErrorResult(errors.New("rows affected error"))
				mock.ExpectExec("UPDATE users SET").
					WithArgs(du.FirstName, du.LastName, du.Email, du.Type, du.Status, du.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when user was not found",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: errors.New(services.ErrUserNotFound),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE users SET").
					WithArgs(du.FirstName, du.LastName, du.Email, du.Type, du.Status, du.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when user was updated",
			args: args{
				ctx:  context.TODO(),
				user: du,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("UPDATE users SET").
					WithArgs(du.FirstName, du.LastName, du.Email, du.Type, du.Status, du.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewUsersDependencies(db)
			dbHandle := test.setMocks(d)

			usersRepo := NewUsersRepository(db)
			err := usersRepo.FullUpdate(test.args.ctx, test.args.user)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestUsersDelete(t *testing.T) {
	initConstantsFromRepository(t)

	id := uuid.New()

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*usersDependencies) *sql.DB
	}{
		{
			name: "returns error when deletion fails",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: errors.New("id not found"),
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("DELETE FROM users").
					WithArgs(id).
					WillReturnError(errors.New("id not found"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when user was deleted",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *usersDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM users").
					WithArgs(id).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewUsersDependencies(db)
			dbHandle := test.setMocks(d)

			usersRepo := NewUsersRepository(db)
			err := usersRepo.Delete(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}
