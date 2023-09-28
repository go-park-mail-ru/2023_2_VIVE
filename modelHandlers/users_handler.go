package modelHandlers

import (
	"crypto/sha256"
	"encoding/hex"
	"models/errors"
	"models/models"
	"net/http"
)

func CheckPassword(user models.User) error {
	actualPass, ok := models.Users.Load(user.Email)
	if !ok {
		return errors.NO_DATA_FOUND
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPass := hasher.Sum(nil)

	if hex.EncodeToString(hashedPass) != hex.EncodeToString(actualPass.([]byte)) {
		return errors.INCORRECT_CREDENTIALS
	}

	return nil
}

func AddUser(user models.User) error {
	_, exist := models.Users.Load(user.Email)

	if exist {
		return errors.ACCOUNT_ALREADY_EXISTS
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
