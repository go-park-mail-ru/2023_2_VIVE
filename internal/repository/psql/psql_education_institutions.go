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

type IEducationInstitutionRepository interface {
	GetCVInstitutionsIDs(ctx context.Context, cvID int) ([]int, error)
	GetTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int) ([]domain.DbEducationInstitution, error)
	GetTxInstitutionsByIds(ctx context.Context, tx *sql.Tx, cvIDs []int) ([]domain.DbEducationInstitution, error)
	AddTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error
	UpdateTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error
	DeleteTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int) error
	DeleteTxInstitutionsByIDs(ctx context.Context, tx *sql.Tx, instIDs []int) error
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

func (repo *psqlEducationInstitutionRepository) GetCVInstitutionsIDs(ctx context.Context, cvID int) ([]int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting intitutions by 'cv_id'")
	query := `SELECT 
		ei.id
	FROM
		hnh_data.education_institution ei
	WHERE
		ei.cv_id = $1`

	rows, selErr := repo.DB.Query(query, cvID)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	instIDsToReturn := []int{}
	for rows.Next() {
		var instID int
		scanErr := rows.Scan(&instID)
		if scanErr != nil {
			return nil, scanErr
		}
		instIDsToReturn = append(instIDsToReturn, instID)
	}
	if len(instIDsToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return instIDsToReturn, nil
}

func (repo *psqlEducationInstitutionRepository) GetTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int) ([]domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting intitutions by 'cv_id' in transaction")
	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.education_institution ei
	WHERE
		ei.cv_id = $1`

	rows, selErr := tx.Query(query, cvID)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	institutionsToReturn := []domain.DbEducationInstitution{}
	for rows.Next() {
		edInst := domain.DbEducationInstitution{}
		scanErr := rows.Scan(
			&edInst.ID,
			&edInst.CvID,
			&edInst.Name,
			&edInst.MajorField,
			&edInst.GraduationYear,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		institutionsToReturn = append(institutionsToReturn, edInst)
	}

	// if len(institutionsToReturn) == 0 {
	// 	return nil, ErrEntityNotFound
	// }
	return institutionsToReturn, nil
}

func (repo *psqlEducationInstitutionRepository) GetTxInstitutionsByIds(ctx context.Context, tx *sql.Tx, cvIDs []int) ([]domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting intitutions list by 'cv_id' list in transaction")
	if len(cvIDs) == 0 {
		return nil, ErrEntityNotFound
	}

	placeHolderValues := *queryUtils.IntToAnySlice(cvIDs)
	placeHolderString := queryUtils.QueryPlaceHolders(1, len(placeHolderValues))

	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.education_institution ei
	WHERE
		ei.cv_id IN (` + placeHolderString + `)`

	rows, selErr := tx.Query(query, placeHolderValues...)
	if selErr != nil {
		return nil, selErr
	}
	defer rows.Close()

	institutionsToReturn := []domain.DbEducationInstitution{}
	for rows.Next() {
		inst := domain.DbEducationInstitution{}
		scanErr := rows.Scan(
			&inst.ID,
			&inst.CvID,
			&inst.Name,
			&inst.MajorField,
			&inst.GraduationYear,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		institutionsToReturn = append(institutionsToReturn, inst)
	}

	// if len(institutionsToReturn) == 0 {
	// 	return []domain.DbEducationInstitution{}, nil
	// }
	return institutionsToReturn, nil
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

func (repo *psqlEducationInstitutionRepository) AddTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new intitutions by 'cv_id' in transaction")
	if len(institutions) == 0 {
		return nil
	}
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

func (repo *psqlEducationInstitutionRepository) UpdateTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int, institutions []domain.DbEducationInstitution) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("updating intitutions by 'cv_id' in transaction")
	if len(institutions) == 0 {
		return nil
	}
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

func (repo *psqlEducationInstitutionRepository) DeleteTxInstitutions(ctx context.Context, tx *sql.Tx, cvID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("deleting intitutions by 'cv_id' in transaction")
	query := `DELETE
	FROM
		hnh_data.education_institution ei
	WHERE
		ei.cv_id = $1`

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

func (repo *psqlEducationInstitutionRepository) DeleteTxInstitutionsByIDs(ctx context.Context, tx *sql.Tx, instIDs []int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("deleting intitutions by 'cv_id' list in transaction")
	if len(instIDs) == 0 {
		return nil
	}

	placeHoldersValues := *queryUtils.IntToAnySlice(instIDs)
	queryPlaceHolders := queryUtils.QueryPlaceHolders(1, len(instIDs))

	query := `DELETE
	FROM
		hnh_data.education_institution ei
	WHERE
		ei.id IN (` + queryPlaceHolders + `)`

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
