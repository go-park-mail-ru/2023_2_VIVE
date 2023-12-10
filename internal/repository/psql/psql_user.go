package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"context"

	"database/sql"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type IUserRepository interface {
	CheckUser(ctx context.Context, user *domain.DbUser) error
	CheckPasswordById(ctx context.Context, id int, passwordToCheck string) error
	AddUser(ctx context.Context, user *domain.ApiUser, hasher authUtils.HashGenerator) error
	GetUserIdByEmail(ctx context.Context, email string) (int, error)
	GetRoleById(ctx context.Context, userID int) (domain.Role, error)
	GetUserInfo(ctx context.Context, userID int) (*domain.DbUser, *int, *int, error)
	GetUserIDByVacID(ctx context.Context, vacancyID int) (int64, error)
	UpdateUserInfo(ctx context.Context, userID int, user *domain.UserUpdate) error
	GetUserEmpId(ctx context.Context, userID int) (int, error)
	UploadAvatarByUserID(ctx context.Context, userID int, path string) error
	GetAvatarByUserID(ctx context.Context, userID int) (string, error)
}

type psqlUserRepository struct {
	userStorage *sql.DB
}

func NewPsqlUserRepository(db *sql.DB) IUserRepository {
	return &psqlUserRepository{
		userStorage: db,
	}
}

func (p *psqlUserRepository) castRawPasswordAndCompare(rawHash, rawSalt interface{}, passwordToCheck string) error {
	castedHash, ok := rawHash.([]byte)
	if !ok {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	castedSalt, ok := rawSalt.([]byte)
	if !ok {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	isEqual := authUtils.ComparePasswordAndHash(passwordToCheck, castedSalt, castedHash)

	if !isEqual {
		return serverErrors.INCORRECT_CREDENTIALS
	}

	return nil
}

func (p *psqlUserRepository) checkPasswordByEmail(logger *logrus.Entry, email, passwordToCheck string) error {
	var actualHash interface{}
	var salt interface{}

	logger.Info("checking password by 'email'in postgres")
	err := p.userStorage.QueryRow(`SELECT pswd, salt FROM hnh_data.user_profile WHERE email = $1`, email).Scan(&actualHash, &salt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEntityNotFound
	} else if err != nil {
		return err
	}

	return p.castRawPasswordAndCompare(actualHash, salt, passwordToCheck)
}

func (p *psqlUserRepository) checkRole(logger *logrus.Entry, user *domain.DbUser) error {
	logger.Info("checking role of user in postgres")
	if user.Type == domain.Employer {
		var isEmployer bool

		empErr := p.userStorage.QueryRow(`SELECT EXISTS `+
			`(SELECT id FROM hnh_data.employer `+
			`WHERE user_id = (SELECT id FROM hnh_data.user_profile WHERE email = $1))`, user.Email).Scan(&isEmployer)
		if empErr != nil {
			return empErr
		} else if !isEmployer {
			return serverErrors.INCORRECT_ROLE
		}
	} else if user.Type == domain.Applicant {
		var isApplicant bool

		appErr := p.userStorage.QueryRow(`SELECT EXISTS `+
			`(SELECT id FROM hnh_data.applicant `+
			`WHERE user_id = (SELECT id FROM hnh_data.user_profile WHERE email = $1))`, user.Email).Scan(&isApplicant)
		if appErr != nil {
			return appErr
		} else if !isApplicant {
			return serverErrors.INCORRECT_ROLE
		}
	} else {
		return serverErrors.INVALID_ROLE
	}

	return nil
}

func (p *psqlUserRepository) CheckUser(ctx context.Context, user *domain.DbUser) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	passwordStatus := p.checkPasswordByEmail(contextLogger, user.Email, user.Password)
	if passwordStatus != nil {
		return passwordStatus
	}

	roleStatus := p.checkRole(contextLogger, user)
	if roleStatus != nil {
		return roleStatus
	}

	return nil
}

func (p *psqlUserRepository) CheckPasswordById(ctx context.Context, id int, passwordToCheck string) error {
	var actualHash interface{}
	var salt interface{}

	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("checking password by 'user_id' in postgres")
	err := p.userStorage.QueryRow(`SELECT pswd, salt FROM hnh_data.user_profile WHERE id = $1`, id).Scan(&actualHash, &salt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEntityNotFound
	} else if err != nil {
		return err
	}

	return p.castRawPasswordAndCompare(actualHash, salt, passwordToCheck)
}

func (p *psqlUserRepository) AddUser(ctx context.Context, user *domain.ApiUser, hasher authUtils.HashGenerator) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	tx, txErr := p.userStorage.Begin()
	if txErr != nil {
		return txErr
	}

	contextLogger.Info("started transaction for adding new user")

	var exists bool

	contextLogger.Info("checking user in postgres")
	err := tx.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.user_profile WHERE email = $1)`, user.Email).Scan(&exists)
	if exists {
		tx.Rollback()
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	} else if err != nil {
		tx.Rollback()
		return err
	}

	hashedPass, salt, err := hasher(user.Password)
	if err != nil {
		tx.Rollback()
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	contextLogger.WithFields(logrus.Fields{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Type,
	}).
		Debug("new user fields")
	if user.Type == domain.Applicant {
		var userID int
		contextLogger.Info("adding applicant to postgres")
		addErr := tx.QueryRow(`INSERT INTO hnh_data.user_profile `+
			`("email", "pswd", "salt", "first_name", "last_name", "birthday", "phone_number", "location") `+
			`VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			user.Email, hashedPass, salt, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location).
			Scan(&userID)

		if addErr != nil {
			tx.Rollback()
			return addErr
		}

		contextLogger.WithFields(logrus.Fields{
			"inserted_user_id": userID,
		}).
			Debug("inserted new user")

		_, appErr := tx.Exec(`INSERT INTO hnh_data.applicant ("user_id") VALUES ($1)`, userID)
		if appErr != nil {
			tx.Rollback()
			return appErr
		}
	} else if user.Type == domain.Employer {
		var userID int

		// employer := domain.DbEmployer{
		// 	OrganizationName:        user.OrganizationName,
		// 	OrganizationDescription: user.OrganizationDescription,
		// }
		// contextLogger.Info("adding organization for employer")
		// orgID, addOrgErr := p.orgRepo.AddTxOrganization(ctx, tx, &employer)
		// if addOrgErr != nil {
		// 	tx.Rollback()
		// 	return addOrgErr
		// }

		contextLogger.Info("adding employer to postgres")
		addErr := tx.QueryRow(`INSERT INTO hnh_data.user_profile `+
			`("email", "pswd", "salt", "first_name", "last_name", "birthday", "phone_number", "location") `+
			`VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			user.Email, hashedPass, salt, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location).
			Scan(&userID)

		if addErr != nil {
			tx.Rollback()
			return addErr
		}

		contextLogger.WithFields(logrus.Fields{
			"inserted_user_id": userID,
		}).
			Debug("inserted new user")

		result, empErr := tx.Exec(`INSERT
				INTO
				hnh_data.employer (user_id, organization_name, organization_description)
			VALUES ($1, $2, $3)`,
			userID,
			user.OrganizationName,
			user.OrganizationDescription,
		)
		if empErr != nil {
			contextLogger.WithFields(logrus.Fields{
				"error": empErr,
			}).
				Error("adding employer failed")
			tx.Rollback()
			return empErr
		}
		_, err = result.RowsAffected()
		if err != nil {
			contextLogger.WithFields(logrus.Fields{
				"error": err,
			}).
				Error("adding employer failed")
			tx.Rollback()
			return err
		}
	} else {
		tx.Rollback()
		return serverErrors.INVALID_ROLE
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return commitErr
	}

	contextLogger.Info("transaction successful")

	return nil
}

func (p *psqlUserRepository) GetUserInfo(ctx context.Context, userID int) (*domain.DbUser, *int, *int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("getting user's info from postgres")
	query := `SELECT
		up.id,
		a.id,
		e.id,
		up.email,
		up.first_name,
		up.last_name,
		up.birthday,
		up.phone_number,
		up."location",
		up.avatar_path
	FROM
		hnh_data.user_profile up
	LEFT JOIN hnh_data.applicant a ON
		a.user_id = up.id
	LEFT JOIN hnh_data.employer e ON
		e.user_id = up.id
	WHERE
		up.id = $1`

	user := &domain.DbUser{}

	var appId, empID *int
	err := p.userStorage.QueryRow(query, userID).
		Scan(
			&user.ID,
			&appId,
			&empID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Birthday,
			&user.PhoneNumber,
			&user.Location,
			&user.AvatarPath,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, ErrEntityNotFound
	} else if err != nil {
		return nil, nil, nil, err
	}

	role, err := p.GetRoleById(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil, serverErrors.INTERNAL_SERVER_ERROR
	} else if err != nil {
		return nil, nil, nil, err
	}
	user.Type = role

	user.Email = strings.TrimSpace(user.Email)

	if user.PhoneNumber != nil {
		*user.PhoneNumber = strings.TrimSpace(*user.PhoneNumber)
	}

	return user, appId, empID, nil
}

func (p *psqlUserRepository) GetUserIDByVacID(ctx context.Context, vacancyID int) (int64, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("getting user's id by vacancy id from postgres")

	query := `SELECT
				up.id
			FROM
				hnh_data.user_profile up
			JOIN hnh_data.employer e ON
				e.user_id = up.id
			JOIN hnh_data.vacancy v ON
				v.employer_id = e.id
			WHERE
				v.id = $1`

	var userID int64

	// var appId, empID *int
	err := p.userStorage.QueryRow(query, vacancyID).
		Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEntityNotFound
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *psqlUserRepository) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	var userID int
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting 'user_id' by 'email' from postgres")
	err := p.userStorage.QueryRow(`SELECT id FROM hnh_data.user_profile WHERE email = $1`, email).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEntityNotFound
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *psqlUserRepository) GetRoleById(ctx context.Context, userID int) (domain.Role, error) {
	var isApplicant bool
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("checking applicant for given 'user_id'")
	appErr := p.userStorage.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.applicant WHERE user_id = $1)`, userID).Scan(&isApplicant)
	if isApplicant {
		return domain.Applicant, nil
	} else if appErr != nil {
		return "", appErr
	}

	var isEmployer bool

	contextLogger.Info("checking employer for given 'user_id'")
	empErr := p.userStorage.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.employer WHERE user_id = $1)`, userID).Scan(&isEmployer)
	if isEmployer {
		return domain.Employer, nil
	} else if empErr != nil {
		return "", empErr
	}

	return "", ErrEntityNotFound
}

func (p *psqlUserRepository) UpdateUserInfo(ctx context.Context, userID int, user *domain.UserUpdate) error {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.Info("updating user's info")
	if user.NewPassword != "" {
		contextLogger.Info("updating with new password")
		hashedPass, salt, err := authUtils.GenerateHash(user.NewPassword)
		if err != nil {
			return serverErrors.INTERNAL_SERVER_ERROR
		}

		query := `UPDATE
			hnh_data.user_profile
		SET
			"email" = $1,
			"pswd" = $2,
			"salt" = $3,
			"first_name" = $4,
			"last_name" = $5,
			"birthday" = $6,
			"phone_number" = $7,
			"location" = $8
		WHERE
			id = $9`

		_, updErr := p.userStorage.Exec(
			query,
			user.Email,
			hashedPass,
			salt,
			user.FirstName,
			user.LastName,
			user.Birthday,
			user.PhoneNumber,
			user.Location,
			userID,
		)
		if updErr != nil {
			return updErr
		}
	} else {
		contextLogger.Info("updating without new password")
		query := `UPDATE
				hnh_data.user_profile
			SET
				"email" = $1,
				"first_name" = $2,
				"last_name" = $3,
				"birthday" = $4,
				"phone_number" = $5,
				"location" = $6
			WHERE
				id = $7`
		_, updErr := p.userStorage.Exec(
			query,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Birthday,
			user.PhoneNumber,
			user.Location,
			userID,
		)
		if updErr != nil {
			return updErr
		}
	}

	return nil
}

func (p *psqlUserRepository) GetUserEmpId(ctx context.Context, userID int) (int, error) {
	var empID int
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting user's 'employer_id'")
	query := `SELECT
		e.id
	FROM
		hnh_data.employer e
	WHERE
		e.user_id = $1`
	err := p.userStorage.QueryRow(query, userID).Scan(&empID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEntityNotFound
	} else if err != nil {
		return 0, err
	}

	return empID, nil
}

func (p *psqlUserRepository) UploadAvatarByUserID(ctx context.Context, userID int, path string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("uploading new avatar by 'user_id'")
	_, uplErr := p.userStorage.Exec(`UPDATE hnh_data.user_profile SET `+
		`"avatar_path" = $1 `+
		`WHERE id = $2`, path, userID)

	if uplErr != nil {
		return uplErr
	}

	return nil
}

func (p *psqlUserRepository) GetAvatarByUserID(ctx context.Context, userID int) (string, error) {
	var path *string
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("getting avatar by 'user_id'")
	err := p.userStorage.QueryRow(`SELECT avatar_path FROM hnh_data.user_profile WHERE id = $1`, userID).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", IncorrectUserID
	} else if err != nil {
		return "", err
	} else if path == nil {
		return "", nil
	}

	return *path, nil
}
