package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	"fmt"
)

type ICVRepository interface {
	GetCVById(cvID int) (*domain.DbCV, error)
	GetCVsByIds(idList []int) ([]domain.DbCV, error)
	GetCVsByUserId(userID int) ([]domain.DbCV, error)
	AddCV(userID int, cv *domain.DbCV, experiences []domain.DbExperience, insitutions []domain.DbEducationInstitution) (int, error)
	GetOneOfUsersCV(userID, cvID int) (*domain.DbCV, error)
	UpdateOneOfUsersCV(userID, cvID int, cv *domain.DbCV, experiences []domain.DbExperience, insitutions []domain.DbEducationInstitution) error
	DeleteOneOfUsersCV(userID, cvID int) error
}

type psqlCVRepository struct {
	DB       *sql.DB
	expRepo  IExperienceRepository
	instRepo IEducationInstitutionRepository
}

func NewPsqlCVRepository(db *sql.DB) ICVRepository {
	return &psqlCVRepository{
		DB:       db,
		expRepo:  NewPsqlExperienceRepository(db),
		instRepo: NewPsqlEducationInstitutionRepository(db),
	}
}

func (repo *psqlCVRepository) GetCVById(cvID int) (*domain.DbCV, error) {
	query := `SELECT
		id,
		applicant_id,
		profession,
		first_name,
		last_name,
		middle_name,
		gender,
		birthday,
		location,
		description,
		education_level,
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
			&cvToReturn.FirstName,
			&cvToReturn.LastName,
			&cvToReturn.MiddleName,
			&cvToReturn.Gender,
			&cvToReturn.Birthday,
			&cvToReturn.Location,
			&cvToReturn.Description,
			&cvToReturn.EducationLevel,
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
	placeHolderString := queryUtils.QueryPlaceHolders(1, len(placeHolderValues))

	query := `SELECT
			id,
			applicant_id,
			profession,
			first_name,
			last_name,
			middle_name,
			gender,
			birthday,
			location,
			description,
			education_level,
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
			&cv.FirstName,
			&cv.LastName,
			&cv.MiddleName,
			&cv.Gender,
			&cv.Birthday,
			&cv.Location,
			&cv.Description,
			&cv.EducationLevel,
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
	query := `SELECT
		c.id,
		c.applicant_id,
		c.profession,
		c.first_name,
		c.last_name,
		c.middle_name,
		c.gender,
		c.birthday,
		c.location,
		c.description,
		c.education_level,
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
			&cv.FirstName,
			&cv.LastName,
			&cv.MiddleName,
			&cv.Gender,
			&cv.Birthday,
			&cv.Location,
			&cv.Description,
			&cv.EducationLevel,
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

	query := `INSERT
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
		query,
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
	expErr := repo.expRepo.AddTxExperiences(tx, insertedCVID, experiences)
	if expErr != nil {
		tx.Rollback()
		return 0, expErr
	}

	instErr := repo.instRepo.AddTxInstitutions(tx, insertedCVID, insitutions)
	if instErr != nil {
		tx.Rollback()
		return 0, instErr
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
		c.first_name,
		c.last_name,
		c.middle_name,
		c.gender,
		c.birthday,
		c.location,
		c.description,
		c.education_level,
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
			&cvToReturn.FirstName,
			&cvToReturn.LastName,
			&cvToReturn.MiddleName,
			&cvToReturn.Gender,
			&cvToReturn.Birthday,
			&cvToReturn.Location,
			&cvToReturn.Description,
			&cvToReturn.EducationLevel,
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

func (repo *psqlCVRepository) UpdateOneOfUsersCV(
	userID, cvID int,
	cv *domain.DbCV,
	experiences []domain.DbExperience,
	insitutions []domain.DbEducationInstitution,
) error {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return txErr
	}
	query := `UPDATE
		hnh_data.cv c
	SET 
		profession = $1,
		first_name = $2,
		last_name = $3,
		middle_name = $4,
		gender = $5,
		birthday = $6,
		location = $7,
		description = $8, 
		status = $9,
		education_level = $10,
		updated_at = now()
	FROM hnh_data.applicant a
	WHERE 
		c.id = $11
		AND a.user_id = $12
		AND c.applicant_id = a.id`

	result, err := tx.Exec(
		query,
		cv.ProfessionName,
		cv.FirstName,
		cv.LastName,
		cv.MiddleName,
		cv.Gender,
		cv.Birthday,
		cv.Location,
		cv.Description,
		cv.Status,
		cv.EducationLevel,
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
	fmt.Printf("after update cv\n")

	expErr := repo.expRepo.UpdateTxExperiences(tx, cvID, experiences)
	if expErr != nil {
		tx.Rollback()
		return expErr
	}
	fmt.Printf("after update exp\n")

	instErr := repo.instRepo.UpdateTxInstitutions(tx, cvID, insitutions)
	if instErr != nil {
		tx.Rollback()
		return instErr
	}
	fmt.Printf("after update inst\n")

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
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
