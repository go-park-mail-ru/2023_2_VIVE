package usecase

import (
	"HnH/configs"
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	SignUp(ctx context.Context, user *domain.ApiUser, expiryUnixSeconds int64) (string, error)
	GetInfo(ctx context.Context) (*domain.ApiUser, error)
	UpdateInfo(ctx context.Context, user *domain.UserUpdate) error
	UploadAvatar(ctx context.Context, uploadedData multipart.File, header *multipart.FileHeader) error
	GetUserAvatar(ctx context.Context) ([]byte, error)
	GetImage(ctx context.Context, imageID int) ([]byte, error)
}

type UserUsecase struct {
	userRepo    psql.IUserRepository
	sessionRepo grpc.IAuthRepository
}

func NewUserUsecase(
	userRepository psql.IUserRepository,
	sessionRepository grpc.IAuthRepository,
) IUserUsecase {
	return &UserUsecase{
		userRepo:    userRepository,
		sessionRepo: sessionRepository,
	}
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

func (userUsecase *UserUsecase) GetInfo(ctx context.Context) (*domain.ApiUser, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID := contextUtils.GetUserIDFromCtx(ctx)

	contextLogger.Info("getting user information")
	user, appID, empID, getErr := userUsecase.userRepo.GetUserInfo(ctx, userID)
	if getErr != nil {
		return nil, getErr
	}
	// fmt.Printf("")

	apiUser := user.ToAPI(empID, appID)

	return apiUser, nil
}

func (userUsecase *UserUsecase) UpdateInfo(ctx context.Context, user *domain.UserUpdate) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID := contextUtils.GetUserIDFromCtx(ctx)

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

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}

	return false
}

func (userUsecase *UserUsecase) avatarExists(ctx context.Context, userID int) (bool, error) {
	path, err := userUsecase.userRepo.GetAvatarByUserID(ctx, userID)
	if err != nil {
		return false, err
	} else if path == "" && err == nil {
		return false, nil
	}

	return true, nil
}

func (userUsecase *UserUsecase) UploadAvatar(ctx context.Context, uploadedData multipart.File, header *multipart.FileHeader) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID := contextUtils.GetUserIDFromCtx(ctx)

	contextLogger.Info("uploading avatar")

	if header.Size > 2*1024*1024 {
		return BadAvatarSize
	}

	mimeType := header.Header.Get("Content-Type")
	ext, err := mime.ExtensionsByType(mimeType)
	if ext == nil {
		return BadAvatarType
	} else if err != nil {
		return err
	}

	if contains(ext, ".jpeg") {
		header.Filename = strconv.Itoa(userID) + ".jpeg"
	} else if contains(ext, ".png") {
		header.Filename = strconv.Itoa(userID) + ".png"
	} else if contains(ext, ".gif") {
		header.Filename = strconv.Itoa(userID) + ".gif"
	} else {
		return BadAvatarType
	}

	avaExists, err := userUsecase.avatarExists(ctx, userID)
	if err != nil {
		return err
	}

	if avaExists {
		path, err := userUsecase.userRepo.GetAvatarByUserID(ctx, userID)
		if err != nil {
			return err
		}

		delErr := os.Remove(configs.CURRENT_DIR + path)
		if delErr != nil {
			return err
		}
	}

	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())

	dirToSave := configs.UPLOADS_DIR + year + "/" + month + "/" + day
	err = os.MkdirAll(configs.CURRENT_DIR+dirToSave, 0777)
	if err != nil {
		return err
	}

	filePath := dirToSave + "/" + header.Filename
	f, err := os.OpenFile(configs.CURRENT_DIR+filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	contextLogger.Infof("file path: %s", filePath)
	fmt.Println(filePath)

	_, copyErr := io.Copy(f, uploadedData)
	if copyErr != nil {
		return copyErr
	}

	syncErr := f.Sync()
	if syncErr != nil {
		return syncErr
	}
	closeErr := f.Close()
	if closeErr != nil {
		return closeErr
	}

	err = userUsecase.userRepo.UploadAvatarByUserID(ctx, userID, filePath)
	if err != nil {
		return err
	}

	return nil
}

func (userUsecase *UserUsecase) GetUserAvatar(ctx context.Context) ([]byte, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	userID := contextUtils.GetUserIDFromCtx(ctx)

	contextLogger.Info("getting user's avatar")
	path, err := userUsecase.userRepo.GetAvatarByUserID(ctx, userID)

	if path == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(configs.CURRENT_DIR + path)
	if err != nil {
		return nil, ErrReadAvatar
	}

	return fileBytes, nil
}

func (userUsecase *UserUsecase) GetImage(ctx context.Context, imageID int) ([]byte, error) {
	path, err := userUsecase.userRepo.GetAvatarByUserID(ctx, imageID)

	if path == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(configs.CURRENT_DIR + path)
	if err != nil {
		return nil, ErrReadAvatar
	}

	return fileBytes, nil
}
