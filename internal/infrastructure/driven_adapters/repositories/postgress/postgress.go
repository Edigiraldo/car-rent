package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Edigiraldo/car-rent/internal/core/domain"
	"github.com/Edigiraldo/car-rent/internal/infrastructure/driven_adapters/repositories/models"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(URI string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", URI)
	if err != nil {
		fmt.Println("starting connection with postgres:", err)
		return nil, err
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (repo *PostgresDB) InsertCar(ctx context.Context, dc domain.Car) (err error) {
	car := models.LoadCarFromDomain(dc)
	_, err = repo.db.ExecContext(ctx, "INSERT INTO cars (id, type, seats, hourly_rent_cost, city, status) VALUES ($1, $2, $3, $4, $5, $6)",
		car.ID, car.Type, car.Seats, car.HourlyRentCost, car.City, car.Status)

	return err
}
