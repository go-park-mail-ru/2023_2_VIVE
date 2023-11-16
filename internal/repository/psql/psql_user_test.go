package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"
	"HnH/pkg/nullTypes"
	"HnH/pkg/serverErrors"
	"database/sql"
	"database/sql/driver"

	// "HnH/pkg/serverErrors"
	"HnH/pkg/testHelper"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	password1                 = "password_number_1"
	hashedPassword1, salt1, _ = authUtils.GenerateHash(password1)
	applicant1                = domain.DbUser{
		ID:          1,
		Email:       "applicant1@example.com",
		Password:    password1,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    nullTypes.NewNullString("2000-01-01", true),
		PhoneNumber: nullTypes.NewNullString("+71111111111", true),
		Location:    nullTypes.NewNullString("Moscow", true),
		Type:        domain.Applicant,
	}
	employer1 = domain.DbUser{
		ID:          1,
		Email:       "employer1@example.com",
		Password:    password1,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    nullTypes.NewNullString("2000-01-01", true),
		PhoneNumber: nullTypes.NewNullString("+71111111111", true),
		Location:    nullTypes.NewNullString("Moscow", true),
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
		Birthday:    nullTypes.NewNullString("2000-01-01", true),
		PhoneNumber: nullTypes.NewNullString("+71111111111", true),
		Location:    nullTypes.NewNullString("Moscow", true),
		Type:        domain.Applicant,
	}
	employer2 = domain.DbUser{
		ID:          1,
		Email:       "employer2@example.com",
		Password:    password2,
		FirstName:   "Ivan",
		LastName:    "Ivanov",
		Birthday:    nullTypes.NewNullString("2000-01-01", true),
		PhoneNumber: nullTypes.NewNullString("+71111111111", true),
		Location:    nullTypes.NewNullString("Moscow", true),
		Type:        domain.Employer,
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery(testHelper.SELECT_EXISTS_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkRoleRows)

		actual := repo.CheckUser(&testCase.inputUser)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery("SELECT EXISTS(.|\n)+").
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkRoleRows)

		actual := repo.CheckUser(&testCase.inputUser)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnRows(checkPasswordByEmailRows)

		mock.
			ExpectQuery("SELECT EXISTS(.|\n)+").
			WithArgs(testCase.inputUser.Email).
			WillReturnError(testHelper.ErrQuery)

		actual := repo.CheckUser(&testCase.inputUser)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnError(testHelper.ErrQuery)

		actual := repo.CheckUser(&testCase.inputUser)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputUser.Email).
			WillReturnError(sql.ErrNoRows)

		actual := repo.CheckUser(&testCase.inputUser)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputID).
			WillReturnRows(rows)

		actual := repo.CheckPasswordById(testCase.inputID, testCase.passwordToCheck)
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
			ExpectQuery(testHelper.SELECT_QUERY).
			WithArgs(testCase.inputID).
			WillReturnError(testCase.returningQueryError)

		actual := repo.CheckPasswordById(testCase.inputID, testCase.passwordToCheck)
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
	user           domain.DbUser
	hashedPassword []byte
	salt           []byte
}{
	{
		user:           applicant1,
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		user:           applicant2,
		hashedPassword: hashedPassword2,
		salt:           salt2,
	},
	{
		user:           employer1,
		hashedPassword: hashedPassword1,
		salt:           salt1,
	},
	{
		user:           employer2,
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

		mock.
			ExpectQuery(testHelper.SELECT_EXISTS_QUERY).
			WithArgs(testCase.user.Email).
			WillReturnRows(existsRows)

		// insertion into user_profile table
		mock.
			ExpectQuery(testHelper.INSERT_QUERY).
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
		mock.
			ExpectExec(testHelper.INSERT_QUERY).
			WithArgs(testCase.user.ID).
			WillReturnResult(driver.RowsAffected(1))

		actual := repo.AddUser(&testCase.user, hasher)
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
