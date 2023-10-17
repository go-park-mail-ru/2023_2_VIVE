package repository

import (
	"HnH/internal/domain"
	"HnH/internal/repository/mock"
	"HnH/pkg/serverErrors"

	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"unicode"

	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier()

func CheckPassword(user *domain.User) error {
	actualUserIndex, ok := mock.UserDB.EmailToUser.Load(user.Email)

	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	actualUser := mock.UserDB.UsersList[actualUserIndex.(int)]

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	if hex.EncodeToString(hashedPass) != actualUser.Password {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	return nil
}

func CheckRole(user *domain.User) error {
	actualUserIndex, ok := mock.UserDB.EmailToUser.Load(user.Email)

	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	actualUser := mock.UserDB.UsersList[actualUserIndex.(int)]

	if user.Type != actualUser.Type {
		return serverErrors.INCORRECT_ROLE
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 128 {
		return serverErrors.INVALID_PASSWORD
	}

	hasDigit := false
	hasCapital := false
	hasSpecialChar := false

	for _, char := range password {
		if !hasDigit {
			hasDigit = unicode.IsDigit(char)
		}

		if !hasCapital {
			hasCapital = unicode.Is(unicode.Latin, char) && unicode.IsUpper(char)
		}

		if !hasSpecialChar {
			for _, specCh := range domain.SpecialChars {
				hasSpecialChar = (char == specCh)

				if hasSpecialChar {
					break
				}
			}
		}
	}

	if !(hasDigit && hasCapital && hasSpecialChar) {
		return serverErrors.INVALID_PASSWORD
	}

	return nil
}

func CheckUser(user *domain.User) error {
	if len(user.Email) == 0 || len(user.Password) == 0 {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	passwordStatus := CheckPassword(user)
	if passwordStatus != nil {
		return passwordStatus
	}

	roleStatus := CheckRole(user)
	if roleStatus != nil {
		return roleStatus
	}

	return nil
}

func AddUser(user *domain.User) error {
	_, exist := mock.UserDB.EmailToUser.Load(user.Email)

	if exist {
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	}

	ret, err := verifier.Verify(user.Email)
	if err != nil {
		return err
	} else if !ret.Syntax.Valid {
		return serverErrors.INVALID_EMAIL
	}

	if len(user.Email) == 0 || len(user.Password) == 0 {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	validPassStatus := ValidatePassword(user.Password)
	if validPassStatus != nil {
		return validPassStatus
	}

	if !user.Type.IsRole() {
		return serverErrors.INVALID_ROLE
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	mock.UserDB.Mu.Lock()

	defer mock.UserDB.Mu.Unlock()

	mock.UserDB.CurrentID++
	user.ID = mock.UserDB.CurrentID
	user.Password = hex.EncodeToString(hashedPass)

	mock.UserDB.UsersList = append(mock.UserDB.UsersList, user)

	mock.UserDB.EmailToUser.Store(user.Email, len(mock.UserDB.UsersList)-1)
	mock.UserDB.IdToUser.Store(user.ID, len(mock.UserDB.UsersList)-1)

	return nil
}

func GetUserInfo(cookie *http.Cookie) (*domain.User, error) {
	uniqueID := cookie.Value

	userID, exist := mock.SessionDB.SessionsList.Load(uniqueID)
	if !exist {
		return nil, serverErrors.AUTH_REQUIRED
	}

	userIndex, exist := mock.UserDB.IdToUser.Load(userID.(int))
	if !exist {
		return nil, serverErrors.NO_DATA_FOUND
	}

	user := mock.UserDB.UsersList[userIndex.(int)]

	user.Password = ""

	return user, nil
}
