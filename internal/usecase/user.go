package usecase

import (
	"HnH/configs"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"
	"io/ioutil"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	SignUp(user *domain.ApiUser, expiryUnixSeconds int64) (string, error)
	GetInfo(sessionID string) (*domain.ApiUser, error)
	UpdateInfo(sessionID string, user *domain.UserUpdate) error
	UploadAvatar(sessionID, path string) error
	GetAvatar(sessionID string) ([]byte, error)
}

type UserUsecase struct {
	userRepo    psql.IUserRepository
	sessionRepo redisRepo.ISessionRepository
	orgRepo     psql.IOrganizationRepository
}

func NewUserUsecase(
	userRepository psql.IUserRepository,
	sessionRepository redisRepo.ISessionRepository,
	orgRepository psql.IOrganizationRepository,
) IUserUsecase {
	return &UserUsecase{
		userRepo:    userRepository,
		sessionRepo: sessionRepository,
		orgRepo:     orgRepository,
	}
}

func (userUsecase *UserUsecase) validateSessionAndGetUserId(sessionID string) (int, error) {
	validStatus := userUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := userUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (userUsecase *UserUsecase) SignUp(user *domain.ApiUser, expiryUnixSeconds int64) (string, error) {
	validEmailStatus := authUtils.ValidateEmail(user.Email)
	if validEmailStatus != nil {
		return "", validEmailStatus
	}

	validPassStatus := authUtils.ValidatePassword(user.Password)
	if validPassStatus != nil {
		return "", validPassStatus
	}

	if !user.Type.IsRole() {
		return "", serverErrors.INVALID_ROLE
	}
	// fmt.Printf("before add user to db\n")
	if user.Type == domain.Employer {
		organization := domain.DbOrganization{
			Name:        user.OrganizationName,
			Description: "описание организации", // TODO: изменить описание организации по-умолчанию
			Location:    user.Location,
		}
		_, addOrgErr := userUsecase.orgRepo.AddOrganization(&organization)
		if addOrgErr != nil {
			return "", addOrgErr
		}
	}
	addStatus := userUsecase.userRepo.AddUser(user, authUtils.GenerateHash)
	if addStatus != nil {
		return "", addStatus
	}
	// fmt.Printf("after add user to db\n")

	userID, err := userUsecase.userRepo.GetUserIdByEmail(user.Email)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()

	addErr := userUsecase.sessionRepo.AddSession(sessionID, userID, expiryUnixSeconds)
	if addErr != nil {
		return "", addErr
	}

	return sessionID, nil
}

func (userUsecase *UserUsecase) GetInfo(sessionID string) (*domain.ApiUser, error) {
	userID, validStatus := userUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	user, appID, empID, getErr := userUsecase.userRepo.GetUserInfo(userID)
	if getErr != nil {
		return nil, getErr
	}
	// fmt.Printf("")

	apiUser := user.ToAPI(empID, appID)

	return apiUser, nil
}

func (userUsecase *UserUsecase) UpdateInfo(sessionID string, user *domain.UserUpdate) error {
	userID, validStatus := userUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	validPassStatus := userUsecase.userRepo.CheckPasswordById(userID, user.Password)
	if validPassStatus != nil {
		return validPassStatus
	}

	updStatus := userUsecase.userRepo.UpdateUserInfo(userID, user)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (userUsecase *UserUsecase) UploadAvatar(sessionID, path string) error {
	userID, validStatus := userUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	err := userUsecase.userRepo.UploadAvatarByUserID(userID, path)
	if err != nil {
		return err
	}

	return nil
}

func (userUsecase *UserUsecase) GetAvatar(sessionID string) ([]byte, error) {
	userID, validStatus := userUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	path, err := userUsecase.userRepo.GetAvatarByUserID(userID)

	if path == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(configs.CURRENT_DIR + path)
	if err != nil {
		return nil, CAN_NOT_READ_AVATAR
	}

	return fileBytes, nil
}
