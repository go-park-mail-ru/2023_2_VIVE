package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	"HnH/pkg/queryUtils"
	"context"
	"database/sql"
	"strings"

	"github.com/sirupsen/logrus"
)

type ICVRepository interface {
	GetCVById(ctx context.Context, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	GetCVsByIds(ctx context.Context, idList []int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	GetCVsByUserId(ctx context.Context, userID int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	GetApplicantInfo(ctx context.Context, applicantID int) (string, string, string, []domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) //
	AddCV(ctx context.Context, userID int, cv *domain.DbCV, experiences []domain.DbExperience, insitutions []domain.DbEducationInstitution) (int, error)
	GetOneOfUsersCV(ctx context.Context, userID, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error)
	UpdateOneOfUsersCV(ctx context.Context, userID, cvID int, cv *domain.DbCV,
		experiencesIDsToDelete []int, experiencesToUpdate, experiencesToInsert []domain.DbExperience,
		insitutionsIDsToDelete []int, insitutionsToUpdate, insitutionsToInsert []domain.DbEducationInstitution,
	) error
	DeleteOneOfUsersCV(ctx context.Context, userID, cvID int) error
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

func (repo *psqlCVRepository) GetCVById(ctx context.Context, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting cv by 'cv_id'")

	contextLogger.Debug("starting transaction")
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, ErrEntityNotFound
	}
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, err
	}

	expToReturn, expErr := repo.expRepo.GetTxExperiences(ctx, tx, cvID)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, expErr
	}

	edInstToReturn, instErr := repo.instRepo.GetTxInstitutions(ctx, tx, cvID)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}
	contextLogger.Debug("commiting transaction")

	return &cvToReturn, expToReturn, edInstToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByIds(ctx context.Context, idList []int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting cv list by 'cv_id' list")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		return nil, nil, nil, txErr
	}
	if len(idList) == 0 {
		commitErr := tx.Commit()
		if commitErr != nil {
			return nil, nil, nil, commitErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
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
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, nil, nil, rollbackErr
			}

			return nil, nil, nil, err
		}
		cvIDs = append(cvIDs, cv.ID)
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}

		return nil, nil, nil, ErrEntityNotFound
	}

	expsToReturn, expErr := repo.expRepo.GetTxExperiencesByIds(ctx, tx, cvIDs)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}

		return nil, nil, nil, expErr
	}
	instsToReturn, instErr := repo.instRepo.GetTxInstitutionsByIds(ctx, tx, cvIDs)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}

		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return cvsToReturn, expsToReturn, instsToReturn, nil
}

func (repo *psqlCVRepository) GetCVsByUserId(ctx context.Context, userID int) ([]domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting user's cv list")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		commitErr := tx.Commit()
		if commitErr != nil {
			return nil, nil, nil, commitErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
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
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, nil, nil, rollbackErr
			}
			return nil, nil, nil, err
		}
		cvIDs = append(cvIDs, cv.ID)
		cvsToReturn = append(cvsToReturn, cv)
	}
	if len(cvsToReturn) == 0 {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, ErrEntityNotFound
	}

	// fmt.Printf("after cv select\n")
	// fmt.Printf("cvIDs: %v\n", cvIDs)

	expsToReturn, expErr := repo.expRepo.GetTxExperiencesByIds(ctx, tx, cvIDs)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, expErr
	}
	// fmt.Printf("after exp select\n")

	instsToReturn, instErr := repo.instRepo.GetTxInstitutionsByIds(ctx, tx, cvIDs)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, instErr
	}
	// fmt.Printf("after inst select\n")

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return cvsToReturn, expsToReturn, instsToReturn, nil
}

func (repo *psqlCVRepository) isApplicant(logger *logrus.Entry, applicantID int) (bool, error) {
	logger.Info("checking applicant for given 'user_id'")

	var isApplicant bool

	appErr := repo.DB.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.applicant WHERE id = $1)`, applicantID).Scan(&isApplicant)
	if appErr != nil {
		return false, appErr
	}

	return isApplicant, nil
}

func (repo *psqlCVRepository) GetApplicantInfo(ctx context.Context, applicantID int) (string, string, string, []domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"applicant_id": applicantID,
	}).
		Info("getting cvs and personal info by 'applicant_id' from postgres")

	isApp, err := repo.isApplicant(contextLogger, applicantID)
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}

	if !isApp {
		return "", "", "", nil, nil, nil, ErrEntityNotFound
	}

	var userID int

	err = repo.DB.QueryRow(`SELECT user_id FROM hnh_data.applicant WHERE id = $1`, applicantID).Scan(&userID)
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}

	var first_name, last_name, email string

	err = repo.DB.QueryRow(`SELECT first_name, last_name, email FROM hnh_data.user_profile WHERE id = $1`, userID).Scan(&first_name, &last_name, &email)
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}

	email = strings.TrimSpace(email)

	cvs, exp, edu, err := repo.GetCVsByUserId(ctx, userID)
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}

	return first_name, last_name, email, cvs, exp, edu, nil
}

func (repo *psqlCVRepository) AddCV(ctx context.Context,
	userID int, cv *domain.DbCV,
	experiences []domain.DbExperience,
	insitutions []domain.DbEducationInstitution,
) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("adding new cv")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, ErrNotInserted
	}
	if insertCvErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, insertCvErr
	}
	// fmt.Println("after cv insert")
	expErr := repo.expRepo.AddTxExperiences(ctx, tx, insertedCVID, experiences)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, expErr
	}
	// fmt.Println("after exp insert")
	instErr := repo.instRepo.AddTxInstitutions(ctx, tx, insertedCVID, insitutions)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, instErr
	}
	// fmt.Println("after inst insert")

	commitErr := tx.Commit()
	if commitErr != nil {
		return 0, commitErr
	}

	return insertedCVID, nil
}

func (repo *psqlCVRepository) GetOneOfUsersCV(ctx context.Context, userID, cvID int) (*domain.DbCV, []domain.DbExperience, []domain.DbEducationInstitution, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting one of the user's cv by 'user_id' and 'cv_id'")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, ErrEntityNotFound
	}
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, err
	}

	exps, expErr := repo.expRepo.GetTxExperiences(ctx, tx, cvID)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, expErr
	}

	insts, instErr := repo.instRepo.GetTxInstitutions(ctx, tx, cvID)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, nil, nil, rollbackErr
		}
		return nil, nil, nil, instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, nil, nil, commitErr
	}

	return &cvToReturn, exps, insts, nil
}

func (repo *psqlCVRepository) UpdateOneOfUsersCV(
	ctx context.Context, userID, cvID int, cv *domain.DbCV,
	experiencesIDsToDelete []int, experiencesToUpdate, experiencesToInsert []domain.DbExperience,
	insitutionsIDsToDelete []int, insitutionsToUpdate, insitutionsToInsert []domain.DbEducationInstitution,
) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("updating one of the user's cv")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return ErrNoRowsUpdated
	}
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	// fmt.Printf("after update cv\n")

	// updating experiences
	expUpdErr := repo.expRepo.UpdateTxExperiences(ctx, tx, cvID, experiencesToUpdate)
	if expUpdErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return expUpdErr
	}
	// fmt.Printf("after update experiences\n")

	// deleting experiences
	expDelErr := repo.expRepo.DeleteTxExperiencesByIDs(ctx, tx, experiencesIDsToDelete)
	if expDelErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return expDelErr
	}
	// fmt.Printf("after deleting experiences\n")

	// inserting experiences
	expInsErr := repo.expRepo.AddTxExperiences(ctx, tx, cvID, experiencesToInsert)
	if expInsErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return expInsErr
	}
	// fmt.Printf("after inserting experiences\n")

	// updating institutions
	instUpdErr := repo.instRepo.UpdateTxInstitutions(ctx, tx, cvID, insitutionsToUpdate)
	if instUpdErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return instUpdErr
	}
	// fmt.Printf("after updating institutions\n")

	// deleting institutions
	instDelErr := repo.instRepo.DeleteTxInstitutionsByIDs(ctx, tx, insitutionsIDsToDelete)
	if instDelErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return instDelErr
	}
	// fmt.Printf("after deleting institutions\n")

	// inserting institutions
	instInsErr := repo.instRepo.AddTxInstitutions(ctx, tx, cvID, insitutionsToInsert)
	if instInsErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return instInsErr
	}
	// fmt.Printf("after inserting institutions\n")

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}

func (repo *psqlCVRepository) DeleteOneOfUsersCV(ctx context.Context, userID, cvID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("deleting one of the user's cv")
	tx, txErr := repo.DB.Begin()
	if txErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
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
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return ErrNoRowsDeleted
	}
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	expErr := repo.expRepo.DeleteTxExperiences(ctx, tx, cvID)
	if expErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return expErr
	}

	instErr := repo.instRepo.DeleteTxInstitutions(ctx, tx, cvID)
	if instErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return instErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	return nil
}
