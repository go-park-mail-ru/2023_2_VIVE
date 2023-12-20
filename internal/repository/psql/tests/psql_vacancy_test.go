package psql

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/pkg/queryUtils"
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
		`"name"`,
		"description",
		"salary_lower_bound",
		"salary_upper_bound",
		"employment",
		"experience",
		"education_type",
		`"location"`,
		"created_at",
		"updated_at",
		"organization_name",
	}

	salaryLowerBound = 50000
	salaryUpperBound = 100000

	fullVacancyID1 = domain.DbVacancy{
		ID:               1,
		EmployerID:       1,
		VacancyName:      "Vacancy #1",
		Description:      "Description #1",
		SalaryLowerBound: &salaryLowerBound,
		SalaryUpperBound: &salaryUpperBound,
		Employment:       domain.FullTime,
		Experience:       domain.None,
		EducationType:    domain.Higher,
		Location:         &location,
		CreatedAt:        testHelper.Created_at,
		UpdatedAt:        testHelper.Updated_at,
	}
	fullVacancyID2 = domain.DbVacancy{
		ID:               2,
		EmployerID:       2,
		VacancyName:      "Vacancy #2",
		Description:      "Description #2",
		SalaryLowerBound: &salaryLowerBound,
		SalaryUpperBound: &salaryUpperBound,
		Employment:       domain.FullTime,
		Experience:       domain.None,
		EducationType:    domain.Higher,
		Location:         &location,
		CreatedAt:        testHelper.Created_at,
		UpdatedAt:        testHelper.Updated_at,
	}
	incompleteVacancyID3 = domain.DbVacancy{
		ID:               3,
		EmployerID:       1,
		VacancyName:      "Vacancy #1",
		Description:      "Description #1",
		SalaryLowerBound: nil,
		SalaryUpperBound: nil,
		Employment:       domain.FullTime,
		Experience:       domain.None,
		EducationType:    domain.Secondary,
		Location:         nil,
		CreatedAt:        testHelper.Created_at,
		UpdatedAt:        testHelper.Updated_at,
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

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetAllVacanciesSuccessCases {
		rows := sqlmock.NewRows(queryUtils.GetColumnNames(
			vacanciesColumns,
			"organization_name",
		))

		for _, item := range testCase.expected {
			rows = rows.AddRow(
				item.ID,
				item.EmployerID,
				item.VacancyName,
				item.Description,
				item.SalaryLowerBound,
				item.SalaryUpperBound,
				item.Employment,
				item.Experience,
				item.EducationType,
				item.Location,
				item.CreatedAt,
				item.UpdatedAt,
			)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WillReturnRows(rows)

		actual, err := repo.GetAllVacancies(testHelper.СtxWithLogger)
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

	repo := psql.NewPsqlVacancyRepository(db)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WillReturnError(testHelper.ErrQuery)

	_, returnedErr := repo.GetAllVacancies(testHelper.СtxWithLogger)
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

	repo := psql.NewPsqlVacancyRepository(db)

	rows := sqlmock.NewRows(vacanciesColumns)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WillReturnRows(rows)

	_, returnedErr := repo.GetAllVacancies(testHelper.СtxWithLogger)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != psql.ErrEntityNotFound {
		t.Errorf("expected error 'ErrEntityNotFound', got: '%s'", returnedErr)
		return
	}
}

var testGetEmpVacanciesByIdsSuccessCases = []struct {
	empID    int
	idList   []int
	expected []domain.DbVacancy
}{
	{
		empID:    1,
		idList:   []int{fullVacancyID1.ID, fullVacancyID2.ID},
		expected: []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
	},
	{
		empID:    2,
		idList:   []int{incompleteVacancyID3.ID},
		expected: []domain.DbVacancy{incompleteVacancyID3},
	},
}

func TestGetEmpVacanciesByIdsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetEmpVacanciesByIdsSuccessCases {
		rows := sqlmock.NewRows(queryUtils.GetColumnNames(
			vacanciesColumns,
			"organization_name",
		))

		for _, item := range testCase.expected {
			rows = rows.AddRow(
				item.ID,
				item.EmployerID,
				item.VacancyName,
				item.Description,
				item.SalaryLowerBound,
				item.SalaryUpperBound,
				item.Employment,
				item.Experience,
				item.EducationType,
				item.Location,
				item.CreatedAt,
				item.UpdatedAt,
			)
		}

		items := make([]int, len(testCase.idList)+1)
		items[0] = testCase.empID
		for i := 1; i < len(items); i++ {
			items[i] = testCase.idList[i-1]
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(items)...).
			WillReturnRows(rows)

		actual, err := repo.GetEmpVacanciesByIds(testHelper.СtxWithLogger, testCase.empID, testCase.idList)
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

func TestGetEmpVacanciesByIdsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	empID := 2
	idList := []int{1, 2}
	items := make([]int, len(idList)+1)
	items[0] = empID
	for i := 1; i < len(items); i++ {
		items[i] = idList[i-1]
	}
	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnError(testHelper.ErrQuery)

	_, returnedErr := repo.GetEmpVacanciesByIds(testHelper.СtxWithLogger, empID, idList)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr == nil {
		t.Errorf("expected query error, got: '%s'", returnedErr)
		return
	}
}

func TestGetEmpVacanciesByIdsEntityNotFoundError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	empID := 2
	idList := []int{1, 2}
	items := make([]int, len(idList)+1)
	items[0] = empID
	for i := 1; i < len(items); i++ {
		items[i] = idList[i-1]
	}
	rows := sqlmock.NewRows(vacanciesColumns)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnRows(rows)

	_, returnedErr := repo.GetEmpVacanciesByIds(testHelper.СtxWithLogger, empID, idList)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != psql.ErrEntityNotFound {
		t.Errorf("expected error 'ErrEntityNotFound', got: '%s'", returnedErr)
		return
	}
}

var testGetVacanciesByIdsSuccessCases = []struct {
	input    []int
	expected []domain.DbVacancy
}{
	{
		input:    []int{1, 2},
		expected: []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
	},
	{
		input:    []int{1},
		expected: []domain.DbVacancy{incompleteVacancyID3},
	},
}

func TestGetVacanciesByIdsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacanciesByIdsSuccessCases {
		rows := sqlmock.NewRows(vacanciesColumns)

		for _, item := range testCase.expected {
			rows = rows.AddRow(
				item.ID,
				item.EmployerID,
				item.VacancyName,
				item.Description,
				item.SalaryLowerBound,
				item.SalaryUpperBound,
				item.Employment,
				item.Experience,
				item.EducationType,
				item.Location,
				item.CreatedAt,
				item.UpdatedAt,
				item.OrganizationName,
			)
		}

		items := make([]int, len(testCase.input))
		// items[0] = testCase.orgID
		for i := 0; i < len(items); i++ {
			items[i] = testCase.input[i]
		}

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(items)...).
			WillReturnRows(rows)

		actual, err := repo.GetVacanciesByIds(testHelper.СtxWithLogger, testCase.input)
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

	repo := psql.NewPsqlVacancyRepository(db)

	_, err = repo.GetVacanciesByIds(testHelper.СtxWithLogger, []int{})
	if err != nil {
		t.Errorf("expected error 'ErrEntityNotFound', got %s", err)
	}
}

func TestGetVacanciesByIdsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	input := []int{1, 2, 3}
	// orgID := 1
	items := make([]int, len(input))
	// items[0] = orgID
	for i := 0; i < len(items); i++ {
		items[i] = input[i]
	}
	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnError(testHelper.ErrQuery)

	_, returnedErr := repo.GetVacanciesByIds(testHelper.СtxWithLogger, input)
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

	repo := psql.NewPsqlVacancyRepository(db)

	input := []int{1, 2, 3}
	// orgID := 1
	items := make([]int, len(input))
	// items[0] = orgID
	for i := 0; i < len(items); i++ {
		items[i] = input[i]
	}
	rows := sqlmock.NewRows(vacanciesColumns)
	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(testHelper.SliceIntToDriverValue(items)...).
		WillReturnRows(rows)

	_, returnedErr := repo.GetVacanciesByIds(testHelper.СtxWithLogger, input)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if returnedErr != nil {
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

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacancySuccessCases {
		rows := sqlmock.NewRows(vacanciesColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.EmployerID,
				testCase.expected.VacancyName,
				testCase.expected.Description,
				testCase.expected.SalaryLowerBound,
				testCase.expected.SalaryUpperBound,
				testCase.expected.Employment,
				testCase.expected.Experience,
				testCase.expected.EducationType,
				testCase.expected.Location,
				testCase.expected.OrganizationName,
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.input).
			WillReturnRows(rows)

		actual, err := repo.GetVacancy(testHelper.СtxWithLogger, testCase.input)
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
		expectedErr:  psql.ErrEntityNotFound,
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

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetVacancyQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.input).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetVacancy(testHelper.СtxWithLogger, testCase.input)
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

var testGetCompanyNameSuccessCases = []struct {
	vacancyID int
	expected  string
}{
	{
		vacancyID: 1,
		expected:  "organization1",
	},
	{
		vacancyID: 2,
		expected:  "organization1",
	},
	{
		vacancyID: 3,
		expected:  "organization3",
	},
}

func TestGetCompanyNameSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetCompanyNameSuccessCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnRows(
				sqlmock.NewRows([]string{"company_name"}).
					AddRow(testCase.expected),
			)

		actual, err := repo.GetCompanyName(testHelper.СtxWithLogger, testCase.vacancyID)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actual != testCase.expected {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
			return
		}
	}
}

var testGetCompanyNameQueryErrorCases = []struct {
	vacancyID    int
	returningErr error
	expectedErr  error
}{
	{
		vacancyID:    1,
		returningErr: sql.ErrNoRows,
		expectedErr:  psql.ErrEntityNotFound,
	},
	{
		vacancyID:    1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestGetCompanyNameQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetCompanyNameQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetCompanyName(testHelper.СtxWithLogger, testCase.vacancyID)
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

var testGetUserVacanciesSuccessCases = []struct {
	userID   int
	expected []domain.DbVacancy
}{
	{
		userID:   1,
		expected: []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
	},
	{
		userID:   2,
		expected: []domain.DbVacancy{incompleteVacancyID3},
	},
}

func TestGetUserVacanciesSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetUserVacanciesSuccessCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.userID).
			WillReturnRows(
				sqlmock.NewRows([]string{"employer_id"}).
					AddRow(1),
			)

		rows := sqlmock.NewRows(queryUtils.GetColumnNames(
			vacanciesColumns,
			"organization_name",
		))
		for _, vac := range testCase.expected {
			rows.AddRow(
				vac.ID,
				vac.EmployerID,
				vac.VacancyName,
				vac.Description,
				vac.SalaryLowerBound,
				vac.SalaryUpperBound,
				vac.Employment,
				vac.Experience,
				vac.EducationType,
				vac.Location,
				vac.CreatedAt,
				vac.UpdatedAt,
			)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(1).
			WillReturnRows(rows)

		actual, err := repo.GetUserVacancies(testHelper.СtxWithLogger, testCase.userID)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
			return
		}
	}
}

var testGetUserVacanciesFirstQueryErrorCases = []struct {
	userID       int
	returningErr error
	expectedErr  error
}{
	{
		userID:       1,
		returningErr: sql.ErrNoRows,
		expectedErr:  psql.ErrEntityNotFound,
	},
	{
		userID:       1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestGetUserVacanciesFirstQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetUserVacanciesFirstQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.userID).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetUserVacancies(testHelper.СtxWithLogger, testCase.userID)
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

var testGetUserVacanciesSecondQueryErrorCases = []struct {
	userID       int
	returningErr error
	expectedErr  error
}{
	{
		userID:       1,
		returningErr: sql.ErrNoRows,
		expectedErr:  nil,
	},
	{
		userID:       1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestGetUserVacanciesSecondQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetUserVacanciesSecondQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.userID).
			WillReturnRows(
				sqlmock.NewRows([]string{"employer_id"}).
					AddRow(1),
			)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(1).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetUserVacancies(testHelper.СtxWithLogger, testCase.userID)
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

var testGetEmployerInfoSuccessCases = []struct {
	isEmployer          bool
	employerID          int
	expectedFirstName   string
	expectedLastName    string
	expectedCompanyName string
	expectedVacancies   []domain.DbVacancy
	expectedErr         error
}{
	{
		isEmployer:          true,
		employerID:          1,
		expectedFirstName:   "first name",
		expectedLastName:    "last name",
		expectedCompanyName: "company name",
		expectedVacancies:   []domain.DbVacancy{fullVacancyID1, fullVacancyID2},
		expectedErr:         nil,
	},
	{
		isEmployer:          true,
		employerID:          2,
		expectedFirstName:   "first name",
		expectedLastName:    "last name",
		expectedCompanyName: "company name",
		expectedVacancies:   []domain.DbVacancy{incompleteVacancyID3},
		expectedErr:         nil,
	},
	{
		isEmployer:          false,
		employerID:          0,
		expectedFirstName:   "",
		expectedLastName:    "",
		expectedCompanyName: "",
		expectedVacancies:   nil,
		expectedErr:         psql.ErrEntityNotFound,
	},
}

func TestGetEmployerInfoSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetEmployerInfoSuccessCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.employerID).
			WillReturnRows(
				sqlmock.NewRows([]string{"is_employer"}).
					AddRow(testCase.isEmployer),
			)

		if testCase.isEmployer {
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(testCase.employerID).
				WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "company_name"}).
						AddRow(1, "company name"),
				)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(1).
				WillReturnRows(
					sqlmock.NewRows([]string{"first_name", "last_name"}).
						AddRow(testCase.expectedFirstName, testCase.expectedLastName),
				)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(1).
				WillReturnRows(
					sqlmock.NewRows([]string{"employer_id"}).
						AddRow(testCase.employerID),
				)

			rows := sqlmock.NewRows(queryUtils.GetColumnNames(
				vacanciesColumns,
				"organization_name",
			))
			for _, vac := range testCase.expectedVacancies {
				rows.AddRow(
					vac.ID,
					vac.EmployerID,
					vac.VacancyName,
					vac.Description,
					vac.SalaryLowerBound,
					vac.SalaryUpperBound,
					vac.Employment,
					vac.Experience,
					vac.EducationType,
					vac.Location,
					vac.CreatedAt,
					vac.UpdatedAt,
				)
			}
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(testCase.employerID).
				WillReturnRows(rows)
		}

		actualFirstName, actualLastName, actualCompName, actualVacancies, err := repo.GetEmployerInfo(testHelper.СtxWithLogger, testCase.employerID)
		if err != nil && err != testCase.expectedErr {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualFirstName != testCase.expectedFirstName {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedFirstName, actualFirstName))
			return
		}
		if actualLastName != testCase.expectedLastName {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedLastName, actualLastName))
			return
		}
		if actualCompName != testCase.expectedCompanyName {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedCompanyName, actualCompName))
			return
		}
		if !reflect.DeepEqual(actualVacancies, testCase.expectedVacancies) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedVacancies, actualVacancies))
			return
		}
	}
}

func TestGetEmployerInfoFirstQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	employerID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(employerID).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, err = repo.GetEmployerInfo(testHelper.СtxWithLogger, employerID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestGetEmployerInfoSecondQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	employerID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(employerID).
		WillReturnRows(
			sqlmock.NewRows([]string{"is_employer"}).
				AddRow(true),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(employerID).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, err = repo.GetEmployerInfo(testHelper.СtxWithLogger, employerID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestGetEmployerInfoThirdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	employerID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(employerID).
		WillReturnRows(
			sqlmock.NewRows([]string{"is_employer"}).
				AddRow(true),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(employerID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "company_name"}).
				AddRow(1, "company name"),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(1).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, err = repo.GetEmployerInfo(testHelper.СtxWithLogger, employerID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

var testGetEmpIdSuccessCases = []struct {
	vacancyID int
	expected  int
}{
	{
		vacancyID: 1,
		expected:  1,
	},
	{
		vacancyID: 2,
		expected:  2,
	},
	{
		vacancyID: 3,
		expected:  3,
	},
}

func TestGetEmpIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetEmpIdSuccessCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnRows(
				sqlmock.NewRows([]string{"company_name"}).
					AddRow(testCase.expected),
			)

		actual, err := repo.GetEmpId(testHelper.СtxWithLogger, testCase.vacancyID)
		if err != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actual != testCase.expected {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
			return
		}
	}
}

var testGetEmpIdQueryErrorCases = []struct {
	vacancyID    int
	returningErr error
	expectedErr  error
}{
	{
		vacancyID:    1,
		returningErr: sql.ErrNoRows,
		expectedErr:  psql.ErrEntityNotFound,
	},
	{
		vacancyID:    1,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestGetEmpIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testGetEmpIdQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.GetEmpId(testHelper.СtxWithLogger, testCase.vacancyID)
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

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testAddVacancySuccessCases {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(testCase.expected)

		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.Experience,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputUserID,
			).
			WillReturnRows(rows)

		actual, err := repo.AddVacancy(testHelper.СtxWithLogger, testCase.inputUserID, &testCase.inputVacancy)
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
		expectedErr:  psql.ErrNotInserted,
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
		expectedErr:  psql.ErrNotInserted,
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

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testAddVacancyErrorCases {
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.Experience,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningErr)

		_, actualErr := repo.AddVacancy(testHelper.СtxWithLogger, testCase.inputUserID, &testCase.inputVacancy)
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

var testUpdateEmpVacancySuccessCases = []struct {
	inputEmpID   int
	inputVacID   int
	inputVacancy domain.DbVacancy
}{
	{
		inputEmpID:   1,
		inputVacID:   1,
		inputVacancy: fullVacancyID1,
	},
	{
		inputEmpID:   1,
		inputVacID:   2,
		inputVacancy: fullVacancyID2,
	},
	{
		inputEmpID:   1,
		inputVacID:   3,
		inputVacancy: incompleteVacancyID3,
	},
}

func TestUpdateEmpVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testUpdateEmpVacancySuccessCases {
		mock.
			ExpectExec(testHelper.UpdateQuery).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.Experience,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputVacID,
				testCase.inputEmpID,
			).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.UpdateEmpVacancy(testHelper.СtxWithLogger, testCase.inputEmpID, testCase.inputVacID, &testCase.inputVacancy)
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

var testUpdateEmpVacancyErrorCases = []struct {
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
		expectedErr:  psql.ErrNoRowsUpdated,
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
		expectedErr:  psql.ErrNoRowsUpdated,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		inputVacancy: incompleteVacancyID3,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestUpdateEmpVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testUpdateEmpVacancyErrorCases {
		mock.
			ExpectExec(testHelper.UpdateQuery).
			WithArgs(
				testCase.inputVacancy.VacancyName,
				testCase.inputVacancy.Description,
				testCase.inputVacancy.SalaryLowerBound,
				testCase.inputVacancy.SalaryUpperBound,
				testCase.inputVacancy.Employment,
				testCase.inputVacancy.Experience,
				testCase.inputVacancy.EducationType,
				testCase.inputVacancy.Location,
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnError(testCase.returningErr)

		actualErr := repo.UpdateEmpVacancy(testHelper.СtxWithLogger, testCase.inputOrgID, testCase.inputVacID, &testCase.inputVacancy)
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

var testDeleteEmpVacancySuccessCases = []struct {
	inputEmpID int
	inputVacID int
}{
	{
		inputEmpID: 1,
		inputVacID: 1,
	},
	{
		inputEmpID: 1,
		inputVacID: 2,
	},
	{
		inputEmpID: 1,
		inputVacID: 3,
	},
}

func TestDeleteEmpVacancySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testDeleteEmpVacancySuccessCases {
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(
				testCase.inputVacID,
				testCase.inputEmpID,
			).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.DeleteEmpVacancy(testHelper.СtxWithLogger, testCase.inputEmpID, testCase.inputVacID)
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

var testDeleteEmpVacancyErrorCases = []struct {
	inputOrgID   int
	inputVacID   int
	returningErr error
	expectedErr  error
}{
	{
		inputOrgID:   1,
		inputVacID:   1,
		returningErr: sql.ErrNoRows,
		expectedErr:  psql.ErrNoRowsDeleted,
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
		expectedErr:  psql.ErrNoRowsDeleted,
	},
	{
		inputOrgID:   3,
		inputVacID:   3,
		returningErr: testHelper.ErrQuery,
		expectedErr:  testHelper.ErrQuery,
	},
}

func TestDeleteEmpVacancyQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlVacancyRepository(db)

	for _, testCase := range testDeleteEmpVacancyErrorCases {
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(
				testCase.inputVacID,
				testCase.inputOrgID,
			).
			WillReturnError(testCase.returningErr)

		actualErr := repo.DeleteEmpVacancy(testHelper.СtxWithLogger, testCase.inputOrgID, testCase.inputVacID)
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
