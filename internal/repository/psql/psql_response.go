package psql

import (
	"database/sql"
	"errors"
)

type IResponseRepository interface {
	RespondToVacancy(vacancyID, cvID int) error
	GetVacanciesIdsByCVId(cvID int) ([]int, error)
	GetAttachedCVs(vacancyID int) ([]int, error)
}

type psqlResponseRepository struct {
	responseStorage *sql.DB
}

func NewPsqlResponseRepository(db *sql.DB) IResponseRepository {
	return &psqlResponseRepository{
		responseStorage: db,
	}
}

func (p *psqlResponseRepository) RespondToVacancy(vacancyID, cvID int) error {
	err := p.responseStorage.QueryRow(`INSERT INTO hnh_data.response ("vacancy_id", "cv_id") VALUES ($1, $2)`, vacancyID, cvID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlResponseRepository) GetVacanciesIdsByCVId(cvID int) ([]int, error) {
	rows, err := p.responseStorage.Query(`SELECT vacancy_id FROM hnh_data.response WHERE cv_id = $1`, cvID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []int{}
	for rows.Next() {
		var vacancyID int

		err := rows.Scan(&vacancyID)
		if err != nil {
			return nil, err
		}
		result = append(result, vacancyID)
	}

	return result, nil
}

func (p *psqlResponseRepository) GetAttachedCVs(vacancyID int) ([]int, error) {
	rows, err := p.responseStorage.Query(`SELECT cv_id FROM hnh_data.response WHERE vacancy_id = $1`, vacancyID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []int{}
	for rows.Next() {
		var cvID int

		err := rows.Scan(&cvID)
		if err != nil {
			return nil, err
		}
		result = append(result, cvID)
	}

	return result, nil
}
