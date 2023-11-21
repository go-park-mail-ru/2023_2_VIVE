package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type IOrganizationRepository interface {
	AddOrganization(ctx context.Context, organization *domain.DbOrganization) (int, error)
	AddTxOrganization(ctx context.Context, tx *sql.Tx, organization *domain.DbOrganization) (int, error)
	// CheckUser(user *domain.DbUser) error
	// CheckPasswordById(id int, passwordToCheck string) error
	// GetUserIdByEmail(email string) (int, error)
	// GetRoleById(userID int) (domain.Role, error)
	// GetUserInfo(userID int) (*domain.DbUser, error)
	// UpdateUserInfo(userID int, user *domain.UserUpdate) error
	// GetUserOrgId(userID int) (int, error)
	// UploadAvatarByUserID(userID int, path string) error
	// GetAvatarByUserID(userID int) (string, error)
}

type psqlOrganizationRepository struct {
	organizationStorage *sql.DB
}

func NewPsqlOrganizationRepository(db *sql.DB) IOrganizationRepository {
	return &psqlOrganizationRepository{
		organizationStorage: db,
	}
}

func (repo *psqlOrganizationRepository) AddOrganization(ctx context.Context, organization *domain.DbOrganization) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new organization")
	query := `INSERT 
	INTO hnh_data.organization 
		("name", description, "location") 
	VALUES 
		($1, $2, $3) 
	RETURNING id`

	var orgID int
	err := repo.organizationStorage.QueryRow(query, organization.Name, organization.Description, organization.Location).Scan(&orgID)
	if err == sql.ErrNoRows {
		return 0, ErrNotInserted
	}
	if err != nil {
		return 0, err
	}

	return orgID, nil
}

func (repo *psqlOrganizationRepository) AddTxOrganization(ctx context.Context, tx *sql.Tx, organization *domain.DbOrganization) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new organization")
	query := `INSERT 
	INTO hnh_data.organization 
		("name", description, "location") 
	VALUES 
		($1, $2, $3) 
	RETURNING id`

	var orgID int
	err := tx.QueryRow(query, organization.Name, organization.Description, organization.Location).Scan(&orgID)
	if err == sql.ErrNoRows {
		return 0, ErrNotInserted
	}
	if err != nil {
		return 0, err
	}

	contextLogger.WithFields(logrus.Fields{
		"inserted_organization_id": orgID,
	}).
		Debug("inserted new orgaization")

	return orgID, nil
}
