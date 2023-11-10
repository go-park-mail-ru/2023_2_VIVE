package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/serverErrors"
)

type IResponseUsecase interface {
	RespondToVacancy(sessionID string, vacancyID, cvID int) error
	GetApplicantsList(sessionID string, vacancyID int) ([]domain.ApplicantInfo, error)
}

type ResponseUsecase struct {
	responseRepo psql.IResponseRepository
	sessionRepo  redisRepo.ISessionRepository
	userRepo     psql.IUserRepository
	vacancyRepo  psql.IVacancyRepository
	cvRepo       psql.ICVRepository
}

func NewResponseUsecase(respondRepository psql.IResponseRepository,
	sessionRepository redisRepo.ISessionRepository,
	userRepository psql.IUserRepository,
	vacancyRepository psql.IVacancyRepository,
	cvRepository psql.ICVRepository) IResponseUsecase {
	return &ResponseUsecase{
		responseRepo: respondRepository,
		sessionRepo:  sessionRepository,
		userRepo:     userRepository,
		vacancyRepo:  vacancyRepository,
		cvRepo:       cvRepository,
	}
}

func (responseUsecase *ResponseUsecase) validateSessionAndGetUserId(sessionID string) (int, error) {
	validStatus := responseUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := responseUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (responseUsecase *ResponseUsecase) RespondToVacancy(sessionID string, vacancyID, cvID int) error {
	userID, validStatus := responseUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
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
	userID, validStatus := responseUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return nil, validStatus
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

	// vacancy, err := responseUsecase.vacancyRepo.GetVacancy(vacancyID)
	// if err != nil {
	// 	return nil, err
	// }

	orgID, _ := responseUsecase.vacancyRepo.GetOrgId(vacancyID)
	if orgID != userOrgID {
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

func (responseUsecase *ResponseUsecase) makeSummary(CVs []domain.DbCV) []domain.ApplicantInfo {
	infoToReturn := make([]domain.ApplicantInfo, 0, len(CVs))

	for _, cv := range CVs {
		info := domain.ApplicantInfo{
			CVid:      cv.ID,
			// FirstName: cv.FirstName,
			// LastName:  cv.LastName,
			// Skills:    cv.Skills,	// FIXME: remove this field
		}

		infoToReturn = append(infoToReturn, info)
	}

	return infoToReturn
}
