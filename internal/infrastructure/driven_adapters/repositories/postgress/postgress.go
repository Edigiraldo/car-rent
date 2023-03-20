package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func (p *PostgresDB) GetDBHandle() *sql.DB {
	return p.db
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
