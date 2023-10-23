package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
)

type IResponseUsecase interface {
	RespondToVacancy(sessionID string, vacancyID, cvID int) error
}

type ResponseUsecase struct {
	responseRepo repository.IResponseRepository
	sessionRepo  repository.ISessionRepository
	userRepo     repository.IUserRepository
}

func NewResponseUsecase(respondRepository repository.IResponseRepository,
	sessionRepository repository.ISessionRepository,
	userRepository repository.IUserRepository) IResponseUsecase {
	return &ResponseUsecase{
		responseRepo: respondRepository,
		sessionRepo:  sessionRepository,
		userRepo:     userRepository,
	}
}

func (responseUsecase *ResponseUsecase) RespondToVacancy(sessionID string, vacancyID, cvID int) error {
	validStatus := responseUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return validStatus
	}

	userID, err := responseUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return err
	}

	userRole, err := responseUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return err
	}

	if userRole != domain.Applicant {
		return INAPPROPRIATE_ROLE
	}

	responseUsecase.responseRepo.RespondToVacancy(vacancyID, cvID)
	return nil
}
