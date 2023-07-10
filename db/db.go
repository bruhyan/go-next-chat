package db

import (
	"database/sql"
	_ "github.com/lib/pq" // postgres driver
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	driverName := "postgres"
	dataSourceName := "postgresql://root:password@localhost:5433/go-chat?sslmode=disable"

	db, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// Close receiver function to close db connection
func (d *Database) Close() {
	_ = d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
