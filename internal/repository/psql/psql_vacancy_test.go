package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/nullTypes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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
		WillReturnError(fmt.Errorf("some query error"))

	_, queryErr := repo.GetAllVacancies()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if queryErr == nil {
		t.Errorf("expected query error, got nil")
		return
	}
}

func TestGetAllVacanciesEnriryNotFoundError(t *testing.T) {
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
	if returnedErr == nil {
		t.Errorf("expected error 'ErrEntityNotFound', got nil")
		return
	}
}
