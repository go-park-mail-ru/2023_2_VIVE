package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	"strings"
)

type ICVRepository interface {
	GetCVById(cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	GetCVsByIds(idList []int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	GetCVsByUserId(userID int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	AddCV(userID int, cv *domain.DbCV, experiences []domain.DbExperience, insitutions []domain.DbEducationInstitution) (int, error)
	GetOneOfUsersCV(userID, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	UpdateOneOfUsersCV(userID, cvID int, cv *domain.DbCV,
		experiencesIDsToDelete []int, experiencesToUpdate, experiencesToInsert []domain.DbExperience,
		insitutionsIDsToDelete []int, insitutionsToUpdate, insitutionsToInsert []domain.DbEducationInstitution,
	) error
	DeleteOneOfUsersCV(userID, cvID int) error
}

type psqlCVRepository struct {
	DB          *sql.DB
	expRepo     IExperienceRepository
	instRepo    IEducationInstitutionRepository
	ColumnNames []string
}

func NewPsqlCVRepository(db *sql.DB) ICVRepository {
	return &psqlCVRepository{
		DB:       db,
		expRepo:  NewPsqlExperienceRepository(db),
		instRepo: NewPsqlEducationInstitutionRepository(db),
		ColumnNames: []string{
			"id",
			"applicant_id",
			"profession",
			"first_name",
			"last_name",
			"middle_name",
			"gender",
			"birthday",
			"location",
			"description",
			"education_level",
			"status",
			"created_at",
			"updated_at",
		},
	}
}

func (repo *psqlCVRepository) GetCVById(cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return nil, nil, nil, txErr
	}
	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.cv c
	WHERE
		c.id = $1`

	cvToReturn := domain.DbCV{}

	err := tx.QueryRow(query, cvID).
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
		return nil, nil, nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, nil, nil, err
	}

	expToReturn, expErr := repo.expRepo.GetTxExperiences(tx, cvID)
	if expErr != nil {
		tx.Rollback()
		return nil, nil, nil, expErr
	}

	edInstToReturn, instErr := repo.instRepo.GetTxInstitutions(tx, cvID)
	if instErr != nil {
		tx.Rollback()
		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return &cvToReturn, expToReturn, edInstToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByIds(idList []int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return nil, nil, nil, txErr
	}
	if len(idList) == 0 {
		return nil, nil, nil, ErrEntityNotFound
	}

	placeHolderValues := *queryUtils.IntToAnySlice(idList)
	placeHolderString := queryUtils.QueryPlaceHolders(1, len(placeHolderValues))

	query := `SELECT ` +
		strings.Join(queryUtils.GetColumnNames(repo.ColumnNames), ", ") +
		` FROM
		hnh_data.cv c
	WHERE
		c.id IN (` + placeHolderString + `)`

	cvRows, err := tx.Query(query, placeHolderValues...)
	if err != nil {
		return nil, nil, nil, err
	}
	defer cvRows.Close()

	cvsToReturn := []domain.DbCV{}
	cvIDs := []int{}
	for cvRows.Next() {
		cv := domain.DbCV{}
		err := cvRows.Scan(
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
			return nil, nil, nil, err
		}
		cvIDs = append(cvIDs, cv.ID)
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		return nil, nil, nil, ErrEntityNotFound
	}

	expsToReturn, expErr := repo.expRepo.GetTxExperiencesByIds(tx, cvIDs)
	if expErr != nil {
		tx.Rollback()
		return nil, nil, nil, expErr
	}
	instsToReturn, instErr := repo.instRepo.GetTxExperiencesByIds(tx, cvIDs)
	if instErr != nil {
		tx.Rollback()
		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return cvsToReturn, expsToReturn, instsToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByUserId(userID int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return nil, nil, nil, txErr
	}
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

	rows, err := tx.Query(query, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	cvsToReturn := []domain.DbCV{}
	cvIDs := []int{}
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
			return nil, nil, nil, err
		}
		cvIDs = append(cvIDs, cv.ID)
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		return nil, nil, nil, ErrEntityNotFound
	}

	// fmt.Printf("after cv select\n")
	// fmt.Printf("cvIDs: %v\n", cvIDs)

	expsToReturn, expErr := repo.expRepo.GetTxExperiencesByIds(tx, cvIDs)
	if expErr != nil {
		tx.Rollback()
		return nil, nil, nil, expErr
	}
	// fmt.Printf("after exp select\n")

	instsToReturn, instErr := repo.instRepo.GetTxExperiencesByIds(tx, cvIDs)
	if instErr != nil {
		tx.Rollback()
		return nil, nil, nil, instErr
	}
	// fmt.Printf("after inst select\n")

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return cvsToReturn, expsToReturn, instsToReturn, nil
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
	// fmt.Println("after cv insert")
	expErr := repo.expRepo.AddTxExperiences(tx, insertedCVID, experiences)
	if expErr != nil {
		tx.Rollback()
		return 0, expErr
	}
	// fmt.Println("after exp insert")
	instErr := repo.instRepo.AddTxInstitutions(tx, insertedCVID, insitutions)
	if instErr != nil {
		tx.Rollback()
		return 0, instErr
	}
	// fmt.Println("after inst insert")

	commitErr := tx.Commit()
	if commitErr != nil {
		return 0, commitErr
	}

	return insertedCVID, nil
}

func (repo *psqlCVRepository) GetOneOfUsersCV(userID, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return nil, nil, nil, txErr
	}
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
		return nil, nil, nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, nil, nil, err
	}

	exps, expErr := repo.expRepo.GetTxExperiences(tx, cvID)
	if expErr != nil {
		tx.Rollback()
		return nil, nil, nil, expErr
	}

	insts, instErr := repo.instRepo.GetTxInstitutions(tx, cvID)
	if instErr != nil {
		tx.Rollback()
		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return &cvToReturn, exps, insts, nil
}

func (repo *psqlCVRepository) UpdateOneOfUsersCV(
	userID, cvID int, cv *domain.DbCV,
	experiencesIDsToDelete []int, experiencesToUpdate, experiencesToInsert []domain.DbExperience,
	insitutionsIDsToDelete []int, insitutionsToUpdate, insitutionsToInsert []domain.DbEducationInstitution,
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
	// fmt.Printf("after update cv\n")

	// updating experiences
	expUpdErr := repo.expRepo.UpdateTxExperiences(tx, cvID, experiencesToUpdate)
	if expUpdErr != nil {
		tx.Rollback()
		return expUpdErr
	}
	// fmt.Printf("after update experiences\n")

	// deleting experiences
	expDelErr := repo.expRepo.DeleteTxExperiencesByIDs(tx, experiencesIDsToDelete)
	if expDelErr != nil {
		tx.Rollback()
		return expDelErr
	}
	// fmt.Printf("after deleting experiences\n")

	// inserting experiences
	expInsErr := repo.expRepo.AddTxExperiences(tx, cvID, experiencesToInsert)
	if expInsErr != nil {
		tx.Rollback()
		return expInsErr
	}
	// fmt.Printf("after inserting experiences\n")

	// updating institutions
	instUpdErr := repo.instRepo.UpdateTxInstitutions(tx, cvID, insitutionsToUpdate)
	if instUpdErr != nil {
		tx.Rollback()
		return instUpdErr
	}
	// fmt.Printf("after updating institutions\n")

	// deleting institutions
	instDelErr := repo.instRepo.DeleteTxExperiencesByIDs(tx, insitutionsIDsToDelete)
	if instDelErr != nil {
		tx.Rollback()
		return instDelErr
	}
	// fmt.Printf("after deleting institutions\n")

	// inserting institutions
	instInsErr := repo.instRepo.AddTxInstitutions(tx, cvID, insitutionsToInsert)
	if instInsErr != nil {
		tx.Rollback()
		return instInsErr
	}
	// fmt.Printf("after inserting institutions\n")

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}

func (repo *psqlCVRepository) DeleteOneOfUsersCV(userID, cvID int) error {
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return txErr
	}
	query := `DELETE
	FROM
		hnh_data.cv c
			USING hnh_data.applicant a
	WHERE
		c.id = $1
		AND a.user_id = $2
		AND c.applicant_id = a.id`

	result, err := tx.Exec(query, cvID, userID)
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

	expErr := repo.expRepo.DeleteTxExperiences(tx, cvID)
	if expErr != nil {
		tx.Rollback()
		return expErr
	}

	instErr := repo.instRepo.DeleteTxExperiences(tx, cvID)
	if instErr != nil {
		tx.Rollback()
		return instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}
