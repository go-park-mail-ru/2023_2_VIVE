package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
)

type IOrganizationUsecase interface {
	AddOrganization(organization *domain.DbOrganization) (int, error)
	// SignUp(user *domain.ApiUserReg, expiryUnixSeconds int64) (string, error)
	// GetInfo(sessionID string) (*domain.DbUser, error)
	// UpdateInfo(sessionID string, user *domain.UserUpdate) error
	// UploadAvatar(sessionID, path string) error
	// GetAvatar(sessionID string) ([]byte, error)
}

type OrganizationUsecase struct {
	orgRepo psql.IOrganizationRepository
	// cvRepo       psql.ICVRepository
	// sessionRepo  redisRepo.ISessionRepository
	// userRepo     psql.IUserRepository
	// responseRepo psql.IResponseRepository
	// vacancyRepo  psql.IVacancyRepository
}

func NewOrganizationUsecase(organizationRepo psql.IOrganizationRepository) IOrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo: organizationRepo,
		// userRepo:    userRepository,
		// sessionRepo: sessionRepository,
	}
}

func (orgUsecase *OrganizationUsecase) AddOrganization(organization *domain.DbOrganization) (int, error) {
	orgID, err := orgUsecase.orgRepo.AddOrganization(organization)
	return orgID, err
}
