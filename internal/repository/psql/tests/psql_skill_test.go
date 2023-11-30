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

var testAddSkillsByVacID = []struct {
	vacancyID int
	skills    []string
}{
	{
		vacancyID: 1,
		skills:    []string{"python", "c++"},
	},
	{
		vacancyID: 2,
		skills:    []string{"js", "golang"},
	},
}

func TestAddSkillsByVacID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByVacID {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.vacancyID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByVacIDQueryErrorCases = []struct {
	vacancyID     int
	skills        []string
	queryError    error
	expectedError error
}{
	{
		vacancyID:     1,
		skills:        []string{"python", "c++"},
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
	{
		vacancyID:     1,
		skills:        []string{"python", "c++"},
		queryError:    sql.ErrNoRows,
		expectedError: psql.ErrNotInserted,
	},
}

func TestAddSkillsByVacIDFirstQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByVacIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(testHelper.ErrQuery)

		addErr := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
		if addErr != testHelper.ErrQuery {
			t.Errorf("unexpected err: %v", addErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestAddSkillsByVacIDSecondQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByVacIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(sql.ErrNoRows)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(testCase.queryError)

		addErr := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
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

func TestAddSkillsByVacIDThirdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByVacIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(0),
			)

		mock.
			ExpectExec(testHelper.InsertQuery).
			WithArgs(testCase.vacancyID, 0).
			WillReturnError(testCase.queryError)

		addErr := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
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

var testAddSkillsByVacIDOnExistance = []struct {
	vacancyID int
	skills    []string
}{
	{
		vacancyID: 1,
		skills:    []string{"python", "c++"},
	},
	{
		vacancyID: 2,
		skills:    []string{"js", "golang"},
	},
}

func TestAddSkillsByVacIDOnExistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByVacIDOnExistance {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnError(sql.ErrNoRows)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.vacancyID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByCvID = []struct {
	cvID   int
	skills []string
}{
	{
		cvID:   1,
		skills: []string{"python", "c++"},
	},
	{
		cvID:   2,
		skills: []string{"js", "golang"},
	},
}

func TestAddSkillsByCvID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByCvID {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.cvID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByCvID(ctxWithLogger, testCase.cvID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByCvIDQueryErrorCases = []struct {
	cvID          int
	skills        []string
	queryError    error
	expectedError error
}{
	{
		cvID:          1,
		skills:        []string{"python", "c++"},
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
	{
		cvID:          1,
		skills:        []string{"python", "c++"},
		queryError:    sql.ErrNoRows,
		expectedError: psql.ErrNotInserted,
	},
}

func TestAddSkillsByCvIDFirstQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByCvIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(testHelper.ErrQuery)

		addErr := repo.AddSkillsByCvID(ctxWithLogger, testCase.cvID, testCase.skills)
		if addErr != testHelper.ErrQuery {
			t.Errorf("unexpected err: %v", addErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestAddSkillsByCvIDSecondQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByCvIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(sql.ErrNoRows)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.skills[0]).
			WillReturnError(testCase.queryError)

		addErr := repo.AddSkillsByCvID(ctxWithLogger, testCase.cvID, testCase.skills)
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

func TestAddSkillsByCvIDThirdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testAddSkillsByCvIDQueryErrorCases {
		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(testCase.skills[0]).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(0),
			)

		mock.
			ExpectExec(testHelper.InsertQuery).
			WithArgs(testCase.cvID, 0).
			WillReturnError(testCase.queryError)

		addErr := repo.AddSkillsByCvID(ctxWithLogger, testCase.cvID, testCase.skills)
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

var testAddSkillsByCvIDOnExistance = []struct {
	cvID   int
	skills []string
}{
	{
		cvID:   1,
		skills: []string{"python", "c++"},
	},
	{
		cvID:   2,
		skills: []string{"js", "golang"},
	},
}

func TestAddSkillsByCvIDOnExistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByCvIDOnExistance {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnError(sql.ErrNoRows)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.cvID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.cvID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testGetSkillsByVacID = []struct {
	vacancyID      int
	expectedSkills []string
}{
	{
		vacancyID:      1,
		expectedSkills: []string{"python", "c++"},
	},
	{
		vacancyID:      2,
		expectedSkills: []string{"js", "golang"},
	},
	{
		vacancyID:      2,
		expectedSkills: []string{},
	},
}

func TestGetSkillsByVacID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testGetSkillsByVacID {
		rows := sqlmock.NewRows([]string{"skill"})
		for _, skill := range testCase.expectedSkills {
			rows.AddRow(skill)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnRows(rows)

		actualSkills, getErr := repo.GetSkillsByVacID(ctxWithLogger, testCase.vacancyID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actualSkills, testCase.expectedSkills) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedSkills, actualSkills))
		}
	}
}

var testGetSkillsByVacIDQueryErrorCases = []struct {
	vacancyID     int
	queryError    error
	expectedError error
}{
	{
		vacancyID:     1,
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
}

func TestGetSkillsByVacIDQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testGetSkillsByVacIDQueryErrorCases {

		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnError(testCase.queryError)

		_, getErr := repo.GetSkillsByVacID(ctxWithLogger, testCase.vacancyID)
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

var testGetSkillsByCvID = []struct {
	cvID           int
	expectedSkills []string
}{
	{
		cvID:           1,
		expectedSkills: []string{"python", "c++"},
	},
	{
		cvID:           2,
		expectedSkills: []string{"js", "golang"},
	},
	{
		cvID:           2,
		expectedSkills: []string{},
	},
}

func TestGetSkillsByCvID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testGetSkillsByCvID {
		rows := sqlmock.NewRows([]string{"skill"})
		for _, skill := range testCase.expectedSkills {
			rows.AddRow(skill)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnRows(rows)

		actualSkills, getErr := repo.GetSkillsByCvID(ctxWithLogger, testCase.cvID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actualSkills, testCase.expectedSkills) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedSkills, actualSkills))
		}
	}
}

var testGetSkillsByCvIDQueryErrorCases = []struct {
	cvID          int
	queryError    error
	expectedError error
}{
	{
		cvID:          1,
		queryError:    testHelper.ErrQuery,
		expectedError: testHelper.ErrQuery,
	},
}

func TestGetSkillsByCvIDQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, testCase := range testGetSkillsByCvIDQueryErrorCases {

		repo := psql.NewPsqlSkillRepository(db)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnError(testCase.queryError)

		_, getErr := repo.GetSkillsByCvID(ctxWithLogger, testCase.cvID)
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
