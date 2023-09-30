package modelHandlers

import (
	"crypto/sha256"
	"encoding/hex"
	"models/models"
	"models/serverErrors"
	"net/http"
)

func CheckPassword(user models.User) error {
	actualPass, ok := models.Users.Load(user.Email)
	if !ok {
		return serverErrors.NO_DATA_FOUND
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	if hex.EncodeToString(hashedPass) != hex.EncodeToString(actualPass.([]byte)) {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	return nil
}

func CheckUser(user models.User) error {
	if len(user.Email) == 0 || len(user.Password) == 0 {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	passwordStatus := CheckPassword(user)
	if passwordStatus != nil {
		return passwordStatus
	}

	return nil
}

func AddUser(user models.User) error {
	_, exist := models.Users.Load(user.Email)

	if exist {
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	models.Users.Store(user.Email, hashedPass)

	return nil
}

func GetUserInfo(cookie *http.Cookie) models.User {
	uniqueID := cookie.Value

	user, _ := models.Sessions.Load(uniqueID)
	return user.(models.User)
}
