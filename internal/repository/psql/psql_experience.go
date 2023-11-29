package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	"HnH/pkg/queryUtils"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type IExperienceRepository interface {
	GetCVExperiencesIDs(ctx context.Context, cvID int) ([]int, error)
	GetTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) ([]domain.DbExperience, error)
	GetTxExperiencesByIds(ctx context.Context, tx *sql.Tx, cvIDs []int) ([]domain.DbExperience, error)
	AddExperience(ctx context.Context, cvID int, experience domain.DbExperience) (int, error)
	AddTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error
	UpdateTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error
	DeleteTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) error
	DeleteTxExperiencesByIDs(ctx context.Context, tx *sql.Tx, expIds []int) error
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

func (repo *psqlExperienceRepository) GetCVExperiencesIDs(ctx context.Context, cvID int) ([]int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting experiences by 'cv_id'")
	query := `SELECT 
		e.id
	FROM
		hnh_data.experience e
	WHERE
		e.cv_id = $1`

	rows, selErr := repo.DB.Query(query, cvID)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	expIDsToReturn := []int{}
	for rows.Next() {
		// exp := domain.DbExperience{}
		var expID int
		scanErr := rows.Scan(&expID)
		if scanErr != nil {
			return nil, scanErr
		}
		expIDsToReturn = append(expIDsToReturn, expID)
	}
	if len(expIDsToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return expIDsToReturn, nil
}

func (repo *psqlExperienceRepository) GetTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) ([]domain.DbExperience, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting experiences by 'cv_id' in transaction")
	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.experience e
	WHERE
		e.cv_id = $1`

	rows, selErr := tx.Query(query, cvID)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	experiencesToReturn := []domain.DbExperience{}
	for rows.Next() {
		exp := domain.DbExperience{}
		scanErr := rows.Scan(
			&exp.ID,
			&exp.CvID,
			&exp.OrganizationName,
			&exp.Position,
			&exp.Description,
			&exp.StartDate,
			&exp.EndDate,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		experiencesToReturn = append(experiencesToReturn, exp)
	}

	// if len(experiencesToReturn) == 0 {
	// 	return nil, ErrEntityNotFound
	// }
	return experiencesToReturn, nil
}

func (repo *psqlExperienceRepository) GetTxExperiencesByIds(ctx context.Context, tx *sql.Tx, cvIDs []int) ([]domain.DbExperience, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting experiences by 'cv_id' list in transaction")
	if len(cvIDs) == 0 {
		return nil, ErrEntityNotFound
	}

	placeHolderValues := *queryUtils.IntToAnySlice(cvIDs)
	placeHolderString := queryUtils.QueryPlaceHolders(1, len(placeHolderValues))

	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.experience e
	WHERE
		e.cv_id IN (` + placeHolderString + `)`

	rows, selErr := tx.Query(query, placeHolderValues...)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	experiencesToReturn := []domain.DbExperience{}
	for rows.Next() {
		exp := domain.DbExperience{}
		scanErr := rows.Scan(
			&exp.ID,
			&exp.CvID,
			&exp.OrganizationName,
			&exp.Position,
			&exp.Description,
			&exp.StartDate,
			&exp.EndDate,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		experiencesToReturn = append(experiencesToReturn, exp)
	}

	// if len(experiencesToReturn) == 0 {
	// 	return []domain.DbExperience{}, nil
	// }
	return experiencesToReturn, nil
}

func (repo *psqlExperienceRepository) AddExperience(ctx context.Context, cvID int, experience domain.DbExperience) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new experience by 'cv_id'")
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

func (repo *psqlExperienceRepository) AddTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new experiences by 'cv_id' in transaction")
	if len(experiences) == 0 {
		return nil
	}
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

// TODO: check when chenging number of experiences in cv
func (repo *psqlExperienceRepository) UpdateTxExperiences(ctx context.Context, tx *sql.Tx, cvID int, experiences []domain.DbExperience) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("updating experiences by 'cv_id' in transaction")
	if len(experiences) == 0 {
		return nil
	}
	ids := repo.getIDs(experiences)
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

func (repo *psqlExperienceRepository) DeleteTxExperiences(ctx context.Context, tx *sql.Tx, cvID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("deleting experiences by 'cv_id' in transaction")
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

func (repo *psqlExperienceRepository) DeleteTxExperiencesByIDs(ctx context.Context, tx *sql.Tx, expIds []int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("deleting experiences by 'cv_id' list in transaction")
	if len(expIds) == 0 {
		return nil
	}

	placeHoldersValues := *queryUtils.IntToAnySlice(expIds)
	queryPlaceHolders := queryUtils.QueryPlaceHolders(1, len(expIds))

	query := `DELETE
	FROM
		hnh_data.experience e
	WHERE
		e.id IN (` + queryPlaceHolders + `)`

	result, delErr := tx.Exec(query, placeHoldersValues...)

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
