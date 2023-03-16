package postgres

import "database/sql"

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(URI string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", URI)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
