package modelHandlers

import (
	"crypto/sha256"
	"encoding/hex"
	"models/models"
	"models/serverErrors"
	"net/http"
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

	if len(user.Email) == 0 || len(user.Password) == 0 {
		return serverErrors.INCORRECT_CREDENTIALS
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

	return user
}
