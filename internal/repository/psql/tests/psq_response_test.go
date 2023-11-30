package psql

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/testHelper"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var testRespondToVacancy = []struct {
	cvID      int
	vacancyID int
}{
	{
		cvID:      1,
		vacancyID: 2,
	},
	{
		cvID:      2,
		vacancyID: 1,
	},
}

func TestRespondToVacancy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlResponseRepository(db)
	for _, testCase := range testRespondToVacancy {
		mock.
			ExpectExec(testHelper.InsertQuery).
			WithArgs(testCase.vacancyID, testCase.cvID).
			WillReturnResult(driver.RowsAffected(1))

		addErr := repo.RespondToVacancy(testHelper.СtxWithLogger, testCase.vacancyID, testCase.cvID)
		if addErr != nil {
			t.Errorf("unexpected err: %v", addErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testRespondToVacancyQueryErrorCases = []struct {
	vacancyID     int
	cvID          int
	queryError    error
	expectedError error
}{
	{
		vacancyID:     1,
		cvID:          1,
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
	{
		vacancyID:     1,
		cvID:          1,
		queryError:    sql.ErrNoRows,
		expectedError: psql.ErrNotInserted,
	},
}

func TestRespondToVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testRespondToVacancyQueryErrorCases {

		repo := psql.NewPsqlResponseRepository(db)
		mock.
			ExpectExec(testHelper.InsertQuery).
			WithArgs(testCase.vacancyID, testCase.cvID).
			WillReturnError(testCase.queryError)

		addErr := repo.RespondToVacancy(testHelper.СtxWithLogger, testCase.vacancyID, testCase.cvID)
		if addErr != testCase.expectedError {
			t.Errorf("unexpected err: %v", addErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testGetVacanciesIdsByCVId = []struct {
	cvID     int
	expected []int
}{
	{
		cvID:     1,
		expected: []int{1, 2, 3},
	},
	{
		cvID:     2,
		expected: []int{1},
	},
}

func TestGetVacanciesIdsByCVId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlResponseRepository(db)
	for _, testCase := range testGetVacanciesIdsByCVId {
		rows := sqlmock.NewRows([]string{"vacancy_id"})
		for _, vacID := range testCase.expected {
			rows.AddRow(vacID)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnRows(rows)

		actual, getErr := repo.GetVacanciesIdsByCVId(testHelper.СtxWithLogger, testCase.cvID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
		}
	}
}

var testGetVacanciesIdsByCVIdQueryErrorCases = []struct {
	cvID          int
	queryError    error
	expectedError error
}{
	{
		cvID:          1,
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
	{
		cvID:          1,
		queryError:    sql.ErrNoRows,
		expectedError: psql.ErrEntityNotFound,
	},
}

func TestGetVacanciesIdsByCVIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testGetVacanciesIdsByCVIdQueryErrorCases {

		repo := psql.NewPsqlResponseRepository(db)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnError(testCase.queryError)

		_, getErr := repo.GetVacanciesIdsByCVId(testHelper.СtxWithLogger, testCase.cvID)
		if getErr != testCase.expectedError {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testGetAttachedCVs = []struct {
	vacancyID int
	expected  []int
}{
	{
		vacancyID: 1,
		expected:  []int{1, 2, 3},
	},
	{
		vacancyID: 2,
		expected:  []int{1},
	},
}

func TestGetAttachedCVs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlResponseRepository(db)
	for _, testCase := range testGetAttachedCVs {
		rows := sqlmock.NewRows([]string{"vacancy_id"})
		for _, vacID := range testCase.expected {
			rows.AddRow(vacID)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnRows(rows)

		actual, getErr := repo.GetAttachedCVs(testHelper.СtxWithLogger, testCase.vacancyID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
		}
	}
}

var testGetAttachedCVsQueryErrorCases = []struct {
	cvID          int
	queryError    error
	expectedError error
}{
	{
		cvID:          1,
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
	{
		cvID:          1,
		queryError:    sql.ErrNoRows,
		expectedError: psql.ErrEntityNotFound,
	},
}

func TestGetAttachedCVsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testGetAttachedCVsQueryErrorCases {

		repo := psql.NewPsqlResponseRepository(db)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnError(testCase.queryError)

		_, getErr := repo.GetAttachedCVs(testHelper.СtxWithLogger, testCase.cvID)
		if getErr != testCase.expectedError {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}
