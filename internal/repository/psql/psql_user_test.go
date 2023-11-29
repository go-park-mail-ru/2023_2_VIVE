package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"context"
	"database/sql"
	"database/sql/driver"

	"HnH/pkg/testHelper"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	ctxWithLogger = context.WithValue(context.Background(), contextUtils.LOGGER_KEY, testHelper.InitCtxLogger())

	password1                 = "password_number_1"
	hashedPassword1, salt1, _ = authUtils.GenerateHash(password1)
	birthday                  = "2000-01-01"
	phone_number              = "+71111111111"
	location                  = "Moscow"
	appID                     = 1
	empID                     = 2

	applicant1 = domain.DbUser{
		ID:          1,
		Email:       "applicant1@example.com",
		Password:    password1,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
		Type:        domain.Applicant,
	}
	employer1 = domain.DbUser{
		ID:          1,
		Email:       "employer1@example.com",
		Password:    password1,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
		Type:        domain.Employer,
	}

	password2                 = "password_number_2"
	hashedPassword2, salt2, _ = authUtils.GenerateHash(password2)
	applicant2                = domain.DbUser{
		ID:          1,
		Email:       "applicant2@example.com",
		Password:    password2,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
		Type:        domain.Applicant,
	}
	employer2 = domain.DbUser{
		ID:          1,
		Email:       "employer2@example.com",
		Password:    password2,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
		Type:        domain.Employer,
	}
	updateUser1 = domain.UserUpdate{
		Email:       "user@example.com",
		Password:    password1,
		NewPassword: password2,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
	}
	updateUser2 = domain.UserUpdate{
		Email:       "user@example.com",
		Password:    password1,
		NewPassword: "",
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    &birthday,
		PhoneNumber: &phone_number,
		Location:    &location,
	}

	dbUserRows = []string{
		"id",
		"app_id",
		"emp_id",
		"email",
		"first_name",
		"last_name",
		"birthday",
		"phone_number",
		"location",
		"avatar_path",
	}
)

var testCheckUserCases = []struct {
	inputUser      domain.DbUser
	hashedPassword []byte
	salt           []byte
}{
	{
		inputUser:      applicant1,
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		inputUser:      applicant2,
		hashedPassword: hashedPassword2,
		salt:           salt2,
	},
	{
		inputUser:      employer1,
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		inputUser:      employer2,
		hashedPassword: hashedPassword2,
		salt:           salt2,
	},
}

func TestCheckUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckUserCases {
		checkPasswordByEmailRows := sqlmock.NewRows([]string{"actual_hash", "salt"}).
			AddRow(testCase.hashedPassword, testCase.salt)

		checkRoleRows := sqlmock.NewRows([]string{"is_employer"}).
			AddRow(true)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery(testHelper.SelectExistQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkRoleRows)

		actual := repo.CheckUser(ctxWithLogger, &testCase.inputUser)
		if actual != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestCheckUserIncorrectRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckUserCases {
		checkPasswordByEmailRows := sqlmock.NewRows([]string{"actual_hash", "salt"}).
			AddRow(testCase.hashedPassword, testCase.salt)

		checkRoleRows := sqlmock.NewRows([]string{"is_employer"}).
			AddRow(false)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery("SELECT EXISTS(.|\n)+").
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkRoleRows)

		actual := repo.CheckUser(ctxWithLogger, &testCase.inputUser)
		if actual != serverErrors.INCORRECT_ROLE {
			t.Errorf("got unexpected err: %s\nexpected: %s", err, serverErrors.INCORRECT_ROLE)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestCheckUserQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckUserCases {
		checkPasswordByEmailRows := sqlmock.NewRows([]string{"actual_hash", "salt"}).
			AddRow(testCase.hashedPassword, testCase.salt)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery("SELECT EXISTS(.|\n)+").
			WithArgs(testCase.inputUser.Email).
			WillReturnError(testHelper.ErrQuery)

		actual := repo.CheckUser(ctxWithLogger, &testCase.inputUser)
		if actual != testHelper.ErrQuery {
			t.Errorf("got unexpected err: %s\nexpected: %s", err, testHelper.ErrQuery)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}

	for _, testCase := range testCheckUserCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnError(testHelper.ErrQuery)

		actual := repo.CheckUser(ctxWithLogger, &testCase.inputUser)
		if actual != testHelper.ErrQuery {
			t.Errorf("got unexpected err: %s\nexpected: %s", err, testHelper.ErrQuery)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestCheckUserErrEntityNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckUserCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUser.Email).
			WillReturnError(sql.ErrNoRows)

		actual := repo.CheckUser(ctxWithLogger, &testCase.inputUser)
		if actual != ErrEntityNotFound {
			t.Errorf("got unexpected err: %s\nexpected: %s", err, testHelper.ErrQuery)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

var testCheckPasswordByIdCases = []struct {
	inputID             int
	passwordToCheck     string
	hashedPassword      []byte
	salt                []byte
	returningQueryError error
	expectedError       error
}{
	{
		inputID:             1,
		passwordToCheck:     password1,
		hashedPassword:      hashedPassword1,
		salt:                salt1,
		returningQueryError: sql.ErrNoRows,
		expectedError:       ErrEntityNotFound,
	},
	{
		inputID:             2,
		passwordToCheck:     password2,
		hashedPassword:      hashedPassword2,
		salt:                salt2,
		returningQueryError: testHelper.ErrQuery,
		expectedError:       testHelper.ErrQuery,
	},
}

func TestCheckPasswordByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckPasswordByIdCases {
		rows := sqlmock.NewRows([]string{"actual_hash", "salt"}).
			AddRow(testCase.hashedPassword, testCase.salt)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputID).
			WillReturnRows(rows)

		actual := repo.CheckPasswordById(ctxWithLogger, testCase.inputID, testCase.passwordToCheck)
		if actual != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestCheckPasswordByIdQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)
	for _, testCase := range testCheckPasswordByIdCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputID).
			WillReturnError(testCase.returningQueryError)

		actual := repo.CheckPasswordById(ctxWithLogger, testCase.inputID, testCase.passwordToCheck)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actual != testCase.expectedError {
			t.Errorf("expected query error, got: '%s'", actual)
			return
		}
	}
}

var testAddUserCases = []struct {
	user           domain.ApiUser
	hashedPassword []byte
	salt           []byte
}{
	{
		user:           *applicant1.ToAPI(nil, &appID),
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		user:           *applicant2.ToAPI(nil, &appID),
		hashedPassword: hashedPassword2,
		salt:           salt2,
	},
	{
		user:           *employer1.ToAPI(&empID, nil),
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		user:           *employer2.ToAPI(&empID, nil),
		hashedPassword: hashedPassword2,
		salt:           salt2,
	},
}

func TestAddUserSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	hasher := func(password string) (hash []byte, salt []byte, err error) {
		return hashedPassword1, salt1, nil
	}
	for _, testCase := range testAddUserCases {
		existsRows := sqlmock.NewRows([]string{"exists"}).
			AddRow(false)

		mock.ExpectBegin()

		mock.
			ExpectQuery(testHelper.SelectExistQuery).
			WithArgs(testCase.user.Email).
			WillReturnRows(existsRows)

		// insertion into user_profile table
		mock.
			ExpectQuery(testHelper.InsertQuery).
			WithArgs(
				testCase.user.Email,
				hashedPassword1,
				salt1,
				testCase.user.FirstName,
				testCase.user.LastName,
				testCase.user.Birthday,
				testCase.user.PhoneNumber,
				testCase.user.Location,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testCase.user.ID))

		// insertion into applicant/employer table
		if testCase.user.EmployerID == nil {
			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.user.ID).
				WillReturnResult(driver.RowsAffected(1))
		} else if testCase.user.ApplicantID == nil {
			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(
					testCase.user.ID,
					testCase.user.OrganizationName,
					testCase.user.OrganizationDescription,
				).
				WillReturnResult(driver.RowsAffected(1))
		}

		mock.ExpectCommit()

		actual := repo.AddUser(ctxWithLogger, &testCase.user, hasher)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if actual != nil {
			t.Errorf("unexpected err: %s", err)
			return
		}
	}
}

var testGetUserInfoCases = []struct {
	inputUserID int
	appID       *int
	empID       *int
	expected    domain.DbUser
}{
	{
		inputUserID: 1,
		appID:       &appID,
		empID:       nil,
		expected:    applicant1,
	},
	{
		inputUserID: 1,
		appID:       &appID,
		empID:       nil,
		expected:    applicant2,
	},
	{
		inputUserID: 1,
		appID:       nil,
		empID:       &empID,
		expected:    employer1,
	},
	{
		inputUserID: 1,
		appID:       nil,
		empID:       &empID,
		expected:    employer2,
	},
}

func TestGetUserInfoSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	for _, testCase := range testGetUserInfoCases {
		rows := sqlmock.NewRows(dbUserRows).
			AddRow(
				testCase.inputUserID,
				testCase.appID,
				testCase.empID,
				testCase.expected.Email,
				testCase.expected.FirstName,
				testCase.expected.LastName,
				testCase.expected.Birthday,
				testCase.expected.PhoneNumber,
				testCase.expected.Location,
				testCase.expected.AvatarPath,
			)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID).
			WillReturnRows(rows)

		mock.
			ExpectQuery(testHelper.SelectExistQuery).
			WithArgs(testCase.inputUserID).
			WillReturnRows(
				sqlmock.NewRows([]string{"bool"}).
					AddRow(true),
			)

		actualUser, actualAppID, actualEmpID, getErr := repo.GetUserInfo(ctxWithLogger, testCase.inputUserID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if getErr != nil {
			t.Errorf("unexpected err: %s", getErr)
			return
		}
		if actualUser.ID != testCase.expected.ID ||
			actualUser.Email != testCase.expected.Email ||
			actualUser.FirstName != testCase.expected.FirstName ||
			actualUser.LastName != testCase.expected.LastName ||
			actualUser.Birthday != testCase.expected.Birthday ||
			actualUser.PhoneNumber != testCase.expected.PhoneNumber ||
			actualUser.Location != testCase.expected.Location ||
			actualUser.AvatarPath != testCase.expected.AvatarPath {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actualUser))
			return
		}
		if actualAppID != testCase.appID {
			t.Errorf(testHelper.ErrNotEqual(testCase.appID, actualAppID))
			return
		}
		if actualEmpID != testCase.empID {
			t.Errorf(testHelper.ErrNotEqual(testCase.empID, actualEmpID))
			return
		}
	}
}

var testGetUserIdByEmailCases = []struct {
	inputEmail string
	expected   int
}{
	{
		inputEmail: applicant1.Email,
		expected:   applicant1.ID,
	},
	{
		inputEmail: applicant2.Email,
		expected:   applicant2.ID,
	},
	{
		inputEmail: employer1.Email,
		expected:   employer1.ID,
	},
	{
		inputEmail: employer2.Email,
		expected:   employer2.ID,
	},
}

func TestGetUserIdByEmailSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	for _, testCase := range testGetUserIdByEmailCases {
		rows := sqlmock.NewRows([]string{"user_id"}).
			AddRow(testCase.expected)

		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputEmail).
			WillReturnRows(rows)

		actual, getErr := repo.GetUserIdByEmail(ctxWithLogger, testCase.inputEmail)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if getErr != nil {
			t.Errorf("unexpected err: %s", getErr)
			return
		}
		if actual != testCase.expected {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
			return
		}
	}
}

var testGetRoleByIdCases = []struct {
	inputUserID int
	expected    domain.Role
}{
	{
		inputUserID: applicant1.ID,
		expected:    domain.Applicant,
	},
	{
		inputUserID: applicant2.ID,
		expected:    domain.Applicant,
	},
	{
		inputUserID: employer1.ID,
		expected:    domain.Employer,
	},
	{
		inputUserID: employer2.ID,
		expected:    domain.Employer,
	},
}

func TestGetRoleByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	for _, testCase := range testGetRoleByIdCases {
		if testCase.expected == domain.Applicant {
			appRows := sqlmock.NewRows([]string{"is_applicant"}).
				AddRow(true)

			mock.
				ExpectQuery(testHelper.SelectExistQuery).
				WithArgs(testCase.inputUserID).
				WillReturnRows(appRows)
		} else if testCase.expected == domain.Employer {
			appRows := sqlmock.NewRows([]string{"is_applicant"}).
				AddRow(false)
			empRows := sqlmock.NewRows([]string{"is_employer"}).
				AddRow(true)
			mock.
				ExpectQuery(testHelper.SelectExistQuery).
				WithArgs(testCase.inputUserID).
				WillReturnRows(appRows)
			mock.
				ExpectQuery(testHelper.SelectExistQuery).
				WithArgs(testCase.inputUserID).
				WillReturnRows(empRows)
		}

		actual, getErr := repo.GetRoleById(ctxWithLogger, testCase.inputUserID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if getErr != nil {
			t.Errorf("unexpected err: %s", getErr)
			return
		}
		if actual != testCase.expected {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
			return
		}
	}
}

var testUpdateUserInfoCases = []struct {
	inputUserID int
	inputUser   domain.UserUpdate
}{
	{
		inputUserID: 1,
		inputUser:   updateUser1,
	},
	{
		inputUserID: 2,
		inputUser:   updateUser2,
	},
	{
		inputUserID: 1,
		inputUser:   updateUser1,
	},
	{
		inputUserID: 2,
		inputUser:   updateUser2,
	},
}

func TestUpdateUserInfoSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	for _, testCase := range testUpdateUserInfoCases {
		if testCase.inputUser.NewPassword != "" {
			mock.
				ExpectExec(testHelper.UpdateQuery).
				WithArgs(
					testCase.inputUser.Email,
					hashedPassword2,
					salt2,
					testCase.inputUser.FirstName,
					testCase.inputUser.LastName,
					testCase.inputUser.Birthday,
					testCase.inputUser.PhoneNumber,
					testCase.inputUser.Location,
					testCase.inputUserID,
				).
				WillReturnResult(driver.RowsAffected(1))
		} else {
			mock.
				ExpectExec(testHelper.UpdateQuery).
				WithArgs(
					testCase.inputUser.Email,
					testCase.inputUser.FirstName,
					testCase.inputUser.LastName,
					testCase.inputUser.Birthday,
					testCase.inputUser.PhoneNumber,
					testCase.inputUser.Location,
					testCase.inputUserID,
				)
			// WillReturnResult(driver.RowsAffected(1))
		}

		repo.UpdateUserInfo(ctxWithLogger, testCase.inputUserID, &testCase.inputUser)
		// if err := mock.ExpectationsWereMet(); err != nil {
		// 	t.Errorf("there were unfulfilled expectations: %s", err)
		// 	return
		// }
		// if updErr != nil {
		// 	t.Errorf("unexpected err: %s", updErr)
		// 	return
		// }
	}
}

var testGetUserEmpIdCases = []struct {
	inputUserID int
	expected    int
}{
	{
		inputUserID: 1,
		expected:    1,
	},
	{
		inputUserID: 1,
		expected:    2,
	},
}

func TestGetUserEmpIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPsqlUserRepository(db)

	for _, testCase := range testGetUserEmpIdCases {
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.inputUserID).
			WillReturnRows(
				sqlmock.NewRows([]string{"emp_id"}).
					AddRow(testCase.expected),
			)

		actual, empErr := repo.GetUserEmpId(ctxWithLogger, testCase.inputUserID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if empErr != nil {
			t.Errorf("unexpected err: %s", empErr)
			return
		}
		if actual != testCase.expected {
			t.Errorf(testHelper.ErrNotEqual(testCase.expected, actual))
		}
	}
}
