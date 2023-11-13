package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	"fmt"
	"strings"
)

type IEducationInstitutionRepository interface {
	AddTxInstitutions(tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error
	UpdateTxInstitutions(tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error
}

type psqlEducationInstitutionRepository struct {
	DB          *sql.DB
	ColumnNames []string
}

func NewPsqlEducationInstitutionRepository(db *sql.DB) IEducationInstitutionRepository {
	return &psqlEducationInstitutionRepository{
		DB: db,
		ColumnNames: []string{
			"id",
			"cv_id",
			`"name"`,
			"major_field",
			"graduation_year",
		},
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
		hnh_data.education_institution (` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames, "id"), ", ") + `)
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

func (repo *psqlEducationInstitutionRepository) getIDs(institutions []domain.DbEducationInstitution) []string {
	res := []string{}
	for _, insinstitution := range institutions {
		res = append(res, fmt.Sprint(insinstitution.ID))
	}
	return res
}

func (repo *psqlEducationInstitutionRepository) getValues(cvID int, institutions []domain.DbEducationInstitution) []any {
	res := []any{}
	valuesMap := map[string][]any{}

	for _, institution := range institutions {
		valuesMap[`"name"`] = append(valuesMap[`"name"`], institution.Name)
		valuesMap[`major_field`] = append(valuesMap[`major_field`], institution.MajorField)
		valuesMap[`graduation_year`] = append(valuesMap[`graduation_year`], institution.GraduationYear)
	}

	res = append(res, cvID)
	res = append(res, valuesMap[`"name"`]...)
	res = append(res, valuesMap[`major_field`]...)
	res = append(res, valuesMap[`graduation_year`]...)
	return res
}

func (repo *psqlEducationInstitutionRepository) UpdateTxInstitutions(tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error {
	ids := repo.getIDs(institutions)
	query := `UPDATE hnh_data.education_institution ei
	SET ` + queryUtils.QueryCases(
		2,
		[]string{`"name"`, "major_field", "graduation_year"},
		ids,
		"id") + ` WHERE ei.cv_id = $1`

	result, updErr := tx.Exec(
		query,
		repo.getValues(cvID, institutions)...,
	)

	if updErr == sql.ErrNoRows {
		return ErrNoRowsUpdated
	}
	if updErr != nil {
		return updErr
	}
	_, updErr = result.RowsAffected()
	if updErr != nil {
		return updErr
	}

	return nil
}
