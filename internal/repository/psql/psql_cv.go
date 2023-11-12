package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	// "fmt"
)

type ICVRepository interface {
	GetCVById(cvID int) (*domain.DbCV, error)
	GetCVsByIds(idList []int) ([]domain.DbCV, error)
	GetCVsByUserId(userID int) ([]domain.DbCV, error)
	AddCV(userID int, cv *domain.DbCV, experiences []domain.DbExperience, insitutions []domain.DbEducationInstitution) (int, error)
	GetOneOfUsersCV(userID, cvID int) (*domain.DbCV, error)
	UpdateOneOfUsersCV(userID, cvID int, cv *domain.DbCV) error
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

func (repo *psqlCVRepository) GetCVById(cvID int) (*domain.DbCV, error) {
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

	cvToReturn := domain.DbCV{}

	err := repo.DB.QueryRow(query, cvID).
		Scan(
			&cvToReturn.ID,
			&cvToReturn.ApplicantID,
			&cvToReturn.ProfessionName,
			&cvToReturn.Description,
			&cvToReturn.Status,
			&cvToReturn.CreatedAt,
			&cvToReturn.UpdatedAt,
		)
	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &cvToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByIds(idList []int) ([]domain.DbCV, error) {
	if len(idList) == 0 {
		return nil, ErrEntityNotFound
	}

	placeHolderValues := *queryUtils.IntToAnySlice(idList)
	placeHolderString := queryUtils.QueryPlaceHolders(1, placeHolderValues...)

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

	cvsToReturn := []domain.DbCV{}
	for rows.Next() {
		cv := domain.DbCV{}
		err := rows.Scan(
			&cv.ID,
			&cv.ApplicantID,
			&cv.ProfessionName,
			&cv.Description,
			&cv.Status,
			&cv.CreatedAt,
			&cv.UpdatedAt,
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

func (repo *psqlCVRepository) GetCVsByUserId(userID int) ([]domain.DbCV, error) {
	// fmt.Printf("userID: %v\n", userID)
	query := `SELECT
		c.id,
		c.applicant_id,
		c.profession,
		c.description,
		c.status,
		c.created_at,
		c.updated_at
	FROM
		hnh_data.cv c
	INNER JOIN (
			SELECT
				a.id
			FROM
				hnh_data.applicant a
			WHERE
				a.user_id = $1
		) AS w ON
		c.applicant_id = w.id`

	rows, err := repo.DB.Query(query, userID)
	if err != nil {
		// fmt.Printf("err: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	cvsToReturn := []domain.DbCV{}
	for rows.Next() {
		cv := domain.DbCV{}
		err := rows.Scan(
			&cv.ID,
			&cv.ApplicantID,
			&cv.ProfessionName,
			&cv.Description,
			&cv.Status,
			&cv.CreatedAt,
			&cv.UpdatedAt,
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

func (repo *psqlCVRepository) AddCV(
	userID int, cv *domain.DbCV,
	experiences []domain.DbExperience,
	insitutions []domain.DbEducationInstitution,
) (int, error) {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return 0, txErr
	}

	cvInsertQuery := `INSERT
		INTO
		hnh_data.cv (
			applicant_id,
			profession,
			first_name,
			last_name,
			middle_name,
			gender,
			birthday,
			location,
			description,
			education_level
		)
	SELECT
		a.id, $1, $2, $3, $4, $5, $6, $7, $8, $9
	FROM
		hnh_data.applicant a
	WHERE
		a.user_id = $10
	RETURNING id`

	var insertedCVID int
	insertCvErr := tx.QueryRow(
		cvInsertQuery,
		cv.ProfessionName,
		cv.FirstName,
		cv.LastName,
		cv.MiddleName,
		cv.Gender,
		cv.Birthday,
		cv.Location,
		cv.Description,
		cv.EducationLevel,
		userID,
	).
		Scan(&insertedCVID)

	if insertCvErr == sql.ErrNoRows {
		tx.Rollback()
		return 0, ErrNotInserted
	}
	if insertCvErr != nil {
		tx.Rollback()
		return 0, insertCvErr
	}

	// TODO: change to one query not in loop
	insertExpQuery := `INSERT
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
		($1, $2, $3, $4, $5, $6)`
	for _, experience := range experiences {
		_, insertExpErr := tx.Exec(
			insertExpQuery,
			insertedCVID,
			experience.OrganizationName,
			experience.Position,
			experience.Description,
			experience.StartDate,
			experience.EndDate,
		)
		if insertExpErr == sql.ErrNoRows {
			tx.Rollback()
			return 0, ErrNotInserted
		}
		if insertExpErr != nil {
			tx.Rollback()
			return 0, insertExpErr
		}
	}

	// TODO: change to one query not in loop
	insertInstQuery := `INSERT
		INTO
		hnh_data.education_institution (
			cv_id,
			"name",
			major_field,
			graduation_year
		)
	VALUES ($1, $2, $3, $4)`
	for _, institution := range insitutions {
		_, insertInstErr := tx.Exec(
			insertInstQuery,
			insertedCVID,
			institution.Name,
			institution.MajorField,
			institution.GraduationYear,
		)
		if insertInstErr == sql.ErrNoRows {
			tx.Rollback()
			return 0, ErrNotInserted
		}
		if insertInstErr != nil {
			tx.Rollback()
			return 0, insertInstErr
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return 0, commitErr
	}

	return insertedCVID, nil
}

func (repo *psqlCVRepository) GetOneOfUsersCV(userID, cvID int) (*domain.DbCV, error) {
	query := `SELECT
		c.id,
		c.applicant_id,
		c.profession,
		c.description,
		c.status,
		c.created_at,
		c.updated_at
	FROM
		hnh_data.cv c
	INNER JOIN (
			SELECT
				a.id
			FROM
				hnh_data.applicant a
			WHERE
				a.user_id = $1
		) AS w ON
		c.applicant_id = w.id
	WHERE
		c.id = $2`

	var cvToReturn = domain.DbCV{}
	err := repo.DB.QueryRow(
		query,
		userID,
		cvID,
	).
		Scan(
			&cvToReturn.ID,
			&cvToReturn.ApplicantID,
			&cvToReturn.ProfessionName,
			&cvToReturn.Description,
			&cvToReturn.Status,
			&cvToReturn.CreatedAt,
			&cvToReturn.UpdatedAt,
		)

	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &cvToReturn, nil
}

func (repo *psqlCVRepository) UpdateOneOfUsersCV(userID, cvID int, cv *domain.DbCV) error {
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
	if err == sql.ErrNoRows {
		return ErrNoRowsUpdated
	}
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
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
	if err == sql.ErrNoRows {
		return ErrNoRowsDeleted
	}
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
