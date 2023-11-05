package repository

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"

	// "fmt"
	// "strings"

	// "context"
	// "HnH/internal/repository/mock"
	"database/sql"
	// "github.com/jackc/pgx/stdlib"
	// _ "github.com/jackc/pgx/stdlib"
)

type IVacancyRepository interface {
	GetAllVacancies() ([]domain.Vacancy, error)
	GetVacanciesByIds(idList []int) ([]domain.Vacancy, error)
	GetVacancy(vacancyID int) (*domain.Vacancy, error)
	// GetVacancyByUserID(userID int, vacancyID int) (*domain.Vacancy, error)
	GetOrgId(vacancyID int) (int, error)
	AddVacancy(userID int, vacancy *domain.Vacancy) (int, error)
	UpdateOrgVacancy(orgID, vacancyID int, vacancy *domain.Vacancy) error
	DeleteOrgVacancy(orgID, vacancyID int) error
}

type psqlVacancyRepository struct {
	DB *sql.DB
}

func NewPsqlVacancyRepository(db *sql.DB) IVacancyRepository {
	return &psqlVacancyRepository{
		DB: db,
	}
}

func (repo *psqlVacancyRepository) GetAllVacancies() ([]domain.Vacancy, error) {
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

	vacanciesToReturn := []domain.Vacancy{}

	for rows.Next() {
		vacancy := domain.Vacancy{}
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.Employer_id,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.Salary_lower_bound,
			&vacancy.Salary_upper_bound,
			&vacancy.Employment,
			&vacancy.Experience_lower_bound,
			&vacancy.Experience_upper_bound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.Created_at,
			&vacancy.Updated_at,
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

func (repo *psqlVacancyRepository) GetVacanciesByIds(idList []int) ([]domain.Vacancy, error) {
	placeHolderValues := *queryUtils.IntToAnySlice(idList)
	placeHolderString := queryUtils.QueryPlaceHolders(placeHolderValues...)

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
		v.id IN (` + placeHolderString + `)`

	rows, err := repo.DB.Query(query, placeHolderValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vacanciesToReturn := []domain.Vacancy{}

	for rows.Next() {
		vacancy := domain.Vacancy{}
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.Employer_id,
			&vacancy.VacancyName,
			&vacancy.Description,
			&vacancy.Salary_lower_bound,
			&vacancy.Salary_upper_bound,
			&vacancy.Employment,
			&vacancy.Experience_lower_bound,
			&vacancy.Experience_upper_bound,
			&vacancy.EducationType,
			&vacancy.Location,
			&vacancy.Created_at,
			&vacancy.Updated_at,
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

func (repo *psqlVacancyRepository) GetVacancy(vacancyID int) (*domain.Vacancy, error) {
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

	vacancyToReturn := domain.Vacancy{}

	err := repo.DB.QueryRow(query, vacancyID).
		Scan(
			&vacancyToReturn.ID,
			&vacancyToReturn.Employer_id,
			&vacancyToReturn.VacancyName,
			&vacancyToReturn.Description,
			&vacancyToReturn.Salary_lower_bound,
			&vacancyToReturn.Salary_upper_bound,
			&vacancyToReturn.Employment,
			&vacancyToReturn.Experience_lower_bound,
			&vacancyToReturn.Experience_upper_bound,
			&vacancyToReturn.EducationType,
			&vacancyToReturn.Location,
			&vacancyToReturn.Created_at,
			&vacancyToReturn.Updated_at,
		)
	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}
	if err != nil {
		return nil, err
	}

	return &vacancyToReturn, nil
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


func (repo *psqlVacancyRepository) GetOrgId(vacancyID int) (int, error) {
	query := `SELECT
		organization_id
	FROM
	(
	SELECT
			e.organization_id,
			v.id
	FROM
			hnh_data.employer e
	LEFT JOIN hnh_data.vacancy v ON
			e.id = v.employer_id
	WHERE
			v.id IS NOT NULL
	) AS w
	WHERE
		id = $1`

	var organizationIDToReturn int
	err := repo.DB.QueryRow(query, vacancyID).
		Scan(&organizationIDToReturn)

	if err == sql.ErrNoRows {
		return 0, ErrEntityNotFound
	}
	if err != nil {
		return 0, err
	}

	return organizationIDToReturn, nil
}

// Add new vacancy and return new id if successful
func (repo *psqlVacancyRepository) AddVacancy(userID int, vacancy *domain.Vacancy) (int, error) {
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
    e.user_id = $10
    RETURNING id`

	var insertedVacancyID int
	err := repo.DB.QueryRow(
		query,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary_lower_bound,
		vacancy.Salary_upper_bound,
		vacancy.Employment,
		vacancy.Experience_lower_bound,
		vacancy.Experience_upper_bound,
		vacancy.EducationType,
		vacancy.Location,
		userID,
	).
		Scan(&insertedVacancyID)

	if err == sql.ErrNoRows {
		return 0, ErrNotInserted
	}
	if err != nil {
		return 0, err
	}

	return insertedVacancyID, nil
}

func (repo *psqlVacancyRepository) UpdateOrgVacancy(orgID, vacancyID int, vacancy *domain.Vacancy) error {
	query := `UPDATE
		hnh_data.vacancy v
	SET
		employer_id = $1,
		"name" = $2,
		description = $3,
		salary_lower_bound = $4,
		salary_upper_bound = $5,
		employment = $6,
		experience_lower_bound = $7,
		experience_upper_bound = $8,
		education_type = $9,
		"location" = $10,
		updated_at = now()
	FROM
		hnh_data.employer e
	WHERE
		v.id = $11
		AND e.organization_id = $12
		AND v.employer_id = e.id`

	result, err := repo.DB.Exec(
		query,
		vacancy.Employer_id,
		vacancy.VacancyName,
		vacancy.Description,
		vacancy.Salary_lower_bound,
		vacancy.Salary_upper_bound,
		vacancy.Employment,
		vacancy.Experience_lower_bound,
		vacancy.Experience_upper_bound,
		vacancy.EducationType,
		vacancy.Location,
		vacancyID,
		orgID,
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

func (repo *psqlVacancyRepository) DeleteOrgVacancy(orgID, vacancyID int) error {
	query := `DELETE
	FROM
		hnh_data.vacancy v
		USING hnh_data.employer e
	WHERE
		v.id = $1
		AND e.organization_id = $2
		AND v.employer_id = e.id`

	result, err := repo.DB.Exec(query, vacancyID, orgID)
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
