package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository"
	"HnH/pkg/serverErrors"
)

type IResponseUsecase interface {
	RespondToVacancy(sessionID string, vacancyID, cvID int) error
	GetApplicantsList(sessionID string, vacancyID int) ([]domain.ApplicantInfo, error)
}

type ResponseUsecase struct {
	responseRepo repository.IResponseRepository
	sessionRepo  repository.ISessionRepository
	userRepo     repository.IUserRepository
	vacancyRepo  repository.IVacancyRepository
	cvRepo       repository.ICVRepository
}

func NewResponseUsecase(respondRepository repository.IResponseRepository,
	sessionRepository repository.ISessionRepository,
	userRepository repository.IUserRepository,
	vacancyRepository repository.IVacancyRepository,
	cvRepository repository.ICVRepository) IResponseUsecase {
	return &ResponseUsecase{
		responseRepo: respondRepository,
		sessionRepo:  sessionRepository,
		userRepo:     userRepository,
		vacancyRepo:  vacancyRepository,
		cvRepo:       cvRepository,
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
	} else if userRole != domain.Applicant {
		return INAPPROPRIATE_ROLE
	}

	responseUsecase.responseRepo.RespondToVacancy(vacancyID, cvID)
	return nil
}

func (responseUsecase *ResponseUsecase) GetApplicantsList(sessionID string, vacancyID int) ([]domain.ApplicantInfo, error) {
	validStatus := responseUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	userID, err := responseUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return nil, err
	}

	userRole, err := responseUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return nil, err
	} else if userRole != domain.Employer {
		return nil, INAPPROPRIATE_ROLE
	}

	userOrgID, err := responseUsecase.userRepo.GetUserOrgId(userID)
	if err != nil {
		return nil, err
	}

	vacancy, err := responseUsecase.vacancyRepo.GetVacancy(vacancyID)
	if err != nil {
		return nil, err
	}

	if vacancy.CompanyID != userOrgID {
		return nil, serverErrors.FORBIDDEN
	}

	cvIDs, err := responseUsecase.responseRepo.GetAttachedCVs(vacancyID)
	if err != nil {
		return nil, err
	}

	CVs, err := responseUsecase.cvRepo.GetCVsByIds(cvIDs)
	if err != nil {
		return nil, err
	}

	return responseUsecase.makeSummary(CVs), nil
}

func (responseUsecase *ResponseUsecase) makeSummary(CVs []domain.CV) []domain.ApplicantInfo {
	infoToReturn := make([]domain.ApplicantInfo, 0, len(CVs))

	for _, cv := range CVs {
		info := domain.ApplicantInfo{
			CVid:      cv.ID,
			FirstName: cv.FirstName,
			LastName:  cv.LastName,
			Skills:    cv.Skills,
		}

		infoToReturn = append(infoToReturn, info)
	}

	return infoToReturn
}
