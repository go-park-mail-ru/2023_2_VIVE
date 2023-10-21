package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
)

type ICVUsecase interface {
	GetCVById(sessionID string, cvID int) (*domain.CV, error)
	GetCVList(sessionID string) ([]domain.CV, error)
	AddNewCV(sessionID string, cv *domain.CV) (int, error)
	GetCVOfUserById(sessionID string, cvID int) (*domain.CV, error)
	UpdateCVOfUserById(sessionID string, cvID int) error
	DeleteCVOfUserById(sessionID string, cvID int) error
}

type CVUsecase struct {
	cvRepo      repository.ICVRepository
	sessionRepo repository.ISessionRepository
	userRepo    repository.IUserRepository
}

func NewCVUsecase(cvRepository repository.ICVRepository, sessionRepository repository.ISessionRepository, userRepository repository.IUserRepository) ICVUsecase {
	return &CVUsecase{
		cvRepo:      cvRepository,
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (cvUsecase *CVUsecase) GetCVById(sessionID string, cvID int) (*domain.CV, error) {
	return nil, nil
}

func (cvUsecase *CVUsecase) GetCVList(sessionID string) ([]domain.CV, error) {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, err
	}

	userRole, err := cvUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return nil, err
	} else if userRole != domain.Applicant {
		return nil, INAPPROPRIATE_ROLE
	}

	cvs, err := cvUsecase.cvRepo.GetByUserId(userID)
	if err != nil {
		return nil, err
	}

	return cvs, nil
}

func (cvUsecase *CVUsecase) AddNewCV(sessionID string, cv *domain.CV) (int, error) {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

	cv.UserID = userID

	cvID, addErr := cvUsecase.cvRepo.AddCV(cv)
	if addErr != nil {
		return 0, addErr
	}

	return cvID, nil
}

func (cvUsecase *CVUsecase) GetCVOfUserById(sessionID string, cvID int) (*domain.CV, error) {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, err
	}

	cv, err := cvUsecase.cvRepo.GetOneOfUsersCV(userID, cvID)
	if err != nil {
		return nil, err
	}

	return cv, nil
}

func (cvUsecase *CVUsecase) UpdateCVOfUserById(sessionID string, cvID int) error {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return err
	}

	updStatus := cvUsecase.cvRepo.UpdateOneOfUsersCV(userID, cvID)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (cvUsecase *CVUsecase) DeleteCVOfUserById(sessionID string, cvID int) error {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return err
	}

	delStatus := cvUsecase.cvRepo.DeleteOneOfUsersCV(userID, cvID)
	if delStatus != nil {
		return delStatus
	}

	return nil
}
