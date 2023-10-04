package modelHandlers

import (
	"HnH/models"
	"HnH/serverErrors"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/mail"
	"unicode"
)

func CheckPassword(user *models.User) error {
	actualUserIndex, ok := models.EmailToUser.Load(user.Email)

	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	actualUser := models.UserDB.UsersList[actualUserIndex.(int)]

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	if hex.EncodeToString(hashedPass) != actualUser.Password {
		return serverErrors.INCORRECT_CREDENTIALS
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
			for _, specCh := range models.SpecialChars {
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

func CheckUser(user *models.User) error {
	if len(user.Email) == 0 || len(user.Password) == 0 {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	passwordStatus := CheckPassword(user)
	if passwordStatus != nil {
		return passwordStatus
	}

	return nil
}

func AddUser(user *models.User) error {
	_, exist := models.EmailToUser.Load(user.Email)

	if exist {
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
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

	models.UserDB.Mu.Lock()

	defer models.UserDB.Mu.Unlock()

	models.UserDB.CurrentID++
	user.ID = models.UserDB.CurrentID
	user.Password = hex.EncodeToString(hashedPass)

	models.UserDB.UsersList = append(models.UserDB.UsersList, user)

	models.EmailToUser.Store(user.Email, len(models.UserDB.UsersList)-1)
	models.IdToUser.Store(user.ID, len(models.UserDB.UsersList)-1)

	return nil
}

func GetUserInfo(cookie *http.Cookie) *models.User {
	uniqueID := cookie.Value

	userID, _ := models.Sessions.Load(uniqueID)
	userIndex, _ := models.IdToUser.Load(userID.(int))
	user := models.UserDB.UsersList[userIndex.(int)]

	user.Password = ""

	return user
}
