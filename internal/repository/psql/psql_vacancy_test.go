package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/nullTypes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	ErrQuery = fmt.Errorf("some query error")
)

func sliceIntToDriverValue(slice []int) []driver.Value {
	result := make([]driver.Value, len(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = slice[i]
	}

	return result
}

var vacanciesColumns = []string{
	"id",
	"employer_id",
	"name",
	"description",
	"salary_lower_bound",
	"salary_upper_bound",
	"employment",
	"experience_lower_bound",
	"experience_upper_bound",
	"education_type",
	"location",
	"created_at",
	"updated_at",
}

var location, _ = time.LoadLocation("Local")
var created_at = time.Date(2023, 11, 1, 0, 0, 0, 0, location)
var updated_at = time.Date(2023, 11, 2, 0, 0, 0, 0, location)

var testGetAllVacanciesSuccessCases = []struct {
	expected []domain.Vacancy
}{
	{
		expected: []domain.Vacancy{
			{
				ID:                     1,
				Employer_id:            1,
				VacancyName:            "Vacancy #1",
				Description:            "Description #1",
				Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
				Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
				Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
				Experience_lower_bound: nullTypes.NewNullInt(0, true),
				Experience_upper_bound: nullTypes.NewNullInt(12, true),
				EducationType:          domain.Higher,
				Location:               nullTypes.NewNullString("Moscow", true),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
			{
				ID:                     2,
				Employer_id:            2,
				VacancyName:            "Vacancy #2",
				Description:            "Description #2",
				Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
				Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
				Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
				Experience_lower_bound: nullTypes.NewNullInt(0, true),
				Experience_upper_bound: nullTypes.NewNullInt(12, true),
				EducationType:          domain.Higher,
				Location:               nullTypes.NewNullString("Moscow", true),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
		},
	},
	{
		expected: []domain.Vacancy{
			{
				ID:                     1,
				Employer_id:            1,
				VacancyName:            "Vacancy #1",
				Description:            "Description #1",
				Salary_lower_bound:     nullTypes.NewNullInt(0, false),
				Salary_upper_bound:     nullTypes.NewNullInt(0, false),
				Employment:             nullTypes.NewNullString("", false),
				Experience_lower_bound: nullTypes.NewNullInt(0, false),
				Experience_upper_bound: nullTypes.NewNullInt(0, false),
				EducationType:          domain.Secondary,
				Location:               nullTypes.NewNullString("", false),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
		},
	},
}

func TestGetAllVacanciesSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetAllVacanciesSuccessCases {
		rows := sqlmock.NewRows(vacanciesColumns)

		for _, item := range testCase.expected {
			rows = rows.AddRow(
				item.ID,
				item.Employer_id,
				item.VacancyName,
				item.Description,
				item.Salary_lower_bound,
				item.Salary_upper_bound,
				item.Employment,
				item.Experience_lower_bound,
				item.Experience_upper_bound,
				item.EducationType,
				item.Location,
				item.Created_at,
				item.Updated_at,
			)
		}
		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+").
			WillReturnRows(rows)

		actual, err := repo.GetAllVacancies()
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("results not match, want %v, have %v", testCase.expected, actual)
			return
		}
	}
}

func TestGetAllVacanciesQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	mock.
		ExpectQuery("SELECT(.|\n)+FROM(.|\n)+").
		WillReturnError(ErrQuery)

	_, returnedErr := repo.GetAllVacancies()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr == nil {
		t.Errorf("expected query error, got: '%s'", returnedErr)
		return
	}
}

func TestGetAllVacanciesEntityNotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	rows := sqlmock.NewRows(vacanciesColumns)

	mock.
		ExpectQuery("SELECT(.|\n)+FROM(.|\n)+").
		WillReturnRows(rows)

	_, returnedErr := repo.GetAllVacancies()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != ErrEntityNotFound {
		t.Errorf("expected error 'ErrEntityNotFound', got: '%s'", returnedErr)
		return
	}
}

var testGetVacanciesByIdsSuccessCases = []struct {
	input    []int
	expected []domain.Vacancy
}{
	{
		input: []int{1, 2},
		expected: []domain.Vacancy{
			{
				ID:                     1,
				Employer_id:            1,
				VacancyName:            "Vacancy #1",
				Description:            "Description #1",
				Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
				Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
				Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
				Experience_lower_bound: nullTypes.NewNullInt(0, true),
				Experience_upper_bound: nullTypes.NewNullInt(12, true),
				EducationType:          domain.Higher,
				Location:               nullTypes.NewNullString("Moscow", true),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
			{
				ID:                     2,
				Employer_id:            2,
				VacancyName:            "Vacancy #2",
				Description:            "Description #2",
				Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
				Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
				Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
				Experience_lower_bound: nullTypes.NewNullInt(0, true),
				Experience_upper_bound: nullTypes.NewNullInt(12, true),
				EducationType:          domain.Higher,
				Location:               nullTypes.NewNullString("Moscow", true),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
		},
	},
	{
		input: []int{1},
		expected: []domain.Vacancy{
			{
				ID:                     1,
				Employer_id:            1,
				VacancyName:            "Vacancy #1",
				Description:            "Description #1",
				Salary_lower_bound:     nullTypes.NewNullInt(0, false),
				Salary_upper_bound:     nullTypes.NewNullInt(0, false),
				Employment:             nullTypes.NewNullString("", false),
				Experience_lower_bound: nullTypes.NewNullInt(0, false),
				Experience_upper_bound: nullTypes.NewNullInt(0, false),
				EducationType:          domain.Secondary,
				Location:               nullTypes.NewNullString("", false),
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
		},
	},
}

func TestGetVacanciesByIdsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacanciesByIdsSuccessCases {
		rows := sqlmock.NewRows(vacanciesColumns)

		for _, item := range testCase.expected {
			rows = rows.AddRow(
				item.ID,
				item.Employer_id,
				item.VacancyName,
				item.Description,
				item.Salary_lower_bound,
				item.Salary_upper_bound,
				item.Employment,
				item.Experience_lower_bound,
				item.Experience_upper_bound,
				item.EducationType,
				item.Location,
				item.Created_at,
				item.Updated_at,
			)
		}
		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
			WithArgs(sliceIntToDriverValue(testCase.input)...).
			WillReturnRows(rows)

		actual, err := repo.GetVacanciesByIds(testCase.input)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("results not match, want %v, have %v", testCase.expected, actual)
			return
		}
	}
}

func TestGetVacanciesByIdsEmptyInput(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	_, err = repo.GetVacanciesByIds([]int{})
	if err != ErrEntityNotFound {
		t.Errorf("expected error 'ErrEntityNotFound', got %s", err)
	}
}

func TestGetVacanciesByIdsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	input := []int{1, 2, 3}
	mock.
		ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
		WithArgs(sliceIntToDriverValue(input)...).
		WillReturnError(ErrQuery)

	_, returnedErr := repo.GetVacanciesByIds(input)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != ErrQuery {
		t.Errorf("expected query error, got: '%s'", returnedErr)
		return
	}
}

func TestGetVacanciesByIdsEntityNotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	input := []int{1, 2, 3}
	rows := sqlmock.NewRows(vacanciesColumns)
	mock.
		ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
		WithArgs(sliceIntToDriverValue(input)...).
		WillReturnRows(rows)

	_, returnedErr := repo.GetVacanciesByIds(input)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != ErrEntityNotFound {
		t.Errorf("expected 'ErrEntityNotFound', got: '%s'", returnedErr)
		return
	}
}

var testGetVacancySuccessCases = []struct {
	input    int
	expected domain.Vacancy
}{
	{
		input: 1,
		expected: domain.Vacancy{
			ID:                     1,
			Employer_id:            1,
			VacancyName:            "Vacancy #1",
			Description:            "Description #1",
			Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
			Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
			Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
			Experience_lower_bound: nullTypes.NewNullInt(0, true),
			Experience_upper_bound: nullTypes.NewNullInt(12, true),
			EducationType:          domain.Higher,
			Location:               nullTypes.NewNullString("Moscow", true),
			Created_at:             created_at,
			Updated_at:             updated_at,
		},
	},
	{
		input: 2,
		expected: domain.Vacancy{
			ID:                     2,
			Employer_id:            2,
			VacancyName:            "Vacancy #2",
			Description:            "Description #2",
			Salary_lower_bound:     nullTypes.NewNullInt(10000, true),
			Salary_upper_bound:     nullTypes.NewNullInt(20000, true),
			Employment:             nullTypes.NewNullString(string(domain.FullTime), true),
			Experience_lower_bound: nullTypes.NewNullInt(0, true),
			Experience_upper_bound: nullTypes.NewNullInt(12, true),
			EducationType:          domain.Higher,
			Location:               nullTypes.NewNullString("Moscow", true),
			Created_at:             created_at,
			Updated_at:             updated_at,
		},
	},
	{
		input: 1,
		expected: domain.Vacancy{
			ID:                     1,
			Employer_id:            1,
			VacancyName:            "Vacancy #1",
			Description:            "Description #1",
			Salary_lower_bound:     nullTypes.NewNullInt(0, false),
			Salary_upper_bound:     nullTypes.NewNullInt(0, false),
			Employment:             nullTypes.NewNullString("", false),
			Experience_lower_bound: nullTypes.NewNullInt(0, false),
			Experience_upper_bound: nullTypes.NewNullInt(0, false),
			EducationType:          domain.Secondary,
			Location:               nullTypes.NewNullString("", false),
			Created_at:             created_at,
			Updated_at:             updated_at,
		},
	},
}

func TestGetVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacancySuccessCases {
		rows := sqlmock.NewRows(vacanciesColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.Employer_id,
				testCase.expected.VacancyName,
				testCase.expected.Description,
				testCase.expected.Salary_lower_bound,
				testCase.expected.Salary_upper_bound,
				testCase.expected.Employment,
				testCase.expected.Experience_lower_bound,
				testCase.expected.Experience_upper_bound,
				testCase.expected.EducationType,
				testCase.expected.Location,
				testCase.expected.Created_at,
				testCase.expected.Updated_at,
			)

		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
			WithArgs(testCase.input).
			WillReturnRows(rows)

		actual, err := repo.GetVacancy(testCase.input)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(*actual, testCase.expected) {
			t.Errorf("results not match, want %v, have %v", testCase.expected, actual)
			return
		}
	}
}

var testGetVacancyQueryErrorCases = []struct {
	input        int
	returningErr error
	expectedErr  error
}{
	{
		input:        1,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrEntityNotFound,
	},
	{
		input:        1,
		returningErr: ErrQuery,
		expectedErr:  ErrQuery,
	},
}

func TestGetVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacancyQueryErrorCases {
		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
			WithArgs(testCase.input).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetVacancy(testCase.input)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedErr {
			t.Errorf("expected query error: '%s'\ngot: '%s'", testCase.expectedErr, actualErr)
			return
		}
	}
}

var testGetOrgIdSuccessCases = []struct {
	input    int
	expected int
}{
	{
		input:    1,
		expected: 1,
	},
	{
		input:    2,
		expected: 2,
	},
	{
		input:    3,
		expected: 3,
	},
}

func TestGetOrgIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetOrgIdSuccessCases {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(
				testCase.expected,
			)

		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
			WithArgs(testCase.input).
			WillReturnRows(rows)

		actual, err := repo.GetOrgId(testCase.input)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("results not match, want %v, have %v", testCase.expected, actual)
			return
		}
	}
}

var testGetOrgIdQueryErrorCases = []struct {
	input        int
	returningErr error
	expectedErr  error
}{
	{
		input:        1,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrEntityNotFound,
	},
	{
		input:        1,
		returningErr: ErrQuery,
		expectedErr:  ErrQuery,
	},
}

func TestGetOrgIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testGetOrgIdQueryErrorCases {
		mock.
			ExpectQuery("SELECT(.|\n)+FROM(.|\n)+WHERE(.|\n)+").
			WithArgs(testCase.input).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetOrgId(testCase.input)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedErr {
			t.Errorf("expected query error: '%s'\ngot: '%s'", testCase.expectedErr, actualErr)
			return
		}
	}
}
