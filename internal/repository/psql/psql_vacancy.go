package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	"HnH/pkg/queryUtils"
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

type IVacancyRepository interface {
	GetAllVacancies(ctx context.Context) ([]domain.DbVacancy, error)
	GetEmpVacanciesByIds(ctx context.Context, empID int, idList []int) ([]domain.DbVacancy, error)
	GetVacanciesByIds(ctx context.Context, idList []int) ([]domain.DbVacancy, error)
	GetVacancy(ctx context.Context, vacancyID int) (*domain.DbVacancy, error)
	GetUserVacancies(ctx context.Context, userID int) ([]domain.DbVacancy, error)
	// GetVacancyByUserID(userID int, vacancyID int) (*domain.Vacancy, error)
	GetEmpId(ctx context.Context, vacancyID int) (int, error)
	AddVacancy(ctx context.Context, userID int, vacancy *domain.DbVacancy) (int, error)
	UpdateEmpVacancy(ctx context.Context, empID, vacancyID int, vacancy *domain.DbVacancy) error
	DeleteEmpVacancy(ctx context.Context, empID, vacancyID int) error
}

type psqlVacancyRepository struct {
	DB *sql.DB
}

func NewPsqlVacancyRepository(db *sql.DB) IVacancyRepository {
	return &psqlVacancyRepository{
		DB: db,
	}
}

func (repo *psqlVacancyRepository) GetAllVacancies(ctx context.Context) ([]domain.DbVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting all vacancies from postgres")
	query := `SELECT
	    id,
	    employer_id,
	    "name",
	    description,
	    salary_lower_bound,
	    salary_upper_bound,
	    employment,
	    experience_lower_bound,
	    experience_upper_bound,
	    education_type,
	    "location",
	    created_at,
	    updated_at
	FROM
	    hnh_data.vacancy v`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacanciesToReturn := []domain.DbVacancy{}

	for rows.Next() {
		vacancy := domain.DbVacancy{}
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.EmployerID,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.SalaryLowerBound,
			&vacancy.SalaryUpperBound,
			&vacancy.Employment,
			&vacancy.ExperienceLowerBound,
			&vacancy.ExperienceUpperBound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.CreatedAt,
			&vacancy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vacanciesToReturn = append(vacanciesToReturn, vacancy)
	}
	if len(vacanciesToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return vacanciesToReturn, nil
}

func (repo *psqlVacancyRepository) GetEmpVacanciesByIds(ctx context.Context, empID int, idList []int) ([]domain.DbVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"id_list": idList,
		"emp_id":  empID,
	}).
		Info("getting vacancies by 'id_list' and 'emp_id' from postgres")
	if len(idList) == 0 {
		return nil, ErrEntityNotFound
	}

	items := make([]int, len(idList)+1)
	items[0] = empID
	for i := 1; i < len(items); i++ {
		items[i] = idList[i-1]
	}

	placeHolderValues := *queryUtils.IntToAnySlice(items)
	placeHolderString := queryUtils.QueryPlaceHolders(2, len(placeHolderValues)-1)

	query := `SELECT
		v.id,
		v.employer_id,
		v."name",
		v.description,
		v.salary_lower_bound,
		v.salary_upper_bound,
		v.employment,
		v.experience_lower_bound,
		v.experience_upper_bound,
		v.education_type,
		v."location",
		v.created_at,
		v.updated_at
	FROM
		hnh_data.vacancy v
	WHERE
		v.employer_id = $1
		AND v.id IN (` + placeHolderString + `)`

	rows, err := repo.DB.Query(query, placeHolderValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacanciesToReturn := []domain.DbVacancy{}

	for rows.Next() {
		vacancy := domain.DbVacancy{}
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.EmployerID,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.SalaryLowerBound,
			&vacancy.SalaryUpperBound,
			&vacancy.Employment,
			&vacancy.ExperienceLowerBound,
			&vacancy.ExperienceUpperBound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.CreatedAt,
			&vacancy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vacanciesToReturn = append(vacanciesToReturn, vacancy)
	}
	if len(vacanciesToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return vacanciesToReturn, nil
}

func (repo *psqlVacancyRepository) GetVacanciesByIds(ctx context.Context, idList []int) ([]domain.DbVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"id_list": idList,
	}).
		Info("getting vacancies by 'id_list' from postgres")
	if len(idList) == 0 {
		return nil, ErrEntityNotFound
	}

	placeHolderValues := *queryUtils.IntToAnySlice(idList)
	placeHolderString := queryUtils.QueryPlaceHolders(1, len(placeHolderValues))

	query := `SELECT
		v.id,
		v.employer_id,
		v."name",
		v.description,
		v.salary_lower_bound,
		v.salary_upper_bound,
		v.employment,
		v.experience_lower_bound,
		v.experience_upper_bound,
		v.education_type,
		v."location",
		v.created_at,
		v.updated_at
	FROM
		hnh_data.vacancy v
	WHERE
		v.id IN (` + placeHolderString + `)`

	contextLogger.WithFields(logrus.Fields{
		"db_query": query,
	}).
		Debug("query to db")

	rows, err := repo.DB.Query(query, placeHolderValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacanciesToReturn := []domain.DbVacancy{}

	for rows.Next() {
		vacancy := domain.DbVacancy{}
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.EmployerID,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.SalaryLowerBound,
			&vacancy.SalaryUpperBound,
			&vacancy.Employment,
			&vacancy.ExperienceLowerBound,
			&vacancy.ExperienceUpperBound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.CreatedAt,
			&vacancy.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vacanciesToReturn = append(vacanciesToReturn, vacancy)
	}
	if len(vacanciesToReturn) == 0 {
		return nil, ErrEntityNotFound
	}
	return vacanciesToReturn, nil
}

func (repo *psqlVacancyRepository) GetVacancy(ctx context.Context, vacancyID int) (*domain.DbVacancy, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("getting vacanciy by 'vacancy_id' from postgres")
	query := `SELECT
		id,
		employer_id,
		"name",
		description,
		salary_lower_bound,
		salary_upper_bound,
		employment,
		experience_lower_bound,
		experience_upper_bound,
		education_type,
		"location",
		created_at,
		updated_at
	FROM
		hnh_data.vacancy v
	WHERE
		v.id = $1`

	vacancyToReturn := domain.DbVacancy{}

	err := repo.DB.QueryRow(query, vacancyID).
		Scan(
			&vacancyToReturn.ID,
			&vacancyToReturn.EmployerID,
			&vacancyToReturn.VacancyName,
			&vacancyToReturn.Description,
			&vacancyToReturn.SalaryLowerBound,
			&vacancyToReturn.SalaryUpperBound,
			&vacancyToReturn.Employment,
			&vacancyToReturn.ExperienceLowerBound,
			&vacancyToReturn.ExperienceUpperBound,
			&vacancyToReturn.EducationType,
			&vacancyToReturn.Location,
			&vacancyToReturn.CreatedAt,
			&vacancyToReturn.UpdatedAt,
		)
	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &vacancyToReturn, nil
}

func (repo *psqlVacancyRepository) GetUserVacancies(ctx context.Context, userID int) ([]domain.DbVacancy, error) {
	var empID int
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"id_list": userID,
	}).
		Info("getting vacancies by 'user_id' from postgres")

	empErr := repo.DB.QueryRow(`SELECT id FROM hnh_data.employer WHERE user_id = $1`, userID).Scan(&empID)
	if empErr != nil {
		return nil, empErr
	}

	rows, err := repo.DB.Query(`SELECT
		id,
		employer_id,
		"name",
		description,
		salary_lower_bound,
		salary_upper_bound,
		employment,
		experience_lower_bound,
		experience_upper_bound,
		education_type,
		"location",
		created_at,
		updated_at
	FROM
		hnh_data.vacancy v
	WHERE
		v.employer_id = $1`, empID)

	if errors.Is(err, sql.ErrNoRows) {
		return []domain.DbVacancy{}, nil
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	listToReturn := []domain.DbVacancy{}
	for rows.Next() {
		var vacancy domain.DbVacancy

		err := rows.Scan(
			&vacancy.ID,
			&vacancy.EmployerID,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.SalaryLowerBound,
			&vacancy.SalaryUpperBound,
			&vacancy.Employment,
			&vacancy.ExperienceLowerBound,
			&vacancy.ExperienceUpperBound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.CreatedAt,
			&vacancy.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		listToReturn = append(listToReturn, vacancy)
	}

	return listToReturn, nil
}

// func (repo *psqlVacancyRepository) GetVacancyByUserID(userID int, vacancyID int) (*domain.Vacancy, error) {
// 	query := `SELECT
// 		v.id,
// 		v.employer_id,
// 		v."name",
// 		v.description,
// 		v.salary_lower_bound,
// 		v.salary_upper_bound,
// 		v.employment,
// 		v.experience_lower_bound,
// 		v.experience_upper_bound,
// 		v.education_type,
// 		v."location",
// 		v.created_at,
// 		v.updated_at
// 	FROM
// 		hnh_data.vacancy v
// 	LEFT JOIN hnh_data.employer e ON
// 		v.employer_id = e.id
// 	WHERE
// 		e.user_id = $1
// 		AND v.id = $2`

// 	vacancyToReturn := domain.Vacancy{}

// 	err := repo.DB.QueryRow(query, userID, vacancyID).
// 		Scan(
// 			&vacancyToReturn.ID,
// 			&vacancyToReturn.Employer_id,
// 			&vacancyToReturn.VacancyName,
// 			&vacancyToReturn.Description,
// 			&vacancyToReturn.Salary_lower_bound,
// 			&vacancyToReturn.Salary_upper_bound,
// 			&vacancyToReturn.Employment,
// 			&vacancyToReturn.Experience_lower_bound,
// 			&vacancyToReturn.Experience_upper_bound,
// 			&vacancyToReturn.EducationType,
// 			&vacancyToReturn.Location,
// 			&vacancyToReturn.Created_at,
// 			&vacancyToReturn.Updated_at,
// 		)
// 	if err == sql.ErrNoRows {
// 		return nil, ErrEntityNotFound
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &vacancyToReturn, nil
// }

func (repo *psqlVacancyRepository) GetEmpId(ctx context.Context, vacancyID int) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("getting employer id by 'vacancy_id' from postgres")
	query := `SELECT
		v.employer_id 
	FROM
		hnh_data.vacancy v
	WHERE
		v.id = $1`

	var employerIDToReturn int
	err := repo.DB.QueryRow(query, vacancyID).
		Scan(&employerIDToReturn)

	if err == sql.ErrNoRows {
		return 0, ErrEntityNotFound
	}
	if err != nil {
		return 0, err
	}

	return employerIDToReturn, nil
}

// Add new vacancy and return new id if successful
func (repo *psqlVacancyRepository) AddVacancy(ctx context.Context, userID int, vacancy *domain.DbVacancy) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("adding new vacancy  by 'user_id' in postgres")
	// fmt.Printf("before inserting vacancy in db: %v\n", vacancy)
	query := `INSERT
		INTO
		hnh_data.vacancy (
			employer_id,
			"name",
			description,
			salary_lower_bound,
			salary_upper_bound,
			employment,
			experience_lower_bound,
			experience_upper_bound,
			education_type,
			"location"
		)
		SELECT
			e.id, $1, $2, $3, $4, $5, $6, $7, $8, $9
	FROM
		hnh_data.employer e
	WHERE
		e.organization_id = $10
		RETURNING id`

	var insertedVacancyID int

	// fmt.Printf("vacancy.VacancyName: %v\n", vacancy.VacancyName)
	// fmt.Printf("vacancy.Description: %v\n", vacancy.Description)
	// fmt.Printf("*vacancy.Salary_lower_bound: %v\n", *vacancy.SalaryLowerBound)
	// fmt.Printf("*vacancy.Salary_upper_bound: %v\n", *vacancy.SalaryUpperBound)
	// fmt.Printf("vacancy.Employment: %v\n", vacancy.Employment)
	// fmt.Printf("*vacancy.Experience_lower_bound: %v\n", *vacancy.ExperienceLowerBound)
	// fmt.Printf("*vacancy.Experience_upper_bound: %v\n", *vacancy.ExperienceUpperBound)
	// fmt.Printf("vacancy.EducationType: %v\n", vacancy.EducationType)
	// fmt.Printf("*vacancy.Location: %v\n", *vacancy.Location)
	// fmt.Printf("userID: %v\n", userID)

	err := repo.DB.QueryRow(
		query,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.SalaryLowerBound,
		vacancy.SalaryUpperBound,
		vacancy.Employment,
		vacancy.ExperienceLowerBound,
		vacancy.ExperienceUpperBound,
		vacancy.EducationType,
		vacancy.Location,
		userID,
	).
		Scan(&insertedVacancyID)

	if err == sql.ErrNoRows {
		// fmt.Printf("err: %s\n", err)
		return 0, ErrNotInserted
	}
	if err != nil {
		// fmt.Printf("err: %s\n", err)
		return 0, err
	}
	// fmt.Printf("after inserting vacancy in db\n")

	return insertedVacancyID, nil
}

func (repo *psqlVacancyRepository) UpdateEmpVacancy(ctx context.Context, empID, vacancyID int, vacancy *domain.DbVacancy) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id":  vacancyID,
		"employer_id": empID,
	}).
		Info("updating vacancy by 'vacancy_id' and 'employer_id' in postgres")

	query := `UPDATE
		hnh_data.vacancy v
	SET
		"name" = $1,
		description = $2,
		salary_lower_bound = $3,
		salary_upper_bound = $4,
		employment = $5,
		experience_lower_bound = $6,
		experience_upper_bound = $7,
		education_type = $8,
		"location" = $9,
		updated_at = now()
	FROM
		hnh_data.employer e
	WHERE
		v.id = $10
		AND v.employer_id = $11`

	// fmt.Printf("vacancy.VacancyName: %v\n", vacancy.VacancyName)
	// fmt.Printf("vacancy.Description: %v\n", vacancy.Description)
	// fmt.Printf("vacancy.Salary_lower_bound: %v\n", vacancy.Salary_lower_bound)
	// fmt.Printf("vacancy.Salary_upper_bound: %v\n", vacancy.Salary_upper_bound)
	// fmt.Printf("vacancy.Employment: %v\n", vacancy.Employment)
	// fmt.Printf("vacancy.Experience_lower_bound: %v\n", vacancy.Experience_lower_bound)
	// fmt.Printf("vacancy.Experience_upper_bound: %v\n", vacancy.Experience_upper_bound)
	// fmt.Printf("vacancy.EducationType: %v\n", vacancy.EducationType)
	// fmt.Printf("vacancy.Location: %v\n", vacancy.Location)
	// fmt.Printf("vacancyID: %v\n", vacancyID)
	// fmt.Printf("orgID: %v\n", orgID)

	result, err := repo.DB.Exec(
		query,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.SalaryLowerBound,
		vacancy.SalaryUpperBound,
		vacancy.Employment,
		vacancy.ExperienceLowerBound,
		vacancy.ExperienceUpperBound,
		vacancy.EducationType,
		vacancy.Location,
		vacancyID,
		empID,
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

func (repo *psqlVacancyRepository) DeleteEmpVacancy(ctx context.Context, empID, vacancyID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id":  vacancyID,
		"employer_id": empID,
	}).
		Info("deleting vacancy by 'vacancy_id' and 'employer_id' in postgres")

	query := `DELETE
	FROM
		hnh_data.vacancy v
	WHERE
		v.id = $1
		AND v.employer_id = $2`

	result, err := repo.DB.Exec(query, vacancyID, empID)
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
