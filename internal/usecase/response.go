package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/serverErrors"
	"context"
)

type IResponseUsecase interface {
	RespondToVacancy(ctx context.Context, sessionID string, vacancyID, cvID int) error
	GetApplicantsList(ctx context.Context, sessionID string, vacancyID int) ([]domain.ApiApplicant, error)
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

func (responseUsecase *ResponseUsecase) validateSessionAndGetUserId(ctx context.Context, sessionID string) (int, error) {
	validStatus := responseUsecase.sessionRepo.ValidateSession(ctx, sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := responseUsecase.sessionRepo.GetUserIdBySession(ctx, sessionID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (responseUsecase *ResponseUsecase) RespondToVacancy(ctx context.Context, sessionID string, vacancyID, cvID int) error {
	userID, validStatus := responseUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return validStatus
	}

	userRole, err := responseUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return err
	} else if userRole != domain.Applicant {
		return ErrInapropriateRole
	}

	respErr := responseUsecase.responseRepo.RespondToVacancy(ctx, vacancyID, cvID)
	if respErr != nil {
		return respErr
	}
	return nil
}

func (responseUsecase *ResponseUsecase) GetApplicantsList(ctx context.Context, sessionID string, vacancyID int) ([]domain.ApiApplicant, error) {
	userID, validStatus := responseUsecase.validateSessionAndGetUserId(ctx, sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	userRole, err := responseUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return nil, err
	} else if userRole != domain.Employer {
		return nil, ErrInapropriateRole
	}

	userEmpID, err := responseUsecase.userRepo.GetUserEmpId(ctx, userID)
	if err != nil {
		return nil, err
	}

	// vacancy, err := responseUsecase.vacancyRepo.GetVacancy(vacancyID)
	// if err != nil {
	// 	return nil, err
	// }

	empID, _ := responseUsecase.vacancyRepo.GetEmpId(ctx, vacancyID)
	if empID != userEmpID {
		return nil, serverErrors.FORBIDDEN
	}

	cvIDs, err := responseUsecase.responseRepo.GetAttachedCVs(ctx, vacancyID)
	if err != nil {
		return nil, err
	}

	CVs, _, _, err := responseUsecase.cvRepo.GetCVsByIds(ctx, cvIDs)
	if err != nil && err != psql.ErrEntityNotFound {
		return nil, err
	}

	return responseUsecase.makeSummary(CVs), nil
}

func (responseUsecase *ResponseUsecase) makeSummary(CVs []domain.DbCV) []domain.ApiApplicant {
	infoToReturn := make([]domain.ApiApplicant, 0, len(CVs))

	for _, cv := range CVs {
		info := domain.ApiApplicant{
			CVid: cv.ID,
			// FirstName: cv.FirstName,
			// LastName:  cv.LastName,
			// Skills:    cv.Skills,	// FIXME: remove this field
		}

		infoToReturn = append(infoToReturn, info)
	}

	return infoToReturn
}
