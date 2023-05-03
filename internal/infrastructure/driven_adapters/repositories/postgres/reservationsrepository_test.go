package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/core/services"
	mocks "github.com/Edigiraldo/car-rent/internal/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type reservationsDependencies struct {
	db *mocks.MockDatabase
}

func NewReservationsDependencies(db *mocks.MockDatabase) *reservationsDependencies {
	return &reservationsDependencies{
		db: db,
	}
}

func TestReservationsInsert(t *testing.T) {
	initConstantsFromRepository(t)

	dr := domain.Reservation{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
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
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when user was not found",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New(services.ErrUserNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("INSERT INTO reservations").
					WithArgs(dr.ID, dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate).
					WillReturnError(&pq.Error{Code: "23503", Message: ".* user_id .*"})

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when reservation was not found",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New(services.ErrCarNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("INSERT INTO reservations").
					WithArgs(dr.ID, dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate).
					WillReturnError(&pq.Error{Code: "23503", Message: ".* car_id .*"})

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when exec context fails",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New("exec context"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("INSERT INTO reservations").
					WithArgs(dr.ID, dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate).
					WillReturnError(errors.New("exec context"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when reservation was successfully inserted",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("INSERT INTO reservations").
					WithArgs(dr.ID, dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			err := reservationsRepo.Insert(test.args.ctx, test.args.reservation)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsGet(t *testing.T) {
	initConstantsFromRepository(t)

	dr := domain.Reservation{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, 0, 7),
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type wants struct {
		reservation domain.Reservation
		err         error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when no rows were found",
			args: args{
				ctx: context.TODO(),
				id:  dr.ID,
			},
			wants: wants{
				reservation: domain.Reservation{},
				err:         errors.New(services.ErrReservationNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM reservations WHERE ID = \$1`).
					WillReturnError(sql.ErrNoRows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when query to reservation repo fails",
			args: args{
				ctx: context.TODO(),
				id:  dr.ID,
			},
			wants: wants{
				reservation: domain.Reservation{},
				err:         errors.New("there was some error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`SELECT \* FROM reservations WHERE ID = \$1`).
					WillReturnError(errors.New("there was some error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns reservation when was successfully found",
			args: args{
				ctx: context.TODO(),
				id:  dr.ID,
			},
			wants: wants{
				reservation: dr,
				err:         nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := dr.ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := dr.UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := dr.CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "reservation_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate)
				mock.ExpectQuery(`SELECT \* FROM reservations WHERE ID = \$1`).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			dr, err := reservationsRepo.Get(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.reservation, dr)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsFullUpdate(t *testing.T) {
	initConstantsFromRepository(t)

	dr := domain.Reservation{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		CarID:         uuid.New(),
		Status:        "Reserved",
		PaymentStatus: "Pending",
		StartDate:     time.Now(),
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
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when user was not found",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New(services.ErrUserNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
					WillReturnError(&pq.Error{Code: "23503", Message: ".* user_id .*"})

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when car was not found",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New(services.ErrCarNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
					WillReturnError(&pq.Error{Code: "23503", Message: ".* car_id .*"})

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when exec context fails",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New("exec context"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
					WillReturnError(errors.New("exec context"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows affected method fails",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New("rows affected error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewErrorResult(errors.New("rows affected error"))
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when reservation was not found",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: errors.New(services.ErrReservationNotFound),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
					WillReturnResult(result)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when reservation was updated",
			args: args{
				ctx:         context.TODO(),
				reservation: dr,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("UPDATE reservations SET").
					WithArgs(dr.UserID, dr.CarID, dr.Status, dr.PaymentStatus, dr.StartDate, dr.EndDate, dr.ID).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			err := reservationsRepo.FullUpdate(test.args.ctx, test.args.reservation)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsDelete(t *testing.T) {
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
		setMocks func(*reservationsDependencies) *sql.DB
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
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectExec("DELETE FROM reservations").
					WithArgs(id).
					WillReturnError(errors.New("id not found"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns nil error when reservation was deleted",
			args: args{
				ctx: context.TODO(),
				id:  id,
			},
			wants: wants{
				err: nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {
				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				result := sqlmock.NewResult(0, 1)
				mock.ExpectExec("DELETE FROM reservations").
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			err := reservationsRepo.Delete(test.args.ctx, test.args.id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsGetByUserID(t *testing.T) {
	initConstantsFromRepository(t)

	drs := []domain.Reservation{
		{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}

	type args struct {
		ctx     context.Context
		user_id uuid.UUID
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when query context fails",
			args: args{
				ctx:     context.TODO(),
				user_id: drs[0].UserID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("query context error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE user_id=\$1$`).
					WillReturnError(errors.New("query context error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows are not as expected",
			args: args{
				ctx:     context.TODO(),
				user_id: drs[0].UserID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("sql: expected 1 destination arguments in Scan, not 7"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(reservationIdByte)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE user_id=\$1$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows.Err fails",
			args: args{
				ctx:     context.TODO(),
				user_id: drs[0].UserID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("rows.Err error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate).
					RowError(0, errors.New("rows.Err error"))
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE user_id=\$1$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns reservations when were successfully found",
			args: args{
				ctx:     context.TODO(),
				user_id: drs[0].UserID,
			},
			wants: wants{
				reservations: drs,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE user_id=\$1$`).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			reservations, err := reservationsRepo.GetByUserID(test.args.ctx, test.args.user_id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsGetByCarID(t *testing.T) {
	initConstantsFromRepository(t)

	drs := []domain.Reservation{
		{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}

	type args struct {
		ctx    context.Context
		car_id uuid.UUID
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when query context fails",
			args: args{
				ctx:    context.TODO(),
				car_id: drs[0].CarID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("query context error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1$`).
					WillReturnError(errors.New("query context error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows are not as expected",
			args: args{
				ctx:    context.TODO(),
				car_id: drs[0].CarID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("sql: expected 1 destination arguments in Scan, not 7"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(reservationIdByte)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows.Err fails",
			args: args{
				ctx:    context.TODO(),
				car_id: drs[0].CarID,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("rows.Err error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate).
					RowError(0, errors.New("rows.Err error"))
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns reservations when were successfully found",
			args: args{
				ctx:    context.TODO(),
				car_id: drs[0].CarID,
			},
			wants: wants{
				reservations: drs,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1$`).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			reservations, err := reservationsRepo.GetByCarID(test.args.ctx, test.args.car_id)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}

func TestReservationsGetByCarIDAndTimeFrame(t *testing.T) {
	initConstantsFromRepository(t)

	drs := []domain.Reservation{
		{
			ID:            uuid.New(),
			UserID:        uuid.New(),
			CarID:         uuid.New(),
			Status:        "Reserved",
			PaymentStatus: "Pending",
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 0, 7),
		},
	}
	start_date := time.Now()
	end_date := time.Now().Add(9 * time.Hour)

	type args struct {
		ctx        context.Context
		car_id     uuid.UUID
		start_date time.Time
		end_date   time.Time
	}
	type wants struct {
		reservations []domain.Reservation
		err          error
	}
	tests := []struct {
		name     string
		args     args
		wants    wants
		setMocks func(*reservationsDependencies) *sql.DB
	}{
		{
			name: "returns error when query context fails",
			args: args{
				ctx:        context.TODO(),
				car_id:     drs[0].CarID,
				start_date: start_date,
				end_date:   end_date,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("query context error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1 AND start_date BETWEEN \$2 AND \$3 AND end_date BETWEEN \$2 AND \$3$`).
					WillReturnError(errors.New("query context error"))

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows are not as expected",
			args: args{
				ctx:        context.TODO(),
				car_id:     drs[0].CarID,
				start_date: start_date,
				end_date:   end_date,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("sql: expected 1 destination arguments in Scan, not 7"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(reservationIdByte)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1 AND start_date BETWEEN \$2 AND \$3 AND end_date BETWEEN \$2 AND \$3$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns error when rows.Err fails",
			args: args{
				ctx:        context.TODO(),
				car_id:     drs[0].CarID,
				start_date: start_date,
				end_date:   end_date,
			},
			wants: wants{
				reservations: nil,
				err:          errors.New("rows.Err error"),
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate).
					RowError(0, errors.New("rows.Err error"))
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1 AND start_date BETWEEN \$2 AND \$3 AND end_date BETWEEN \$2 AND \$3$`).
					WillReturnRows(rows)

				d.db.EXPECT().GetDBHandle().Return(dbHandle)

				return dbHandle
			},
		},
		{
			name: "returns reservations when were successfully found",
			args: args{
				ctx:        context.TODO(),
				car_id:     drs[0].CarID,
				start_date: start_date,
				end_date:   end_date,
			},
			wants: wants{
				reservations: drs,
				err:          nil,
			},
			setMocks: func(d *reservationsDependencies) *sql.DB {

				dbHandle, mock, err := sqlmock.New()
				if err != nil {
					t.Fatal(err)
				}
				reservationIdByte, err := drs[0].ID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				userIdByte, err := drs[0].UserID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				carIdByte, err := drs[0].CarID.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				rows := sqlmock.NewRows([]string{"id", "user_id", "car_id", "status", "payment_status", "start_date", "end_date"}).
					AddRow(reservationIdByte, userIdByte, carIdByte, drs[0].Status, drs[0].PaymentStatus, drs[0].StartDate, drs[0].EndDate)
				mock.ExpectQuery(`^SELECT \* FROM reservations WHERE car_id=\$1 AND start_date BETWEEN \$2 AND \$3 AND end_date BETWEEN \$2 AND \$3$`).
					WithArgs(drs[0].CarID, start_date, end_date).
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
			d := NewReservationsDependencies(db)
			dbHandle := test.setMocks(d)

			reservationsRepo := NewReservationsRepository(db)
			reservations, err := reservationsRepo.GetByCarIDAndTimeFrame(test.args.ctx, test.args.car_id, test.args.start_date, test.args.end_date)

			if dbHandle != nil {
				dbHandle.Close()
			}
			assert.Equal(t, test.wants.reservations, reservations)
			assert.Equal(t, test.wants.err, err)
		})
	}
}
