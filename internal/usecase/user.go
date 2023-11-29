package usecase

import (
	"HnH/configs"
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"context"
	"io/ioutil"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	SignUp(ctx context.Context, user *domain.ApiUser, expiryUnixSeconds int64) (string, error)
	GetInfo(ctx context.Context, sessionID string) (*domain.ApiUser, error)
	UpdateInfo(ctx context.Context, sessionID string, user *domain.UserUpdate) error
	UploadAvatar(ctx context.Context, sessionID, path string) error
	GetAvatar(ctx context.Context, sessionID string) ([]byte, error)
}

type UserUsecase struct {
	userRepo    psql.IUserRepository
	sessionRepo redisRepo.ISessionRepository
}

func NewUserUsecase(
	userRepository psql.IUserRepository,
	sessionRepository redisRepo.ISessionRepository,
) IUserUsecase {
	return &UserUsecase{
		userRepo:    userRepository,
		sessionRepo: sessionRepository,
	}
}

func (userUsecase *UserUsecase) validateSessionAndGetUserId(ctx context.Context, sessionID string) (int, error) {
	validStatus := userUsecase.sessionRepo.ValidateSession(ctx, sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := userUsecase.sessionRepo.GetUserIdBySession(ctx, sessionID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (userUsecase *UserUsecase) SignUp(ctx context.Context, user *domain.ApiUser, expiryUnixSeconds int64) (string, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("validating email")
	validEmailStatus := authUtils.ValidateEmail(user.Email)
	if validEmailStatus != nil {
		contextLogger.Info("validating email failed")
		return "", validEmailStatus
	}

	contextLogger.Info("validating password")
	validPassStatus := authUtils.ValidatePassword(user.Password)
	if validPassStatus != nil {
		contextLogger.Info("validating password failed")
		return "", validPassStatus
	}

	contextLogger.Info("validating role")
	if !user.Type.IsRole() {
		contextLogger.Info("validating role failed")
		return "", serverErrors.INVALID_ROLE
	}
	// fmt.Printf("before add user to db\n")
	// if user.Type == domain.Employer {
	// 	organization := domain.DbOrganization{
	// 		Name:        user.OrganizationName,
	// 		Description: "описание организации", // TODO: изменить описание организации по-умолчанию
	// 		Location:    user.Location,
	// 	}
	// 	contextLogger.Info("adding organization for employer")
	// 	_, addOrgErr := userUsecase.orgRepo.AddOrganization(ctx, &organization)
	// 	if addOrgErr != nil {
	// 		return "", addOrgErr
	// 	}
	// }

	contextLogger.Info("adding user")
	addStatus := userUsecase.userRepo.AddUser(ctx, user, authUtils.GenerateHash)
	if addStatus != nil {
		return "", addStatus
	}
	contextLogger.Info("user added")
	// fmt.Printf("after add user to db\n")

	contextLogger.Info("getting user by email")
	userID, err := userUsecase.userRepo.GetUserIdByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()

	contextLogger.Info("adding session")
	addErr := userUsecase.sessionRepo.AddSession(ctx, sessionID, userID, expiryUnixSeconds)
	if addErr != nil {
		return "", addErr
	}

	return sessionID, nil
}

func (userUsecase *UserUsecase) GetInfo(ctx context.Context, sessionID string) (*domain.ApiUser, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID, validStatus := userUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	contextLogger.Info("getting user information")
	user, appID, empID, getErr := userUsecase.userRepo.GetUserInfo(ctx, userID)
	if getErr != nil {
		return nil, getErr
	}
	// fmt.Printf("")

	apiUser := user.ToAPI(empID, appID)

	return apiUser, nil
}

func (userUsecase *UserUsecase) UpdateInfo(ctx context.Context, sessionID string, user *domain.UserUpdate) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID, validStatus := userUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return validStatus
	}

	validPassStatus := userUsecase.userRepo.CheckPasswordById(ctx, userID, user.Password)
	if validPassStatus != nil {
		return validPassStatus
	}
	contextLogger.Info("updating user info")
	updStatus := userUsecase.userRepo.UpdateUserInfo(ctx, userID, user)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (userUsecase *UserUsecase) UploadAvatar(ctx context.Context, sessionID, path string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID, validStatus := userUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return validStatus
	}

	contextLogger.Info("uploading avatar")
	err := userUsecase.userRepo.UploadAvatarByUserID(ctx, userID, path)
	if err != nil {
		return err
	}

	return nil
}

func (userUsecase *UserUsecase) GetAvatar(ctx context.Context, sessionID string) ([]byte, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	userID, validStatus := userUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	contextLogger.Info("getting user's avatar")
	path, err := userUsecase.userRepo.GetAvatarByUserID(ctx, userID)

	if path == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(configs.CURRENT_DIR + path)
	if err != nil {
		return nil, ErrReadAvatar
	}

	return fileBytes, nil
}
