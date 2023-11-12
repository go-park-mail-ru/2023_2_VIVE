package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
)

type IEducationInstitutionRepository interface {
	AddTxInstitutions(tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error
}

type psqlEducationInstitutionRepository struct {
	DB *sql.DB
}

func NewPsqlEducationInstitutionRepository(db *sql.DB) IEducationInstitutionRepository {
	return &psqlEducationInstitutionRepository{
		DB: db,
	}
}

func (repo *psqlEducationInstitutionRepository) convertToSlice(cvID int, institutions []domain.DbEducationInstitution) []any {
	res := []any{}
	for _, institution := range institutions {
		res = append(res, cvID)
		res = append(res, institution.Name)
		res = append(res, institution.MajorField)
		res = append(res, institution.GraduationYear)
	}
	return res
}

func (repo *psqlEducationInstitutionRepository) AddTxInstitutions(tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error {
	query := `INSERT
		INTO
		hnh_data.education_institution (
			cv_id,
			"name",
			major_field,
			graduation_year
		)
	VALUES ` + queryUtils.QueryPlaceHoldersMultipleRows(1, 4, len(institutions))

	result, insertErr := tx.Exec(query, repo.convertToSlice(cvID, institutions)...)
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
