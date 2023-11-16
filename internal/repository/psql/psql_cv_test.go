package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/testHelper"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	cvColumns = []string{
		"id",
		"applicant_id",
		"profession_name",
		"description",
		"status",
		"created_at",
		"updated_at",
	}

	cvID1 = domain.DbCV{
		ID:             1,
		ApplicantID:    1,
		ProfessionName: "Profession #1",
		Description:    "Description #1",
		Status:         domain.Searching,
		CreatedAt:      testHelper.Created_at,
		UpdatedAt:      testHelper.Updated_at,
	}
	cvID2 = domain.DbCV{
		ID:             2,
		ApplicantID:    2,
		ProfessionName: "Profession #2",
		Description:    "Description #2",
		Status:         domain.NotSearching,
		CreatedAt:      testHelper.Created_at,
		UpdatedAt:      testHelper.Updated_at,
	}
	cvID3 = domain.DbCV{
		ID:             3,
		ApplicantID:    3,
		ProfessionName: "Profession #3",
		Description:    "Description #3",
		Status:         domain.NotSearching,
		CreatedAt:      testHelper.Created_at,
		UpdatedAt:      testHelper.Updated_at,
	}
)

var testGetCVByIdSuccessCases = []struct {
	inputCVID int
	expected  domain.DbCV
}{
	{
		inputCVID: 1,
		expected:  cvID1,
	},
	{
		inputCVID: 2,
		expected:  cvID2,
	},
}

func TestGetCVByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVByIdSuccessCases {
		rows := sqlmock.NewRows(cvColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.ApplicantID,
				testCase.expected.ProfessionName,
				testCase.expected.Description,
				testCase.expected.Status,
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)

		actual, err := repo.GetCVById(testCase.inputCVID)
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

var testGetCVByIdQueryErrorCases = []struct {
	inputCVID      int
	returningError error
	expectedError  error
}{
	{
		inputCVID:      1,
		returningError: sql.ErrNoRows,
		expectedError:  ErrEntityNotFound,
	},
	{
		inputCVID:      2,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
}

func TestGetCVByIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVByIdQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputCVID).
			WillReturnError(testCase.returningError)

		_, actualErr := repo.GetCVById(testCase.inputCVID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

var testGetCVsByIdsSuccessCases = []struct {
	inputCVIDs []int
	expected   []domain.DbCV
}{
	{
		inputCVIDs: []int{1, 2, 3},
		expected:   []domain.DbCV{cvID1, cvID2, cvID3},
	},
	{
		inputCVIDs: []int{1, 2},
		expected:   []domain.DbCV{cvID1, cvID2},
	},
	{
		inputCVIDs: []int{1},
		expected:   []domain.DbCV{cvID1},
	},
	{
		inputCVIDs: []int{2},
		expected:   []domain.DbCV{cvID2},
	},
	{
		inputCVIDs: []int{3},
		expected:   []domain.DbCV{cvID3},
	},
}

func TestGetCVsByIdsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByIdsSuccessCases {
		rows := sqlmock.NewRows(cvColumns)
		for _, cv := range testCase.expected {
			rows = rows.AddRow(
				cv.ID,
				cv.ApplicantID,
				cv.ProfessionName,
				cv.Description,
				cv.Status,
				cv.CreatedAt,
				cv.UpdatedAt,
			)
		}

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnRows(rows)

		actual, err := repo.GetCVsByIds(testCase.inputCVIDs)
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

var testGetCVsByIdsQueryErrorCases = []struct {
	inputCVIDs     []int
	returningError error
	expectedError  error
}{
	{
		inputCVIDs:     []int{1, 2, 3},
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
}

func TestGetCVsByIdsQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByIdsQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnError(testCase.returningError)

		_, actualErr := repo.GetCVsByIds(testCase.inputCVIDs)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

func TestGetCVsByIdsErrEntityNotFoundEmptyArgs(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	inputCVIDs := []int{}
	expectedErr := ErrEntityNotFound

	_, actualErr := repo.GetCVsByIds(inputCVIDs)
	if actualErr != expectedErr {
		t.Errorf("expected query error: %s\ngot: '%s'", expectedErr, actualErr)
		return
	}
}

func TestGetCVsByIdsErrEntityNotFoundEmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	inputCVIDs := []int{1, 2, 3}
	expectedErr := ErrEntityNotFound
	mock.
		ExpectQuery(testHelper.SELECT_QUERY).
		WithArgs(testHelper.SliceIntToDriverValue(inputCVIDs)...).
		WillReturnError(expectedErr)

	_, actualErr := repo.GetCVsByIds(inputCVIDs)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if actualErr != expectedErr {
		t.Errorf("expected query error: %s\ngot: '%s'", expectedErr, actualErr)
		return
	}
}

var testGetCVsByUserIdSuccessCases = []struct {
	inputUserID int
	expected    []domain.DbCV
}{
	{
		inputUserID: 1,
		expected:    []domain.DbCV{cvID1, cvID2},
	},
	{
		inputUserID: 2,
		expected:    []domain.DbCV{cvID3},
	},
	{
		inputUserID: 3,
		expected:    []domain.DbCV{cvID1, cvID2, cvID3},
	},
}

func TestGetCVsByUserIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByUserIdSuccessCases {
		rows := sqlmock.NewRows(cvColumns)
		for _, cv := range testCase.expected {
			rows = rows.AddRow(
				cv.ID,
				cv.ApplicantID,
				cv.ProfessionName,
				cv.Description,
				cv.Status,
				cv.CreatedAt,
				cv.UpdatedAt,
			)
		}

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUserID).
			WillReturnRows(rows)

		actual, err := repo.GetCVsByUserId(testCase.inputUserID)
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

var testGetCVsByUserIdQueryErrorCases = []struct {
	inputUserID    int
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
}

func TestGetCVsByUserIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByUserIdQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUserID).
			WillReturnError(testCase.returningError)

		_, actualErr := repo.GetCVsByUserId(testCase.inputUserID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

func TestGetCVsByUserIdErrEntityNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	inputUserID := 1
	expectedError := ErrEntityNotFound

	mock.
		ExpectQuery(testHelper.SELECT_QUERY).
		WithArgs(inputUserID).
		WillReturnRows(sqlmock.NewRows(cvColumns))

	_, actualErr := repo.GetCVsByUserId(inputUserID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if actualErr != expectedError {
		t.Errorf("expected query error: %s\ngot: '%s'", expectedError, actualErr)
		return
	}
}

var testAddCVSuccessCases = []struct {
	inputUserID int
	inputCV     domain.DbCV
	expected    int
}{
	{
		inputUserID: 1,
		inputCV:     cvID1,
		expected:    1,
	},
	{
		inputUserID: 2,
		inputCV:     cvID2,
		expected:    2,
	},
	{
		inputUserID: 3,
		inputCV:     cvID3,
		expected:    3,
	},
}

func TestAddCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testAddCVSuccessCases {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(testCase.expected)

		mock.
			ExpectQuery(testHelper.INSERT_QUERY).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCV.Status,
				testCase.inputUserID,
			).
			WillReturnRows(rows)

		actual, err := repo.AddCV(testCase.inputUserID, &testCase.inputCV)
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

var testAddCVQueryErrorCases = []struct {
	inputUserID    int
	inputCV        domain.DbCV
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		inputCV:        cvID1,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
	{
		inputUserID:    1,
		inputCV:        cvID1,
		returningError: sql.ErrNoRows,
		expectedError:  ErrNotInserted,
	},
}

func TestAddCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testAddCVQueryErrorCases {
		mock.
			ExpectQuery(testHelper.INSERT_QUERY).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCV.Status,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningError)

		_, actualErr := repo.AddCV(testCase.inputUserID, &testCase.inputCV)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

var testGetOneOfUsersCVSuccessCases = []struct {
	inputUserID int
	inputCVID   int
	expected    domain.DbCV
}{
	{
		inputUserID: 1,
		inputCVID:   1,
		expected:    cvID1,
	},
	{
		inputUserID: 2,
		inputCVID:   2,
		expected:    cvID2,
	},
	{
		inputUserID: 3,
		inputCVID:   3,
		expected:    cvID3,
	},
}

func TestGetOneOfUsersCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetOneOfUsersCVSuccessCases {
		rows := sqlmock.NewRows(cvColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.ApplicantID,
				testCase.expected.ProfessionName,
				testCase.expected.Description,
				testCase.expected.Status,
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUserID, testCase.inputCVID).
			WillReturnRows(rows)

		actual, err := repo.GetOneOfUsersCV(testCase.inputUserID, testCase.inputCVID)
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

var testGetOneOfUsersCVQueryErrorCases = []struct {
	inputUserID    int
	inputCVID      int
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		inputCVID:      1,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
	{
		inputUserID:    1,
		inputCVID:      1,
		returningError: sql.ErrNoRows,
		expectedError:  ErrEntityNotFound,
	},
}

func TestGetOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testGetOneOfUsersCVQueryErrorCases {
		mock.
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUserID, testCase.inputCVID).
			WillReturnError(testCase.returningError)

		_, actualErr := repo.GetOneOfUsersCV(testCase.inputUserID, testCase.inputCVID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

var testUpdateOneOfUsersCVSuccessCases = []struct {
	inputUserID int
	inputCVID   int
	inputCV     domain.DbCV
}{
	{
		inputUserID: 1,
		inputCVID:   1,
		inputCV:     cvID1,
	},
	{
		inputUserID: 2,
		inputCVID:   2,
		inputCV:     cvID2,
	},
	{
		inputUserID: 3,
		inputCVID:   3,
		inputCV:     cvID2,
	},
}

func TestUpdateOneOfUsersCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testUpdateOneOfUsersCVSuccessCases {
		mock.
			ExpectExec(testHelper.UPDATE_QUERY).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCVID,
				testCase.inputUserID,
			).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.UpdateOneOfUsersCV(testCase.inputUserID, testCase.inputCVID, &testCase.inputCV)
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

var testUpdateOneOfUsersCVQueryErrorCases = []struct {
	inputUserID    int
	inputCVID      int
	inputCV        domain.DbCV
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		inputCVID:      1,
		inputCV:        cvID1,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
	{
		inputUserID:    1,
		inputCVID:      1,
		inputCV:        cvID1,
		returningError: sql.ErrNoRows,
		expectedError:  ErrNoRowsUpdated,
	},
}

func TestUpdateOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testUpdateOneOfUsersCVQueryErrorCases {
		mock.
			ExpectExec(testHelper.UPDATE_QUERY).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCVID,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningError)

		actualErr := repo.UpdateOneOfUsersCV(testCase.inputUserID, testCase.inputCVID, &testCase.inputCV)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}

var testDeleteOneOfUsersCVSuccessCases = []struct {
	inputUserID int
	inputCVID   int
}{
	{
		inputUserID: 1,
		inputCVID:   1,
	},
	{
		inputUserID: 2,
		inputCVID:   2,
	},
	{
		inputUserID: 3,
		inputCVID:   3,
	},
}

func TestDeleteOneOfUsersCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testDeleteOneOfUsersCVSuccessCases {
		mock.
			ExpectExec(testHelper.DELETE_QUERY).
			WithArgs(testCase.inputCVID, testCase.inputUserID).
			WillReturnResult(driver.RowsAffected(1))

		err := repo.DeleteOneOfUsersCV(testCase.inputUserID, testCase.inputCVID)
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

var testDeleteOneOfUsersCVQueryErrorCases = []struct {
	inputUserID    int
	inputCVID      int
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		inputCVID:      1,
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
	{
		inputUserID:    1,
		inputCVID:      1,
		returningError: sql.ErrNoRows,
		expectedError:  ErrNoRowsDeleted,
	},
}

func TestDeleteOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlCVRepository(db)

	for _, testCase := range testDeleteOneOfUsersCVQueryErrorCases {
		mock.
			ExpectExec(testHelper.DELETE_QUERY).
			WithArgs(testCase.inputCVID, testCase.inputUserID).
			WillReturnError(testCase.returningError)

		actualErr := repo.DeleteOneOfUsersCV(testCase.inputUserID, testCase.inputCVID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actualErr != testCase.expectedError {
			t.Errorf("expected query error: %s\ngot: '%s'", testCase.expectedError, actualErr)
			return
		}
	}
}
