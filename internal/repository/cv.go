package repository

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
)

type ICVRepository interface {
	GetCVById(cvID int) (*domain.CV, error)
	GetCVsByIds(idList []int) ([]domain.CV, error)
	GetCVsByUserId(userID int) ([]domain.CV, error)
	AddCV(userID int, cv *domain.CV) (int, error)
	GetOneOfUsersCV(userID, cvID int) (*domain.CV, error)
	UpdateOneOfUsersCV(userID, cvID int, cv *domain.CV) error
	DeleteOneOfUsersCV(userID, cvID int) error
}

type psqlCVRepository struct {
	DB *sql.DB
}

func NewPsqlCVRepository(db *sql.DB) ICVRepository {
	return &psqlCVRepository{
		DB: db,
	}
}

func (repo *psqlCVRepository) GetCVById(cvID int) (*domain.CV, error) {
	query := `SELECT
		id,
		applicant_id,
		profession,
		description,
		status,
		created_at,
		updated_at
	FROM
		hnh_data.cv c
	WHERE
		c.id = $1`

	cvToReturn := domain.CV{}

	err := repo.DB.QueryRow(query, cvID).
		Scan(
			&cvToReturn.ID,
			&cvToReturn.ApplicantID,
			&cvToReturn.Status,
			&cvToReturn.Created_at,
			&cvToReturn.Updated_at,
		)
	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &cvToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByIds(idList []int) ([]domain.CV, error) {
	placeHolderValues := *queryUtils.IntToAnySlice(idList)
	placeHolderString := queryUtils.QueryPlaceHolders(placeHolderValues...)

	query := `SELECT
		id,
		applicant_id,
		profession,
		description,
		status,
		created_at,
		updated_at
	FROM
		hnh_data.cv c
	WHERE
		c.id IN (` + placeHolderString + `)`

	rows, err := repo.DB.Query(query, placeHolderValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cvsToReturn := []domain.CV{}
	for rows.Next() {
		cv := domain.CV{}
		err := rows.Scan(
			&cv.ID,
			&cv.ApplicantID,
			&cv.ProfessionName,
			&cv.Description,
			&cv.Status,
			&cv.Created_at,
			&cv.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return cvsToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByUserId(userID int) ([]domain.CV, error) {
	query := `SELECT
		id,
		applicant_id,
		profession,
		description,
		status,
		created_at,
		updated_at
	FROM
		hnh_data.cv c
	INNER JOIN (
			SELECT
				id
			FROM
				hnh_data.applicant a
			WHERE
				a.user_id = $1
		) AS w ON
		c.applicant_id = w.id`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cvsToReturn := []domain.CV{}
	for rows.Next() {
		cv := domain.CV{}
		err := rows.Scan(
			&cv.ID,
			&cv.ApplicantID,
			&cv.ProfessionName,
			&cv.Description,
			&cv.Status,
			&cv.Created_at,
			&cv.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return cvsToReturn, nil
}

func (repo *psqlCVRepository) AddCV(userID int, cv *domain.CV) (int, error) {
	query := `INSERT
		INTO
		hnh_data.cv (
			applicant_id,
			profession,
			description,
			status
		)
	SELECT
		a.id,
		$1,
		$2,
		$3
	FROM
		hnh_data.applicant a
	WHERE
		a.user_id = $4
	RETURNING id`

	var insertedCVID int
	err := repo.DB.QueryRow(
		query,
		cv.ProfessionName,
		cv.Description,
		cv.Status,
		cv.Status,
		userID,
	).
		Scan(&insertedCVID)

	if err == sql.ErrNoRows {
		return 0, ErrNotInserted
	}
	if err != nil {
		return 0, err
	}

	return insertedCVID, nil
}

func (repo *psqlCVRepository) GetOneOfUsersCV(userID, cvID int) (*domain.CV, error) {
	query := `SELECT
		id,
		applicant_id,
		profession,
		description,
		status,
		created_at,
		updated_at
	FROM
		hnh_data.cv c
	INNER JOIN (
			SELECT
				id
			FROM
				hnh_data.applicant a
			WHERE
				a.user_id = $1
		) AS w ON
		c.applicant_id = w.id
	WHERE
		c.id = $2`

	var cvToReturn *domain.CV
	err := repo.DB.QueryRow(
		query,
		userID,
		cvID,
	).
		Scan(&cvToReturn)

	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return cvToReturn, nil
}

func (repo *psqlCVRepository) UpdateOneOfUsersCV(userID, cvID int, cv *domain.CV) error {
	query := `UPDATE
		hnh_data.cv c
	SET 
		profession = $1, 
		description = $2, 
		status = $3,
		updated_at = now()
	FROM hnh_data.applicant a
	WHERE 
		c.id = $4
		AND a.user_id = $5
		AND c.applicant_id = a.id`

	result, err := repo.DB.Exec(
		query,
		cv.ProfessionName,
		cv.Description,
		cv.Status,
		cvID,
		userID,
	)
	if err != nil {
		return err
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return ErrNoRowsUpdated
	}

	return nil
}

func (repo *psqlCVRepository) DeleteOneOfUsersCV(userID, cvID int) error {
	query := `DELETE
	FROM
		hnh_data.cv c
			USING hnh_data.applicant a
	WHERE
		c.id = $1
		AND a.user_id = $2
		AND c.applicant_id = a.id`

	result, err := repo.DB.Exec(query, cvID, userID)
	if err != nil {
		return err
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return ErrNoRowsDeleted
	}

	return nil
}
