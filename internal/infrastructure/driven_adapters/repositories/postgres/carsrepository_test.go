package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/pkg/constants"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var pathToRoot = "./../../../../.."

type carsDependencies struct {
	db         *mocks.MockDatabase
	citiesRepo *mocks.MockCitiesRepo
}

func NewCarsDependencies(db *mocks.MockDatabase, citiesRepo *mocks.MockCitiesRepo) *carsDependencies {
	return &carsDependencies{
		db:         db,
		citiesRepo: citiesRepo,
	}
}

func initConstantsFromRepository(t *testing.T) {
	if err := constants.InitValuesFrom(pathToRoot); err != nil {
		t.Fatal(err)
	}
}

func TestCarsInsert(t *testing.T) {
	initConstantsFromRepository(t)

	dc := domain.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		CityName:       "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		ctx context.Context
		car domain.Car
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies) *sql.DB
	}{
		{
			name: "returns nil error when car has been inserted succesfully",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("INSERT INTO cars").
					WithArgs(dc.ID, dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when car insertion fails",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("there was some error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("INSERT INTO cars").
					WithArgs(dc.ID, dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status).
					WillReturnError(errors.New("there was some error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query to cities repo fails",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("there was some error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(uuid.Nil, errors.New("there was some error"))

				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewCarsDependencies(db, citiesRepo)
			dbHandle := test.setMocks(d)

			carsRepo := NewCarsRepository(db, citiesRepo)
			err := carsRepo.Insert(test.args.ctx, test.args.car)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCarsGet(t *testing.T) {
	initConstantsFromRepository(t)

	dc := domain.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		CityName:       "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type wants struct {
		car domain.Car
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies) *sql.DB
	}{
		{
			name: "returns car when was successfully found",
			args: args{
				ctx: context.TODO(),
				id:  dc.ID,
			},
			wants: wants{
				car: dc,
				err: nil,
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetNameByID(gomock.Any(), gomock.Any()).Return("Los Angeles", nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dc.ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				cityIdByte, err := uuid.New().MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "type", "seats", "hourly_rent_cost", "city_id", "status"}).
					AddRow(carIdByte, dc.Type, dc.Seats, dc.HourlyRentCost, cityIdByte, dc.Status)
				mock.ExpectQuery(`SELECT \* FROM cars WHERE ID = \$1`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query to cities repo fails",
			args: args{
				ctx: context.TODO(),
				id:  dc.ID,
			},
			wants: wants{
				car: domain.Car{},
				err: errors.New("there was some error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetNameByID(gomock.Any(), gomock.Any()).Return("", errors.New("there was some error"))

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dc.ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				cityIdByte, err := uuid.New().MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "type", "seats", "hourly_rent_cost", "city_id", "status"}).
					AddRow(carIdByte, dc.Type, dc.Seats, dc.HourlyRentCost, cityIdByte, dc.Status)
				mock.ExpectQuery(`SELECT \* FROM cars WHERE ID = \$1`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query to car repo fails",
			args: args{
				ctx: context.TODO(),
				id:  dc.ID,
			},
			wants: wants{
				car: domain.Car{},
				err: errors.New("there was some error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM cars WHERE ID = \$1`).
					WillReturnError(errors.New("there was some error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when car was not found",
			args: args{
				ctx: context.TODO(),
				id:  dc.ID,
			},
			wants: wants{
				car: domain.Car{},
				err: errors.New("car not found"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM cars WHERE ID = \$1`).
					WillReturnError(sql.ErrNoRows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewCarsDependencies(db, citiesRepo)
			dbHandle := test.setMocks(d)

			carsRepo := NewCarsRepository(db, citiesRepo)
			dc, err := carsRepo.Get(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.car, dc)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCarsFullUpdate(t *testing.T) {
	initConstantsFromRepository(t)

	dc := domain.Car{
		ID:             uuid.New(),
		Type:           "Sedan",
		Seats:          4,
		HourlyRentCost: 21.1,
		CityName:       "Los Angeles",
		Status:         "Available",
	}

	type args struct {
		ctx context.Context
		car domain.Car
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies) *sql.DB
	}{
		{
			name: "returns error when cities repo fails",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("cities repo error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(uuid.Nil, errors.New("cities repo error"))

				return nil
			},
		},
		{
			name: "returns error when exec method fails",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("exec error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("UPDATE cars SET").
					WithArgs(dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status, dc.ID).
					WillReturnError(errors.New("exec error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows updated fails",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("rows affected error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewErrorResult(errors.New("rows affected error"))
				mock.ExpectExec("UPDATE cars SET").
					WithArgs(dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status, dc.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when car was not found",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: errors.New("car not found"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE cars SET").
					WithArgs(dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status, dc.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when car has been updated succesfully",
			args: args{
				ctx: context.TODO(),
				car: dc,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				city_id := uuid.New()
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), dc.CityName).Return(city_id, nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("UPDATE cars SET").
					WithArgs(dc.Type, dc.Seats, dc.HourlyRentCost, city_id, dc.Status, dc.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewCarsDependencies(db, citiesRepo)
			dbHandle := test.setMocks(d)

			carsRepo := NewCarsRepository(db, citiesRepo)
			err := carsRepo.FullUpdate(test.args.ctx, test.args.car)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCarsDelete(t *testing.T) {
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
		setMocks func(*carsDependencies) *sql.DB
	}{
		{
			name: "returns error when deletion fails",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: errors.New("execContext error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("DELETE FROM cars").
					WithArgs(id).
					WillReturnError(errors.New("execContext error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows affected fails",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: errors.New("rows affected error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewErrorResult(errors.New("rows affected error"))
				mock.ExpectExec("DELETE FROM cars").
					WithArgs(id).WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when car was not found",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: errors.New("car not found"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("DELETE FROM cars").
					WithArgs(id).WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when car was deleted",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM cars").
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
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewCarsDependencies(db, citiesRepo)
			dbHandle := test.setMocks(d)

			carsRepo := NewCarsRepository(db, citiesRepo)
			err := carsRepo.Delete(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestCarsList(t *testing.T) {
	initConstantsFromRepository(t)

	dcs := []domain.Car{
		{
			ID:             uuid.New(),
			Type:           "Sedan",
			Seats:          4,
			HourlyRentCost: 21.1,
			CityName:       "Los Angeles",
			Status:         "Available",
		},
	}

	type args struct {
		ctx         context.Context
		cityName    string
		from_car_id string
		limit       uint16
	}
	type wants struct {
		cars []domain.Car
		err  error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*carsDependencies) *sql.DB
	}{
		{
			name: "returns error when query to cities repo fails",
			args: args{
				ctx:         context.TODO(),
				cityName:    "LosAngeles",
				from_car_id: "",
				limit:       20,
			},
			wants: wants{
				cars: nil,
				err:  errors.New("cities repo error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), "LosAngeles").Return(uuid.Nil, errors.New("cities repo error"))

				return nil
			},
		},
		{
			name: "returns error when query context fails",
			args: args{
				ctx:         context.TODO(),
				cityName:    "Los Angeles",
				from_car_id: "",
				limit:       20,
			},
			wants: wants{
				cars: nil,
				err:  errors.New("query context error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), "Los Angeles").Return(uuid.New(), nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`^SELECT \* FROM cars WHERE city_id=\$1 AND id > \$2 ORDER BY id ASC LIMIT \$3$`).
					WillReturnError(errors.New("query context error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows are not as expected",
			args: args{
				ctx:         context.TODO(),
				cityName:    "Los Angeles",
				from_car_id: "",
				limit:       20,
			},
			wants: wants{
				cars: nil,
				err:  errors.New("sql: expected 1 destination arguments in Scan, not 6"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), "Los Angeles").Return(uuid.New(), nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dcs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(carIdByte)
				mock.ExpectQuery(`^SELECT \* FROM cars WHERE city_id=\$1 AND id > \$2 ORDER BY id ASC LIMIT \$3$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows.Err fails",
			args: args{
				ctx:         context.TODO(),
				cityName:    "Los Angeles",
				from_car_id: "",
				limit:       20,
			},
			wants: wants{
				cars: nil,
				err:  errors.New("rows.Err error"),
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), "Los Angeles").Return(uuid.New(), nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dcs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				cityIdByte, err := uuid.New().MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "type", "seats", "hourly_rent_cost", "city_id", "status"}).
					AddRow(carIdByte, dcs[0].Type, dcs[0].Seats, dcs[0].HourlyRentCost, cityIdByte, dcs[0].Status).
					RowError(0, errors.New("rows.Err error"))
				mock.ExpectQuery(`^SELECT \* FROM cars WHERE city_id=\$1 AND id > \$2 ORDER BY id ASC LIMIT \$3$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns cars when were successfully found",
			args: args{
				ctx:         context.TODO(),
				cityName:    "Los Angeles",
				from_car_id: "",
				limit:       20,
			},
			wants: wants{
				cars: dcs,
				err:  nil,
			},
			setMocks: func(d *carsDependencies) *sql.DB {
				d.citiesRepo.EXPECT().GetIdByName(gomock.Any(), "Los Angeles").Return(uuid.New(), nil)

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dcs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				cityIdByte, err := uuid.New().MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "type", "seats", "hourly_rent_cost", "city_id", "status"}).
					AddRow(carIdByte, dcs[0].Type, dcs[0].Seats, dcs[0].HourlyRentCost, cityIdByte, dcs[0].Status)
				mock.ExpectQuery(`^SELECT \* FROM cars WHERE city_id=\$1 AND id > \$2 ORDER BY id ASC LIMIT \$3$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtlr := gomock.NewController(t)
			citiesRepo := mocks.NewMockCitiesRepo(mockCtlr)
			db := mocks.NewMockDatabase(mockCtlr)
			d := NewCarsDependencies(db, citiesRepo)
			dbHandle := test.setMocks(d)

			carsRepo := NewCarsRepository(db, citiesRepo)
			cars, err := carsRepo.List(test.args.ctx, test.args.cityName, test.args.from_car_id, test.args.limit)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.cars, cars)
			assert.Equal(t, test.wants.err, err)
		})
	}
}
