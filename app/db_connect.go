package app

import (
	"HnH/configs"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
)

func getPostgres() (*sql.DB, error) {
	dsn := configs.HnHPostgresConfig.GetConnectionString()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
