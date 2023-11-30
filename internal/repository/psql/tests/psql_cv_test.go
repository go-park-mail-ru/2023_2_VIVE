package psql

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
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
		"name",
		"first_name",
		"last_name",
		"middle_name",
		"gender",
		"birthday",
		"city",
		"description",
		"status",
		"education_level",
		"created_at",
		"updated_at",
	}

	middleName  = "Middle name"
	description = "description"

	cvID1 = domain.DbCV{
		ID:             1,
		ApplicantID:    1,
		ProfessionName: "Profession #1",
		FirstName:      "First name",
		LastName:       "Last name",
		MiddleName:     &middleName,
		Gender:         domain.Female,
		Birthday:       &birthday,
		Location:       &location,
		Description:    &description,
		Status:         domain.Searching,
		EducationLevel: domain.Bachelor,
		CreatedAt:      testHelper.Created_at,
		UpdatedAt:      testHelper.Updated_at,
	}
	cvID2 = domain.DbCV{
		ID:             2,
		ApplicantID:    2,
		ProfessionName: "Profession #2",
		FirstName:      "First name",
		LastName:       "Last name",
		MiddleName:     &middleName,
		Gender:         domain.Female,
		Birthday:       &birthday,
		Location:       &location,
		Description:    &description,
		Status:         domain.NotSearching,
		EducationLevel: domain.Bachelor,
		CreatedAt:      testHelper.Created_at,
		UpdatedAt:      testHelper.Updated_at,
	}
	cvID3 = domain.DbCV{
		ID:             3,
		ApplicantID:    3,
		ProfessionName: "Profession #3",
		FirstName:      "First name",
		LastName:       "Last name",
		MiddleName:     &middleName,
		Gender:         domain.Female,
		Birthday:       &birthday,
		Location:       &location,
		Description:    &description,
		Status:         domain.NotSearching,
		EducationLevel: domain.Bachelor,
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVByIdSuccessCases {
		rows := sqlmock.NewRows(cvColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.ApplicantID,
				testCase.expected.ProfessionName,
				testCase.expected.FirstName,
				testCase.expected.LastName,
				testCase.expected.MiddleName,
				testCase.expected.Gender,
				testCase.expected.Birthday,
				testCase.expected.Location,
				testCase.expected.Description,
				testCase.expected.EducationLevel,
				testCase.expected.Status,
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)
		mock.ExpectCommit()

		actual, _, _, err := repo.GetCVById(testHelper.СtxWithLogger, testCase.inputCVID)
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
		expectedError:  psql.ErrEntityNotFound,
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVByIdQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		_, _, _, actualErr := repo.GetCVById(testHelper.СtxWithLogger, testCase.inputCVID)
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByIdsSuccessCases {
		rows := sqlmock.NewRows(cvColumns)
		for _, cv := range testCase.expected {
			rows = rows.AddRow(
				cv.ID,
				cv.ApplicantID,
				cv.ProfessionName,
				cv.FirstName,
				cv.LastName,
				cv.MiddleName,
				cv.Gender,
				cv.Birthday,
				cv.Location,
				cv.Description,
				cv.EducationLevel,
				cv.Status,
				cv.CreatedAt,
				cv.UpdatedAt,
			)
		}

		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnRows(rows)
		mock.ExpectCommit()

		actual, _, _, err := repo.GetCVsByIds(testHelper.СtxWithLogger, testCase.inputCVIDs)
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByIdsQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.inputCVIDs)...).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		_, _, _, actualErr := repo.GetCVsByIds(testHelper.СtxWithLogger, testCase.inputCVIDs)
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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	inputCVIDs := []int{}
	expectedErr := psql.ErrEntityNotFound

	mock.ExpectBegin()
	mock.ExpectCommit()

	_, _, _, actualErr := repo.GetCVsByIds(testHelper.СtxWithLogger, inputCVIDs)
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

	repo := psql.NewPsqlCVRepository(db)

	inputCVIDs := []int{1, 2, 3}
	expectedErr := psql.ErrEntityNotFound
	mock.ExpectBegin()
	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(testHelper.SliceIntToDriverValue(inputCVIDs)...).
		WillReturnError(expectedErr)
	mock.ExpectRollback()

	_, _, _, actualErr := repo.GetCVsByIds(testHelper.СtxWithLogger, inputCVIDs)
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
	expectedIDs []int
	expected    []domain.DbCV
}{
	{
		inputUserID: 1,
		expectedIDs: []int{cvID1.ID, cvID2.ID},
		expected:    []domain.DbCV{cvID1, cvID2},
	},
	{
		inputUserID: 2,
		expectedIDs: []int{cvID3.ID},
		expected:    []domain.DbCV{cvID3},
	},
	{
		inputUserID: 3,
		expectedIDs: []int{cvID1.ID, cvID2.ID, cvID3.ID},
		expected:    []domain.DbCV{cvID1, cvID2, cvID3},
	},
}

func TestGetCVsByUserIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByUserIdSuccessCases {
		rows := sqlmock.NewRows(cvColumns)
		for _, cv := range testCase.expected {
			rows = rows.AddRow(
				cv.ID,
				cv.ApplicantID,
				cv.ProfessionName,
				cv.FirstName,
				cv.LastName,
				cv.MiddleName,
				cv.Gender,
				cv.Birthday,
				cv.Location,
				cv.Description,
				cv.EducationLevel,
				cv.Status,
				cv.CreatedAt,
				cv.UpdatedAt,
			)
		}

		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.expectedIDs)...).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testHelper.SliceIntToDriverValue(testCase.expectedIDs)...).
			WillReturnRows(rows)
		mock.ExpectCommit()

		actual, _, _, err := repo.GetCVsByUserId(testHelper.СtxWithLogger, testCase.inputUserID)
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetCVsByUserIdQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		_, _, _, actualErr := repo.GetCVsByUserId(testHelper.СtxWithLogger, testCase.inputUserID)
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

	repo := psql.NewPsqlCVRepository(db)

	inputUserID := 1
	expectedError := psql.ErrEntityNotFound

	mock.ExpectBegin()
	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(inputUserID).
		WillReturnRows(sqlmock.NewRows(cvColumns))
	mock.ExpectRollback()

	_, _, _, actualErr := repo.GetCVsByUserId(testHelper.СtxWithLogger, inputUserID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if actualErr != expectedError {
		t.Errorf("expected query error: %s\ngot: '%s'", expectedError, actualErr)
		return
	}
}

var testGetApplicantInfoSuccessCases = []struct {
	isApplicant       bool
	applicantID       int
	expectedFirstName string
	expectedLastName  string
	expectedCVs       []domain.DbCV
	expectedExps      []domain.DbExperience
	expectedInsts     []domain.DbEducationInstitution
	expectedErr       error
}{
	{
		isApplicant:       true,
		applicantID:       1,
		expectedFirstName: "first name",
		expectedLastName:  "last name",
		expectedCVs:       []domain.DbCV{cvID1, cvID2},
		expectedExps:      []domain.DbExperience{},
		expectedInsts:     []domain.DbEducationInstitution{},
		expectedErr:       nil,
	},
	{
		isApplicant:       true,
		applicantID:       2,
		expectedFirstName: "first name",
		expectedLastName:  "last name",
		expectedCVs:       []domain.DbCV{cvID3},
		expectedExps:      []domain.DbExperience{},
		expectedInsts:     []domain.DbEducationInstitution{},
		expectedErr:       nil,
	},
	{
		isApplicant:       false,
		applicantID:       2,
		expectedFirstName: "",
		expectedLastName:  "",
		expectedCVs:       nil,
		expectedExps:      nil,
		expectedInsts:     nil,
		expectedErr:       psql.ErrEntityNotFound,
	},
}

func TestGetApplicantInfoSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetApplicantInfoSuccessCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.applicantID).
			WillReturnRows(
				sqlmock.NewRows([]string{"is_applicant"}).
					AddRow(testCase.isApplicant),
			)

		if testCase.isApplicant {
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(testCase.applicantID).
				WillReturnRows(
					sqlmock.NewRows([]string{"user_id"}).
						AddRow(1),
				)

			if testCase.isApplicant {
				mock.
					ExpectQuery(testHelper.SelectQuery).
					WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"first_name", "last_name"}).
							AddRow(testCase.expectedFirstName, testCase.expectedLastName),
					)
			}

			rows := sqlmock.NewRows(cvColumns)
			for _, cv := range testCase.expectedCVs {
				rows = rows.AddRow(
					cv.ID,
					cv.ApplicantID,
					cv.ProfessionName,
					cv.FirstName,
					cv.LastName,
					cv.MiddleName,
					cv.Gender,
					cv.Birthday,
					cv.Location,
					cv.Description,
					cv.EducationLevel,
					cv.Status,
					cv.CreatedAt,
					cv.UpdatedAt,
				)
			}

			mock.ExpectBegin()
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(1).
				WillReturnRows(rows)
			expectedIDs := make([]int, len(testCase.expectedCVs))
			for i := range testCase.expectedCVs {
				expectedIDs[i] = testCase.expectedCVs[i].ID
			}
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(testHelper.SliceIntToDriverValue(expectedIDs)...).
				WillReturnRows(rows)
			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(testHelper.SliceIntToDriverValue(expectedIDs)...).
				WillReturnRows(rows)
			mock.ExpectCommit()
		}

		aFirstName, aLastName, aCVs, aExps, aInsts, err := repo.GetApplicantInfo(testHelper.СtxWithLogger, testCase.applicantID)
		if err != nil && err != testCase.expectedErr {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if aFirstName != testCase.expectedFirstName {
			t.Errorf("results not match, want %v, have %v", testCase.expectedFirstName, aFirstName)
			return
		}
		if aLastName != testCase.expectedLastName {
			t.Errorf("results not match, want %v, have %v", testCase.expectedLastName, aLastName)
			return
		}
		if !reflect.DeepEqual(aCVs, testCase.expectedCVs) {
			t.Errorf("results not match, want %v, have %v", testCase.expectedCVs, aCVs)
			return
		}
		if !reflect.DeepEqual(aExps, testCase.expectedExps) {
			t.Errorf("results not match, want %v, have %v", testCase.expectedExps, aExps)
			return
		}
		if !reflect.DeepEqual(aInsts, testCase.expectedInsts) {
			t.Errorf("results not match, want %v, have %v", testCase.expectedInsts, aInsts)
			return
		}
	}
}

func TestGetApplicantInfoFirstQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	applicantID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(applicantID).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, _, err = repo.GetApplicantInfo(testHelper.СtxWithLogger, applicantID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestGetApplicantInfoSecondQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	applicantID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(applicantID).
		WillReturnRows(
			sqlmock.NewRows([]string{"is_applicant"}).
				AddRow(true),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(applicantID).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, _, err = repo.GetApplicantInfo(testHelper.СtxWithLogger, applicantID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestGetApplicantInfoThirdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	applicantID := 1

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(applicantID).
		WillReturnRows(
			sqlmock.NewRows([]string{"is_applicant"}).
				AddRow(true),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(applicantID).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id"}).
				AddRow(1),
		)

	mock.
		ExpectQuery(testHelper.SelectQuery).
		WithArgs(1).
		WillReturnError(testHelper.ErrQuery)

	_, _, _, _, _, err = repo.GetApplicantInfo(testHelper.СtxWithLogger, applicantID)
	if err != testHelper.ErrQuery {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

var testAddCVSuccessCases = []struct {
	inputUserID int
	inputCV     domain.DbCV
	inputExps   []domain.DbExperience
	inputInsts  []domain.DbEducationInstitution
	expected    int
}{
	{
		inputUserID: 1,
		inputCV:     cvID1,
		inputExps:   []domain.DbExperience{},
		inputInsts:  []domain.DbEducationInstitution{},
		expected:    1,
	},
	{
		inputUserID: 2,
		inputCV:     cvID2,
		inputExps:   []domain.DbExperience{},
		inputInsts:  []domain.DbEducationInstitution{},
		expected:    2,
	},
	{
		inputUserID: 3,
		inputCV:     cvID3,
		inputExps:   []domain.DbExperience{},
		inputInsts:  []domain.DbEducationInstitution{},
		expected:    3,
	},
}

func TestAddCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testAddCVSuccessCases {
		cvRows := sqlmock.NewRows([]string{"id"}).
			AddRow(testCase.expected)
		// expRows := sqlmock.NewRows([]string{""})

		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.FirstName,
				testCase.inputCV.LastName,
				testCase.inputCV.MiddleName,
				testCase.inputCV.Gender,
				testCase.inputCV.Birthday,
				testCase.inputCV.Location,
				testCase.inputCV.Description,
				testCase.inputCV.EducationLevel,
				testCase.inputUserID,
			).
			WillReturnRows(cvRows)
		// mock.
		// 	ExpectExec(testHelper.INSERT_QUERY).
		// 	WithArgs(testCase.inputCV.ID)
		// mock.
		// 	ExpectExec(testHelper.INSERT_QUERY).
		// 	WithArgs(testCase.inputCV.ID)

		mock.ExpectCommit()

		actual, err := repo.AddCV(testHelper.СtxWithLogger, testCase.inputUserID, &testCase.inputCV, testCase.inputExps, testCase.inputInsts)
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
	inputExps      []domain.DbExperience
	inputInsts     []domain.DbEducationInstitution
	returningError error
	expectedError  error
}{
	{
		inputUserID:    1,
		inputCV:        cvID1,
		inputExps:      []domain.DbExperience{},
		inputInsts:     []domain.DbEducationInstitution{},
		returningError: testHelper.ErrQuery,
		expectedError:  testHelper.ErrQuery,
	},
	{
		inputUserID:    1,
		inputCV:        cvID1,
		inputExps:      []domain.DbExperience{},
		inputInsts:     []domain.DbEducationInstitution{},
		returningError: sql.ErrNoRows,
		expectedError:  psql.ErrNotInserted,
	},
}

func TestAddCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testAddCVQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.FirstName,
				testCase.inputCV.LastName,
				testCase.inputCV.MiddleName,
				testCase.inputCV.Gender,
				testCase.inputCV.Birthday,
				testCase.inputCV.Location,
				testCase.inputCV.Description,
				testCase.inputCV.EducationLevel,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		_, actualErr := repo.AddCV(testHelper.СtxWithLogger, testCase.inputUserID, &testCase.inputCV, testCase.inputExps, testCase.inputInsts)
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetOneOfUsersCVSuccessCases {
		rows := sqlmock.NewRows(cvColumns).
			AddRow(
				testCase.expected.ID,
				testCase.expected.ApplicantID,
				testCase.expected.ProfessionName,
				testCase.expected.FirstName,
				testCase.expected.LastName,
				testCase.expected.MiddleName,
				testCase.expected.Gender,
				testCase.expected.Birthday,
				testCase.expected.Location,
				testCase.expected.Description,
				testCase.expected.EducationLevel,
				testCase.expected.Status,
				testCase.expected.CreatedAt,
				testCase.expected.UpdatedAt,
			)

		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID, testCase.inputCVID).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputCVID).
			WillReturnRows(rows)
		mock.ExpectCommit()

		actual, _, _, err := repo.GetOneOfUsersCV(testHelper.СtxWithLogger, testCase.inputUserID, testCase.inputCVID)
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
		expectedError:  psql.ErrEntityNotFound,
	},
}

func TestGetOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testGetOneOfUsersCVQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID, testCase.inputCVID).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		_, _, _, actualErr := repo.GetOneOfUsersCV(testHelper.СtxWithLogger, testCase.inputUserID, testCase.inputCVID)
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
	inputUserID           int
	inputCVID             int
	inputCV               domain.DbCV
	inputExpsIDsToDelete  []int
	inputExpsToUpdate     []domain.DbExperience
	inputExpsToInsert     []domain.DbExperience
	inputInstsIDsToDelete []int
	inputInstsToUpdate    []domain.DbEducationInstitution
	inputInstsToInsert    []domain.DbEducationInstitution
}{
	{
		inputUserID:           1,
		inputCVID:             1,
		inputCV:               cvID1,
		inputExpsIDsToDelete:  []int{},
		inputExpsToUpdate:     []domain.DbExperience{},
		inputExpsToInsert:     []domain.DbExperience{},
		inputInstsIDsToDelete: []int{},
		inputInstsToUpdate:    []domain.DbEducationInstitution{},
		inputInstsToInsert:    []domain.DbEducationInstitution{},
	},
	{
		inputUserID:           2,
		inputCVID:             2,
		inputCV:               cvID2,
		inputExpsIDsToDelete:  []int{},
		inputExpsToUpdate:     []domain.DbExperience{},
		inputExpsToInsert:     []domain.DbExperience{},
		inputInstsIDsToDelete: []int{},
		inputInstsToUpdate:    []domain.DbEducationInstitution{},
		inputInstsToInsert:    []domain.DbEducationInstitution{},
	},
	{
		inputUserID:           3,
		inputCVID:             3,
		inputCV:               cvID2,
		inputExpsIDsToDelete:  []int{},
		inputExpsToUpdate:     []domain.DbExperience{},
		inputExpsToInsert:     []domain.DbExperience{},
		inputInstsIDsToDelete: []int{},
		inputInstsToUpdate:    []domain.DbEducationInstitution{},
		inputInstsToInsert:    []domain.DbEducationInstitution{},
	},
}

func TestUpdateOneOfUsersCVSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testUpdateOneOfUsersCVSuccessCases {
		mock.ExpectBegin()
		mock.
			ExpectExec(testHelper.UpdateQuery).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.FirstName,
				testCase.inputCV.LastName,
				testCase.inputCV.MiddleName,
				testCase.inputCV.Gender,
				testCase.inputCV.Birthday,
				testCase.inputCV.Location,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCV.EducationLevel,
				testCase.inputCVID,
				testCase.inputUserID,
			).
			WillReturnResult(driver.RowsAffected(1))
		// mock.
		// 	ExpectExec(testHelper.UPDATE_QUERY).
		// 	WithArgs(
		// 		testCase.inputCV.ProfessionName,
		// 		testCase.inputCV.Description,
		// 		testCase.inputCV.Status,
		// 		testCase.inputCVID,
		// 		testCase.inputUserID,
		// 	)
		mock.ExpectCommit()

		err := repo.UpdateOneOfUsersCV(
			testHelper.СtxWithLogger,
			testCase.inputUserID,
			testCase.inputCVID,
			&testCase.inputCV,
			testCase.inputExpsIDsToDelete,
			testCase.inputExpsToUpdate,
			testCase.inputExpsToInsert,
			testCase.inputInstsIDsToDelete,
			testCase.inputInstsToUpdate,
			testCase.inputInstsToInsert,
		)
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
	inputUserID           int
	inputCVID             int
	inputCV               domain.DbCV
	inputExpsIDsToDelete  []int
	inputExpsToUpdate     []domain.DbExperience
	inputExpsToInsert     []domain.DbExperience
	inputInstsIDsToDelete []int
	inputInstsToUpdate    []domain.DbEducationInstitution
	inputInstsToInsert    []domain.DbEducationInstitution
	returningError        error
	expectedError         error
}{
	{
		inputUserID:           1,
		inputCVID:             1,
		inputCV:               cvID1,
		inputExpsIDsToDelete:  []int{},
		inputExpsToUpdate:     []domain.DbExperience{},
		inputExpsToInsert:     []domain.DbExperience{},
		inputInstsIDsToDelete: []int{},
		inputInstsToUpdate:    []domain.DbEducationInstitution{},
		inputInstsToInsert:    []domain.DbEducationInstitution{},
		returningError:        testHelper.ErrQuery,
		expectedError:         testHelper.ErrQuery,
	},
	{
		inputUserID:           1,
		inputCVID:             1,
		inputCV:               cvID1,
		inputExpsIDsToDelete:  []int{},
		inputExpsToUpdate:     []domain.DbExperience{},
		inputExpsToInsert:     []domain.DbExperience{},
		inputInstsIDsToDelete: []int{},
		inputInstsToUpdate:    []domain.DbEducationInstitution{},
		inputInstsToInsert:    []domain.DbEducationInstitution{},
		returningError:        sql.ErrNoRows,
		expectedError:         psql.ErrNoRowsUpdated,
	},
}

func TestUpdateOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testUpdateOneOfUsersCVQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectExec(testHelper.UpdateQuery).
			WithArgs(
				testCase.inputCV.ProfessionName,
				testCase.inputCV.FirstName,
				testCase.inputCV.LastName,
				testCase.inputCV.MiddleName,
				testCase.inputCV.Gender,
				testCase.inputCV.Birthday,
				testCase.inputCV.Location,
				testCase.inputCV.Description,
				testCase.inputCV.Status,
				testCase.inputCV.EducationLevel,
				testCase.inputCVID,
				testCase.inputUserID,
			).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		actualErr := repo.UpdateOneOfUsersCV(
			testHelper.СtxWithLogger,
			testCase.inputUserID,
			testCase.inputCVID,
			&testCase.inputCV,
			testCase.inputExpsIDsToDelete,
			testCase.inputExpsToUpdate,
			testCase.inputExpsToInsert,
			testCase.inputInstsIDsToDelete,
			testCase.inputInstsToUpdate,
			testCase.inputInstsToInsert,
		)
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

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testDeleteOneOfUsersCVSuccessCases {
		mock.ExpectBegin()
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(testCase.inputCVID, testCase.inputUserID).
			WillReturnResult(driver.RowsAffected(1))
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(testCase.inputCVID).
			WillReturnResult(driver.RowsAffected(1))
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(testCase.inputCVID).
			WillReturnResult(driver.RowsAffected(1))
		mock.ExpectCommit()

		err := repo.DeleteOneOfUsersCV(testHelper.СtxWithLogger, testCase.inputUserID, testCase.inputCVID)
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
		expectedError:  psql.ErrNoRowsDeleted,
	},
}

func TestDeleteOneOfUsersCVQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlCVRepository(db)

	for _, testCase := range testDeleteOneOfUsersCVQueryErrorCases {
		mock.ExpectBegin()
		mock.
			ExpectExec(testHelper.DeleteQuery).
			WithArgs(testCase.inputCVID, testCase.inputUserID).
			WillReturnError(testCase.returningError)
		mock.ExpectRollback()

		actualErr := repo.DeleteOneOfUsersCV(testHelper.СtxWithLogger, testCase.inputUserID, testCase.inputCVID)
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
