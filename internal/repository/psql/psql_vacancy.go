package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/queryUtils"
	"database/sql"
	"errors"
)

type IVacancyRepository interface {
	GetAllVacancies() ([]domain.DbVacancy, error)
	GetVacanciesByIds(orgID int, idList []int) ([]domain.DbVacancy, error)
	GetVacancy(vacancyID int) (*domain.DbVacancy, error)
	GetUserVacancies(userID int) ([]domain.DbVacancy, error)
	// GetVacancyByUserID(userID int, vacancyID int) (*domain.Vacancy, error)
	GetOrgId(vacancyID int) (int, error)
	AddVacancy(userID int, vacancy *domain.DbVacancy) (int, error)
	UpdateOrgVacancy(orgID, vacancyID int, vacancy *domain.DbVacancy) error
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

func (repo *psqlVacancyRepository) GetAllVacancies() ([]domain.DbVacancy, error) {
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

func (repo *psqlVacancyRepository) GetVacanciesByIds(orgID int, idList []int) ([]domain.DbVacancy, error) {
	if len(idList) == 0 {
		return nil, ErrEntityNotFound
	}

	items := make([]int, len(idList)+1)
	items[0] = orgID
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
	JOIN hnh_data.employer e ON
		v.employer_id = e.id
	WHERE
		e.organization_id = $1
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

func (repo *psqlVacancyRepository) GetVacancy(vacancyID int) (*domain.DbVacancy, error) {
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

func (repo *psqlVacancyRepository) GetUserVacancies(userID int) ([]domain.DbVacancy, error) {
	var empID int

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
func (repo *psqlVacancyRepository) AddVacancy(userID int, vacancy *domain.DbVacancy) (int, error) {
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
	// fmt.Printf("*vacancy.Salary_lower_bound: %v\n", *vacancy.Salary_lower_bound)
	// fmt.Printf("*vacancy.Salary_upper_bound: %v\n", *vacancy.Salary_upper_bound)
	// fmt.Printf("vacancy.Employment: %v\n", vacancy.Employment)
	// fmt.Printf("*vacancy.Experience_lower_bound: %v\n", *vacancy.Experience_lower_bound)
	// fmt.Printf("*vacancy.Experience_upper_bound: %v\n", *vacancy.Experience_upper_bound)
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
		return 0, ErrNotInserted
	}
	if err != nil {
		return 0, err
	}

	return insertedVacancyID, nil
}

func (repo *psqlVacancyRepository) UpdateOrgVacancy(orgID, vacancyID int, vacancy *domain.DbVacancy) error {
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
		AND e.organization_id = $11
		AND v.employer_id = e.id`

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
		orgID,
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
