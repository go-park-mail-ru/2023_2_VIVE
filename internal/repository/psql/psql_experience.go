package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
)

type IExperienceRepository interface {
	AddExperience(cvID int, experience domain.DbExperience) (int, error)
	AddTxExperiences(tx *sql.Tx, cvID int, experiences []domain.DbExperience) error
}

type psqlExperienceRepository struct {
	DB *sql.DB
}

func NewPsqlExperienceRepository(db *sql.DB) IExperienceRepository {
	return &psqlExperienceRepository{
		DB: db,
	}
}

func (repo *psqlExperienceRepository) AddExperience(cvID int, experience domain.DbExperience) (int, error) {
	query := `INSERT
		INTO
		hnh_data.experience (
			cv_id,
			organization_name,
			"position",
			description,
			start_date,
			end_date
		)
	VALUES 
		($1, $2, $3, $4, $5, $6)
	RETURNING id`

	var insertedExpID int
	insertErr := repo.DB.QueryRow(
		query,
		cvID,
		experience.OrganizationName,
		experience.Position,
		experience.Description,
		experience.StartDate,
		experience.EndDate,
	).
		Scan(insertedExpID)

	if insertErr == sql.ErrNoRows {
		return 0, ErrNotInserted
	}
	if insertErr != nil {
		return 0, insertErr
	}

	return insertedExpID, nil
}

func (repo *psqlExperienceRepository) convertToSlice(cvID int, experiences []domain.DbExperience) []any {
	res := []any{}
	for _, experience := range experiences {

		res = append(res, cvID)
		res = append(res, experience.OrganizationName)
		res = append(res, experience.Position)
		res = append(res, experience.Description)
		res = append(res, experience.StartDate)
		if experience.EndDate == nil {
			res = append(res, nil)
		} else {
			res = append(res, *experience.EndDate)
		}
	}
	return res
}

func (repo *psqlExperienceRepository) AddTxExperiences(tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	elementsToInsert := repo.convertToSlice(cvID, experiences)
	query := `INSERT
		INTO
		hnh_data.experience (
			cv_id,
			organization_name,
			"position",
			description,
			start_date,
			end_date
		)
	VALUES ` + queryUtils.QueryPlaceHoldersMultipleRows(1, 6, len(experiences))

	result, insertErr := tx.Exec(query, elementsToInsert...)
	if insertErr == sql.ErrNoRows {
		return ErrNotInserted
	}
	if insertErr != nil {
		return insertErr
	}
	_, insertErr = result.RowsAffected()
	if insertErr != nil {
		return insertErr
	}

	return nil
}
