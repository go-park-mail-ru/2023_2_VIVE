package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/nullTypes"
	"HnH/pkg/testHelper"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	vacanciesColumns = []string{
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

	fullVacancyID1 = domain.DbVacancy{
		ID:                   1,
		EmployerID:           1,
		VacancyName:          "Vacancy #1",
		Description:          "Description #1",
		SalaryLowerBound:     nullTypes.NewNullInt(10000, true),
		SalaryUpperBound:     nullTypes.NewNullInt(20000, true),
		Employment:           nullTypes.NewNullString(string(domain.FullTime), true),
		ExperienceLowerBound: nullTypes.NewNullInt(0, true),
		ExperienceUpperBound: nullTypes.NewNullInt(12, true),
		EducationType:        domain.Higher,
		Location:             nullTypes.NewNullString("Moscow", true),
		CreatedAt:            testHelper.Created_at,
		UpdatedAt:            testHelper.Updated_at,
	}
	fullVacancyID2 = domain.DbVacancy{
		ID:                   2,
		EmployerID:           2,
		VacancyName:          "Vacancy #2",
		Description:          "Description #2",
		SalaryLowerBound:     nullTypes.NewNullInt(10000, true),
		SalaryUpperBound:     nullTypes.NewNullInt(20000, true),
		Employment:           nullTypes.NewNullString(string(domain.FullTime), true),
		ExperienceLowerBound: nullTypes.NewNullInt(0, true),
		ExperienceUpperBound: nullTypes.NewNullInt(12, true),
		EducationType:        domain.Higher,
		Location:             nullTypes.NewNullString("Moscow", true),
		CreatedAt:            testHelper.Created_at,
		UpdatedAt:            testHelper.Updated_at,
	}
	incompleteVacancyID3 = domain.DbVacancy{
		ID:                   3,
		EmployerID:           1,
		VacancyName:          "Vacancy #1",
		Description:          "Description #1",
		SalaryLowerBound:     nullTypes.NewNullInt(0, false),
		SalaryUpperBound:     nullTypes.NewNullInt(0, false),
		Employment:           nullTypes.NewNullString("", false),
		ExperienceLowerBound: nullTypes.NewNullInt(0, false),
		ExperienceUpperBound: nullTypes.NewNullInt(0, false),
		EducationType:        domain.Secondary,
		Location:             nullTypes.NewNullString("", false),
		CreatedAt:            testHelper.Created_at,
		UpdatedAt:            testHelper.Updated_at,
	}
)

var testGetAllVacanciesSuccessCases = []struct {
	expected []domain.DbVacancy
}{
	{
		expected: []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
	},
	{
		expected: []domain.DbVacancy{incompleteVacancyID3},
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
				item.EmployerID,
				item.VacancyName,
				item.Description,
				item.SalaryLowerBound.GetValue(),
				item.SalaryUpperBound.GetValue(),
				item.Employment.GetValue(),
				item.ExperienceLowerBound.GetValue(),
				item.ExperienceUpperBound.GetValue(),
				item.EducationType,
				item.Location.GetValue(),
				item.CreatedAt,
				item.UpdatedAt,
			)
		}
		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
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
		ExpectQuery(testHelper.SELECT_QUERY).
		WillReturnError(testHelper.ErrQuery)

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
		ExpectQuery(testHelper.SELECT_QUERY).
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
	orgID    int
	expected []domain.DbVacancy
}{
	{
		input:    []int{1, 2},
		orgID:    1,
		expected: []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
	},
	{
		input:    []int{1},
		orgID:    1,
		expected: []domain.DbVacancy{incompleteVacancyID3},
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
				item.EmployerID,
				item.VacancyName,
				item.Description,
				item.SalaryLowerBound.GetValue(),
				item.SalaryUpperBound.GetValue(),
				item.Employment.GetValue(),
				item.ExperienceLowerBound.GetValue(),
				item.ExperienceUpperBound.GetValue(),
				item.EducationType,
				item.Location.GetValue(),
				item.CreatedAt,
				item.UpdatedAt,
			)
		}

		items := make([]int, len(testCase.input)+1)
		items[0] = testCase.orgID
		for i := 1; i < len(items); i++ {
			items[i] = testCase.input[i-1]
		}

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testHelper.SliceIntToDriverValue(items)...).
			WillReturnRows(rows)

		actual, err := repo.GetVacanciesByIds(1, testCase.input)
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

	_, err = repo.GetVacanciesByIds(1, []int{})
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
	orgID := 1
	items := make([]int, len(input)+1)
	items[0] = orgID
	for i := 1; i < len(items); i++ {
		items[i] = input[i-1]
	}
	mock.
		ExpectQuery(testHelper.SELECT_QUERY).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnError(testHelper.ErrQuery)

	_, returnedErr := repo.GetVacanciesByIds(1, input)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != testHelper.ErrQuery {
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
	orgID := 1
	items := make([]int, len(input)+1)
	items[0] = orgID
	for i := 1; i < len(items); i++ {
		items[i] = input[i-1]
	}
	rows := sqlmock.NewRows(vacanciesColumns)
	mock.
		ExpectQuery(testHelper.SELECT_QUERY).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnRows(rows)

	_, returnedErr := repo.GetVacanciesByIds(1, input)
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
	expected domain.DbVacancy
}{
	{
		input:    1,
		expected: fullVacancyID1,
	},
	{
		input:    2,
		expected: fullVacancyID2,
	},
	{
		input:    1,
		expected: incompleteVacancyID3,
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
				testCase.expected.EmployerID,
				testCase.expected.VacancyName,
				testCase.expected.Description,
				testCase.expected.SalaryLowerBound.GetValue(),
				testCase.expected.SalaryUpperBound.GetValue(),
				testCase.expected.Employment.GetValue(),
				testCase.expected.ExperienceLowerBound.GetValue(),
				testCase.expected.ExperienceUpperBound.GetValue(),
				testCase.expected.EducationType,
				testCase.expected.Location.GetValue(),
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
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
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
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
			ExpectQuery(testHelper.SELECT_QUERY).
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
			ExpectQuery(testHelper.SELECT_QUERY).
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
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
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
			ExpectQuery(testHelper.SELECT_QUERY).
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

var testAddVacancySuccessCases = []struct {
	inputUserID  int
	inputVacancy domain.DbVacancy
	expected     int
}{
	{
		inputUserID:  1,
		inputVacancy: fullVacancyID1,
		expected:     1,
	},
	{
		inputUserID:  1,
		inputVacancy: fullVacancyID2,
		expected:     2,
	},
	{
		inputUserID:  1,
		inputVacancy: incompleteVacancyID3,
		expected:     3,
	},
}

func TestAddVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testAddVacancySuccessCases {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(testCase.expected)

		mock.
			ExpectQuery(testHelper.INSERT_QUERY).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.ExperienceLowerBound,
				testCase.inputVacancy.ExperienceUpperBound,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputUserID,
			).
			WillReturnRows(rows)

		actual, err := repo.AddVacancy(testCase.inputUserID, &testCase.inputVacancy)
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

var testAddVacancyErrorCases = []struct {
	inputUserID  int
	inputVacancy domain.DbVacancy
	returningErr error
	expectedErr  error
}{
	{
		inputUserID:  1,
		inputVacancy: fullVacancyID1,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNotInserted,
	},
	{
		inputUserID:  1,
		inputVacancy: fullVacancyID1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
	{
		inputUserID:  3,
		inputVacancy: incompleteVacancyID3,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNotInserted,
	},
	{
		inputUserID:  3,
		inputVacancy: incompleteVacancyID3,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestAddVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testAddVacancyErrorCases {
		mock.
			ExpectQuery(testHelper.INSERT_QUERY).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.ExperienceLowerBound,
				testCase.inputVacancy.ExperienceUpperBound,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.AddVacancy(testCase.inputUserID, &testCase.inputVacancy)
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

var testUpdateOrgVacancySuccessCases = []struct {
	inputOrgID   int
	inputVacID   int
	inputVacancy domain.DbVacancy
}{
	{
		inputOrgID:   1,
		inputVacID:   1,
		inputVacancy: fullVacancyID1,
	},
	{
		inputOrgID:   1,
		inputVacID:   2,
		inputVacancy: fullVacancyID2,
	},
	{
		inputOrgID:   1,
		inputVacID:   3,
		inputVacancy: incompleteVacancyID3,
	},
}

func TestUpdateOrgVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testUpdateOrgVacancySuccessCases {
		mock.
			ExpectExec(testHelper.UPDATE_QUERY).
			WithArgs(
				testCase.inputVacancy.EmployerID,
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.ExperienceLowerBound,
				testCase.inputVacancy.ExperienceUpperBound,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.UpdateOrgVacancy(testCase.inputOrgID, testCase.inputVacID, &testCase.inputVacancy)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testUpdateOrgVacancyErrorCases = []struct {
	inputOrgID   int
	inputVacID   int
	inputVacancy domain.DbVacancy
	returningErr error
	expectedErr  error
}{
	{
		inputOrgID:   1,
		inputVacID:   1,
		inputVacancy: fullVacancyID1,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNoRowsUpdated,
	},
	{
		inputOrgID:   1,
		inputVacID:   1,
		inputVacancy: fullVacancyID1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		inputVacancy: incompleteVacancyID3,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNoRowsUpdated,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		inputVacancy: incompleteVacancyID3,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestUpdateOrgVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testUpdateOrgVacancyErrorCases {
		mock.
			ExpectExec(testHelper.UPDATE_QUERY).
			WithArgs(
				testCase.inputVacancy.EmployerID,
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.ExperienceLowerBound,
				testCase.inputVacancy.ExperienceUpperBound,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnError(testCase.returningErr)

		actualErr := repo.UpdateOrgVacancy(testCase.inputOrgID, testCase.inputVacID, &testCase.inputVacancy)
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

var testDeleteOrgVacancySuccessCases = []struct {
	inputOrgID int
	inputVacID int
}{
	{
		inputOrgID: 1,
		inputVacID: 1,
	},
	{
		inputOrgID: 1,
		inputVacID: 2,
	},
	{
		inputOrgID: 1,
		inputVacID: 3,
	},
}

func TestDeleteOrgVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testDeleteOrgVacancySuccessCases {
		mock.
			ExpectExec(testHelper.DELETE_QUERY).
			WithArgs(
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.DeleteOrgVacancy(testCase.inputOrgID, testCase.inputVacID)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testDeleteOrgVacancyErrorCases = []struct {
	inputOrgID   int
	inputVacID   int
	returningErr error
	expectedErr  error
}{
	{
		inputOrgID:   1,
		inputVacID:   1,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNoRowsDeleted,
	},
	{
		inputOrgID:   1,
		inputVacID:   1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		returningErr: sql.ErrNoRows,
		expectedErr:  ErrNoRowsDeleted,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestDeleteOrgVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlVacancyRepository(db)

	for _, testCase := range testDeleteOrgVacancyErrorCases {
		mock.
			ExpectExec(testHelper.DELETE_QUERY).
			WithArgs(
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnError(testCase.returningErr)

		actualErr := repo.DeleteOrgVacancy(testCase.inputOrgID, testCase.inputVacID)
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
