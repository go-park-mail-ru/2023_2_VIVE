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
	GetOrgId(vacancyID int) (int, error)
	AddVacancy(vacancy *domain.Vacancy) (int, error)
	UpdateVacancy(vacancy *domain.Vacancy) (int64, error)
	DeleteVacancy(vacancyID int) (int64, error)
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
		return nil, ENTITY_NOT_FOUND
	}
	if err != nil {
		return nil, err
	}

	return &vacancyToReturn, nil
}

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
		return 0, ENTITY_NOT_FOUND
	}
	if err != nil {
		return 0, err
	}

	return organizationIDToReturn, nil
}

// Add new vacancy and return new id if successful
func (repo *psqlVacancyRepository) AddVacancy(vacancy *domain.Vacancy) (int, error) {
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
	VALUES ($1, $1, $1, $1, $1, $1, $1, $1, $1, $1)
	RETURNING id`

	var insertedVacancyID int
	err := repo.DB.QueryRow(
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
	).
		Scan(&insertedVacancyID)

	if err == sql.ErrNoRows {
		return 0, ENTITY_NOT_FOUND
	}
	if err != nil {
		return 0, err
	}

	return insertedVacancyID, nil
}

func (repo *psqlVacancyRepository) UpdateVacancy(vacancy *domain.Vacancy) (int64, error) {
	query := `UPDATE
		hnh_data.vacancy
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
	WHERE
		id = $11`

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
		vacancy.ID,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (repo *psqlVacancyRepository) DeleteVacancy(vacancyID int) (int64, error) {
	query := `DELETE
		FROM
			hnh_data.vacancy
		WHERE
			id = $1;`

	result, err := repo.DB.Exec(query, vacancyID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
