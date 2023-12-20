package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/pkg/contextUtils"
	"HnH/pkg/notificationMessages"
	"HnH/pkg/serverErrors"
	notificationsPB "HnH/services/notifications/api/proto"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type IResponseUsecase interface {
	RespondToVacancy(ctx context.Context, vacancyID, cvID int) error
	GetApplicantsList(ctx context.Context, vacancyID int) ([]domain.ApiApplicant, error)
	GetUserResponses(ctx context.Context, userID int) ([]domain.ApiResponse, error)
}

type ResponseUsecase struct {
	responseRepo      psql.IResponseRepository
	sessionRepo       grpc.IAuthRepository
	userRepo          psql.IUserRepository
	vacancyRepo       psql.IVacancyRepository
	cvRepo            psql.ICVRepository
	notificationsRepo grpc.INotificationRepository
}

func NewResponseUsecase(respondRepository psql.IResponseRepository,
	sessionRepository grpc.IAuthRepository,
	userRepository psql.IUserRepository,
	vacancyRepository psql.IVacancyRepository,
	cvRepository psql.ICVRepository,
	notificationRepo grpc.INotificationRepository,
) IResponseUsecase {
	return &ResponseUsecase{
		responseRepo:      respondRepository,
		sessionRepo:       sessionRepository,
		userRepo:          userRepository,
		vacancyRepo:       vacancyRepository,
		cvRepo:            cvRepository,
		notificationsRepo: notificationRepo,
	}
}

func (responseUsecase *ResponseUsecase) RespondToVacancy(ctx context.Context, vacancyID, cvID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
		"cv_id":      cvID,
	}).
		Info("responding to vacancy")

	currUserID := contextUtils.GetUserIDFromCtx(ctx)

	userRole, err := responseUsecase.userRepo.GetRoleById(ctx, currUserID)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("error while getting role by id")
		return err
	} else if userRole != domain.Applicant {
		contextLogger.Error(ErrInapropriateRole.Error())
		return ErrInapropriateRole
	}

	respErr := responseUsecase.responseRepo.RespondToVacancy(ctx, vacancyID, cvID)
	if respErr != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": respErr,
		}).
			Error("error while responding to vacancy")
		return respErr
	}

	empUserID, err := responseUsecase.userRepo.GetUserIDByVacID(ctx, vacancyID)
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("error while getting user_id by vacancy_id")
		return err
	}

	message := notificationsPB.NotificationMessage{
		UserId:  empUserID,
		Message: notificationMessages.NewVacancyResponse,
		VacancyId: int64(vacancyID),
		CvId: int64(cvID),
	}

	sendErr := responseUsecase.notificationsRepo.SendMessage(ctx, &message)
	if sendErr != nil {
		return sendErr
	}
	return nil
}

func (responseUsecase *ResponseUsecase) GetApplicantsList(ctx context.Context, vacancyID int) ([]domain.ApiApplicant, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

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

func (responseUsecase *ResponseUsecase) GetUserResponses(ctx context.Context, userID int) ([]domain.ApiResponse, error) {
	currUserID := contextUtils.GetUserIDFromCtx(ctx)

	userRole, err := responseUsecase.userRepo.GetRoleById(ctx, currUserID)
	if err != nil {
		return nil, err
	} else if userRole != domain.Applicant {
		return nil, ErrInapropriateRole
	}

	if userID != currUserID {
		return nil, ErrForbidden
	}

	responses, err := responseUsecase.responseRepo.GetUserResponses(ctx, userID)
	if err != nil && !errors.Is(err, psql.ErrEntityNotFound) {
		return nil, err
	}
	return responses, nil
}
