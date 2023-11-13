package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	"fmt"
	"strings"
)

type IExperienceRepository interface {
	AddExperience(cvID int, experience domain.DbExperience) (int, error)
	AddTxExperiences(tx *sql.Tx, cvID int, experiences []domain.DbExperience) error
	UpdateTxExperiences(tx *sql.Tx, cvID int, experiences []domain.DbExperience) error
	DeleteTxExperiences(tx *sql.Tx, cvID int) error
}

type psqlExperienceRepository struct {
	DB          *sql.DB
	ColumnNames []string
}

func NewPsqlExperienceRepository(db *sql.DB) IExperienceRepository {
	return &psqlExperienceRepository{
		DB: db,
		ColumnNames: []string{
			"id",
			"cv_id",
			"organization_name",
			`"position"`,
			"description",
			"start_date",
			"end_date",
		},
	}
}

func (repo *psqlExperienceRepository) AddExperience(cvID int, experience domain.DbExperience) (int, error) {
	query := `INSERT
		INTO
		hnh_data.experience (` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames, "id"), ", ") + `)
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
		hnh_data.experience (` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames, "id"), ", ") + `)
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

func (repo *psqlExperienceRepository) getIDs(experiences []domain.DbExperience) []string {
	res := []string{}
	for _, experience := range experiences {
		res = append(res, fmt.Sprint(experience.ID))
	}
	return res
}

func (repo *psqlExperienceRepository) getValues(cvID int, experiences []domain.DbExperience) []any {
	res := []any{}
	valuesMap := map[string][]any{}

	for _, experience := range experiences {
		valuesMap["organization_name"] = append(valuesMap["organization_name"], experience.OrganizationName)
		valuesMap[`"position"`] = append(valuesMap[`"position"`], experience.Position)
		valuesMap["description"] = append(valuesMap["description"], experience.Description)
		valuesMap["start_date"] = append(valuesMap["start_date"], experience.StartDate)
		valuesMap["end_date"] = append(valuesMap["end_date"], experience.EndDate)
	}

	res = append(res, cvID)
	res = append(res, valuesMap["organization_name"]...)
	res = append(res, valuesMap[`"position"`]...)
	res = append(res, valuesMap["description"]...)
	res = append(res, valuesMap["start_date"]...)
	res = append(res, valuesMap["end_date"]...)
	return res
}

func (repo *psqlExperienceRepository) UpdateTxExperiences(tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	ids := repo.getIDs(experiences)
	// fmt.Printf("ids: %v\n", ids)
	query := `UPDATE hnh_data.experience e
	SET ` + queryUtils.QueryCases(
		2,
		[]string{
			"organization_name",
			`"position"`,
			"description",
			"start_date",
			"end_date",
		},
		ids,
		"id") + ` WHERE e.cv_id = $1`

	// fmt.Println(query)
	// newPlaceHolderValues := make([]any, len())
	result, updErr := tx.Exec(
		query,
		repo.getValues(cvID, experiences)...,
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

func (repo *psqlExperienceRepository) DeleteTxExperiences(tx *sql.Tx, cvID int) error {
	query := `DELETE
	FROM
		hnh_data.experience e
	WHERE
		e.cv_id = $1`

	// fmt.Println(query)
	// newPlaceHolderValues := make([]any, len())
	result, delErr := tx.Exec(query, cvID)

	if delErr == sql.ErrNoRows {
		return ErrNoRowsDeleted
	}
	if delErr != nil {
		return delErr
	}
	_, delErr = result.RowsAffected()
	if delErr != nil {
		return delErr
	}

	return nil
}
