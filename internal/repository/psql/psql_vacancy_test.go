package psql

import (
	"HnH/internal/domain"
	"reflect"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var location, _ = time.LoadLocation("Local")
var created_at = time.Date(2023, 11, 1, 0, 0, 0, 0, location)
var updated_at = time.Date(2023, 11, 2, 0, 0, 0, 0, location)

var rows = []struct {
	colomns  []string
	expected []domain.Vacancy
}{
	{
		colomns: []string{
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
		},
		expected: []domain.Vacancy{
			{
				ID:                     1,
				Employer_id:            1,
				VacancyName:            "Vacancy #1",
				Description:            "Description #1",
				Salary_lower_bound:     10000,
				Salary_upper_bound:     20000,
				Employment:             domain.FullTime,
				Experience_lower_bound: 0,
				Experience_upper_bound: 12,
				EducationType:          domain.Higher,
				Location:               "Moscow",
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
			{
				ID:                     2,
				Employer_id:            2,
				VacancyName:            "Vacancy #2",
				Description:            "Description #2",
				Salary_lower_bound:     10000,
				Salary_upper_bound:     20000,
				Employment:             domain.FullTime,
				Experience_lower_bound: 0,
				Experience_upper_bound: 12,
				EducationType:          domain.Higher,
				Location:               "Moscow",
				Created_at:             created_at,
				Updated_at:             updated_at,
			},
		},
	},
}

func TestGetAllVacancies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range rows {
		rows := sqlmock.NewRows(testCase.colomns)

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
