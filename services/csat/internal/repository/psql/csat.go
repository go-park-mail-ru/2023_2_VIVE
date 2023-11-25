package psql

import (
	"context"
	"database/sql"
	"time"
)

type ICsatRepository interface {
	GetLastUpdate(ctx context.Context, userID int64) (time.Time, error)
	GetQuestions(ctx context.Context) ([]string, error)
}

type psqlCsatRepository struct {
	DB *sql.DB
}

func NewPsqlCsatRepository(db *sql.DB) ICsatRepository {
	return &psqlCsatRepository{
		DB: db,
	}
}

func (repo *psqlCsatRepository) GetLastUpdate(ctx context.Context, userID int64) (time.Time, error) {
	return time.Now(), nil
}

func (repo *psqlCsatRepository) GetQuestions(ctx context.Context) ([]string, error) {
	return []string{}, nil
}
