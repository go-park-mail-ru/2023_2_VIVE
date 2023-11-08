package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"

	"database/sql"
	"errors"
	"strings"
)

type IUserRepository interface {
	CheckUser(user *domain.User) error
	CheckPasswordById(id int, passwordToCheck string) error
	AddUser(user *domain.User, hasher authUtils.HashGenerator) error
	GetUserIdByEmail(email string) (int, error)
	GetRoleById(userID int) (domain.Role, error)
	GetUserInfo(userID int) (*domain.User, error)
	UpdateUserInfo(userID int, user *domain.UserUpdate) error
	GetUserOrgId(userID int) (int, error)
	UploadAvatarByUserID(userID int, path string) error
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

func (p *psqlUserRepository) checkPasswordByEmail(email, passwordToCheck string) error {
	var actualHash interface{}
	var salt interface{}

	err := p.userStorage.QueryRow(`SELECT pswd, salt FROM hnh_data.user_profile WHERE email = $1`, email).Scan(&actualHash, &salt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEntityNotFound
	} else if err != nil {
		return err
	}

	return p.castRawPasswordAndCompare(actualHash, salt, passwordToCheck)
}

func (p *psqlUserRepository) checkRole(user *domain.User) error {
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

func (p *psqlUserRepository) CheckUser(user *domain.User) error {
	passwordStatus := p.checkPasswordByEmail(user.Email, user.Password)
	if passwordStatus != nil {
		return passwordStatus
	}

	roleStatus := p.checkRole(user)
	if roleStatus != nil {
		return roleStatus
	}

	return nil
}

func (p *psqlUserRepository) CheckPasswordById(id int, passwordToCheck string) error {
	var actualHash interface{}
	var salt interface{}

	err := p.userStorage.QueryRow(`SELECT pswd, salt FROM hnh_data.user_profile WHERE id = $1`, id).Scan(&actualHash, &salt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrEntityNotFound
	} else if err != nil {
		return err
	}

	return p.castRawPasswordAndCompare(actualHash, salt, passwordToCheck)
}

func (p *psqlUserRepository) AddUser(user *domain.User, hasher authUtils.HashGenerator) error {
	var exists bool

	err := p.userStorage.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.user_profile WHERE email = $1)`, user.Email).Scan(&exists)
	if exists {
		return serverErrors.ACCOUNT_ALREADY_EXISTS
	} else if err != nil {
		return err
	}

	hashedPass, salt, err := hasher(user.Password)
	if err != nil {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	if user.Type == domain.Applicant {
		var userID int
		addErr := p.userStorage.QueryRow(`INSERT INTO hnh_data.user_profile `+
			`("email", "pswd", "salt", "first_name", "last_name", "birthday", "phone_number", "location") `+
			`VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			user.Email, hashedPass, salt, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location).
			Scan(&userID)
		if addErr != nil {
			return addErr
		}

		_, appErr := p.userStorage.Exec(`INSERT INTO hnh_data.applicant ("user_id") VALUES ($1)`, userID)
		if appErr != nil {
			return appErr
		}
	} else if user.Type == domain.Employer {
		var userID int
		addErr := p.userStorage.QueryRow(`INSERT INTO hnh_data.user_profile `+
			`("email", "pswd", "salt", "first_name", "last_name", "birthday", "phone_number", "location") `+
			`VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			user.Email, hashedPass, salt, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location).
			Scan(&userID)
		if addErr != nil {
			return addErr
		}

		_, empErr := p.userStorage.Exec(`INSERT INTO hnh_data.employer ("user_id") VALUES ($1)`, userID)
		if empErr != nil {
			return empErr
		}
	} else {
		return serverErrors.INVALID_ROLE
	}

	return nil
}

func (p *psqlUserRepository) GetUserInfo(userID int) (*domain.User, error) {
	user := &domain.User{}

	err := p.userStorage.QueryRow(`SELECT email, first_name, last_name, birthday, phone_number, location, avatar_path `+
		`FROM hnh_data.user_profile WHERE id = $1`, userID).
		Scan(&user.Email, &user.FirstName, &user.LastName, &user.Birthday, &user.PhoneNumber, &user.Location, &user.AvatarPath)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}

	role, err := p.GetRoleById(userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, serverErrors.INTERNAL_SERVER_ERROR
	} else if err != nil {
		return nil, err
	}
	user.Type = role

	user.Email = strings.TrimSpace(user.Email)

	if user.PhoneNumber != nil {
		*user.PhoneNumber = strings.TrimSpace(*user.PhoneNumber)
	}

	return user, nil
}

func (p *psqlUserRepository) GetUserIdByEmail(email string) (int, error) {
	var userID int

	err := p.userStorage.QueryRow(`SELECT id FROM hnh_data.user_profile WHERE email = $1`, email).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEntityNotFound
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *psqlUserRepository) GetRoleById(userID int) (domain.Role, error) {
	var isApplicant bool

	appErr := p.userStorage.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.applicant WHERE user_id = $1)`, userID).Scan(&isApplicant)
	if isApplicant {
		return domain.Applicant, nil
	} else if appErr != nil {
		return "", appErr
	}

	var isEmployer bool

	empErr := p.userStorage.QueryRow(`SELECT EXISTS (SELECT id FROM hnh_data.employer WHERE user_id = $1)`, userID).Scan(&isEmployer)
	if isEmployer {
		return domain.Employer, nil
	} else if empErr != nil {
		return "", empErr
	}

	return "", ErrEntityNotFound
}

func (p *psqlUserRepository) UpdateUserInfo(userID int, user *domain.UserUpdate) error {
	if user.NewPassword != "" {
		hashedPass, salt, err := authUtils.GenerateHash(user.NewPassword)
		if err != nil {
			return serverErrors.INTERNAL_SERVER_ERROR
		}

		_, updErr := p.userStorage.Exec(`UPDATE hnh_data.user_profile SET `+
			`"email" = $1, "pswd" = $2, "salt" = $3, "first_name" = $4, "last_name" = $5, `+
			`"birthday" = $6, "phone_number" = $7, "location" = $8 `+
			`WHERE id = $9`,
			user.Email, hashedPass, salt, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location, userID)
		if updErr != nil {
			return updErr
		}
	} else {
		_, updErr := p.userStorage.Exec(`UPDATE hnh_data.user_profile SET `+
			`"email" = $1, "first_name" = $2, "last_name" = $3, `+
			`"birthday" = $4, "phone_number" = $5, "location" = $6 `+
			`WHERE id = $7`,
			user.Email, user.FirstName, user.LastName, user.Birthday, user.PhoneNumber, user.Location, userID)
		if updErr != nil {
			return updErr
		}
	}

	return nil
}

func (p *psqlUserRepository) GetUserOrgId(userID int) (int, error) {
	var orgID int

	err := p.userStorage.QueryRow(`SELECT organization_id FROM hnh_data.employer WHERE user_id = $1`, userID).Scan(&orgID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrEntityNotFound
	} else if err != nil {
		return 0, err
	}

	return orgID, nil
}

func (p *psqlUserRepository) UploadAvatarByUserID(userID int, path string) error {
	_, uplErr := p.userStorage.Exec(`UPDATE hnh_data.user_profile SET `+
		`"avatar_path" = $1 `+
		`WHERE id = $2`, path, userID)

	if uplErr != nil {
		return uplErr
	}

	return nil
}
