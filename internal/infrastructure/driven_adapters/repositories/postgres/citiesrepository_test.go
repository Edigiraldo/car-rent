package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/services"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type citiesDependencies struct {
	db *mocks.MockDatabase
}

func NewCitiesDependencies(db *mocks.MockDatabase) *citiesDependencies {
	return &citiesDependencies{
		db: db,
	}
}

func TestCitiesGetIdByName(t *testing.T) {
	initConstantsFromRepository(t)

	city_id := uuid.New()

	type args struct {
		ctx  context.Context
		name string
	}
	type wants struct {
		id  uuid.UUID
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*citiesDependencies) *sql.DB
	}{
		{
			name: "returns error when the city name was not found",
			args: args{
				ctx:  context.TODO(),
				name: "LosAngeles",
			},
			wants: wants{
				id:  uuid.Nil,
				err: errors.New(services.ErrInvalidCityName),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT id FROM cities WHERE name = \$1`).
					WillReturnError(sql.ErrNoRows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query rows fails",
			args: args{
				ctx:  context.TODO(),
				name: "Los Angeles",
			},
			wants: wants{
				id:  uuid.Nil,
				err: errors.New("query row error"),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT id FROM cities WHERE name = \$1`).
					WillReturnError(errors.New("query row error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns city id when was successfully found",
			args: args{
				ctx:  context.TODO(),
				name: "Los Angeles",
			},
			wants: wants{
				id:  city_id,
				err: nil,
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				cityIdByte, err := city_id.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(cityIdByte)
				mock.ExpectQuery(`SELECT id FROM cities WHERE name = \$1`).
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
			d := NewCitiesDependencies(db)
			dbHandle := test.setMocks(d)

			citiesRepo := NewCitiesRepository(db)
			ID, err := citiesRepo.GetIdByName(test.args.ctx, test.args.name)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.id, ID)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCitiesGetNameByID(t *testing.T) {
	initConstantsFromRepository(t)

	city_name := "Los Angeles"

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type wants struct {
		name string
		err  error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*citiesDependencies) *sql.DB
	}{
		{
			name: "returns error when the city name was not found",
			args: args{
				ctx: context.TODO(),
				id:  uuid.New(),
			},
			wants: wants{
				name: "",
				err:  errors.New(services.ErrCityNotFound),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT name FROM cities WHERE id = \$1`).
					WillReturnError(sql.ErrNoRows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query rows fails",
			args: args{
				ctx: context.TODO(),
				id:  uuid.New(),
			},
			wants: wants{
				name: "",
				err:  errors.New("query row error"),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT name FROM cities WHERE id = \$1`).
					WillReturnError(errors.New("query row error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns city name when was successfully found",
			args: args{
				ctx: context.TODO(),
				id:  uuid.New(),
			},
			wants: wants{
				name: city_name,
				err:  nil,
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(city_name)
				mock.ExpectQuery(`SELECT name FROM cities WHERE id = \$1`).
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
			d := NewCitiesDependencies(db)
			dbHandle := test.setMocks(d)

			citiesRepo := NewCitiesRepository(db)
			name, err := citiesRepo.GetNameByID(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.name, name)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestListNames(t *testing.T) {
	initConstantsFromRepository(t)

	ns := []string{
		"Chicago",
		"New York",
		"Los Angeles",
	}

	type args struct {
		ctx context.Context
	}
	type wants struct {
		names []string
		err   error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*citiesDependencies) *sql.DB
	}{
		{
			name: "returns error when query context fails",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				names: nil,
				err:   errors.New("query context error"),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery("SELECT name FROM cities ORDER BY name LIMIT 100").
					WillReturnError(errors.New("query context error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows are not as expected",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				names: nil,
				err:   errors.New("sql: expected 2 destination arguments in Scan, not 1"),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"noname", "noname"}).
					AddRow("some name", "some name")
				mock.ExpectQuery("SELECT name FROM cities ORDER BY name LIMIT 100").
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows.Err fails",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				names: nil,
				err:   errors.New("rows.Err error"),
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("city name").
					RowError(0, errors.New("rows.Err error"))
				mock.ExpectQuery("SELECT name FROM cities ORDER BY name LIMIT 100").
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns cities name when there were no errors",
			args: args{
				ctx: context.TODO(),
			},
			wants: wants{
				names: ns,
				err:   nil,
			},
			setMocks: func(d *citiesDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(ns[0]).
					AddRow(ns[1]).
					AddRow(ns[2])
				mock.ExpectQuery("SELECT name FROM cities ORDER BY name LIMIT 100").
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
			d := NewCitiesDependencies(db)
			dbHandle := test.setMocks(d)

			citiesRepo := NewCitiesRepository(db)
			names, err := citiesRepo.ListNames(test.args.ctx)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.names, names)
			assert.Equal(t, test.wants.err, err)
		})
	}
}
