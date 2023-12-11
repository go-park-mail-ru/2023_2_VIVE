package psql

import (
	"HnH/internal/domain"
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

type IResponseRepository interface {
	RespondToVacancy(ctx context.Context, vacancyID, cvID int) error
	GetVacanciesIdsByCVId(ctx context.Context, cvID int) ([]int, error)
	GetAttachedCVs(ctx context.Context, vacancyID int) ([]int, error)
	GetUserResponses(ctx context.Context, userID int) ([]domain.ApiResponse, error)
}

type psqlResponseRepository struct {
	responseStorage *sql.DB
}

func NewPsqlResponseRepository(db *sql.DB) IResponseRepository {
	return &psqlResponseRepository{
		responseStorage: db,
	}
}

func (p *psqlResponseRepository) RespondToVacancy(ctx context.Context, vacancyID, cvID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.Info("adding new responce to vacancy from cv to postgres")
	result, err := p.responseStorage.Exec(`INSERT INTO hnh_data.response ("vacancy_id", "cv_id") VALUES ($1, $2)`, vacancyID, cvID)
	if err == sql.ErrNoRows {
		return ErrNotInserted
	}
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (p *psqlResponseRepository) GetVacanciesIdsByCVId(ctx context.Context, cvID int) ([]int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id": cvID,
	}).
		Info("getting 'vacancy_id' list to witch 'cv_id' is responded from postgres")
	rows, err := p.responseStorage.Query(`SELECT vacancy_id FROM hnh_data.response WHERE cv_id = $1`, cvID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []int{}
	for rows.Next() {
		var vacancyID int

		err := rows.Scan(&vacancyID)
		if err != nil {
			return nil, err
		}
		result = append(result, vacancyID)
	}

	return result, nil
}

func (p *psqlResponseRepository) GetAttachedCVs(ctx context.Context, vacancyID int) ([]int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("getting 'cv_id' list attached to  'vacancy_id' from postgres")
	rows, err := p.responseStorage.Query(`SELECT cv_id FROM hnh_data.response WHERE vacancy_id = $1`, vacancyID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := []int{}
	for rows.Next() {
		var cvID int

		err := rows.Scan(&cvID)
		if err != nil {
			return nil, err
		}
		result = append(result, cvID)
	}

	return result, nil
}

func (p *psqlResponseRepository) GetUserResponses(ctx context.Context, userID int) ([]domain.ApiResponse, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)
	contextLogger.WithFields(logrus.Fields{
		"user_id": userID,
	}).
		Info("getting user's responses from postgres")

	query := `SELECT
			r.id,
			v."name",
			r.vacancy_id,
			v.organization_name,
			v.employer_id,
			r.created_at,
			r.updated_at
		FROM
			hnh_data.response r
		JOIN hnh_data.cv c ON
			r.cv_id = c.id
		JOIN hnh_data.applicant a ON
			c.applicant_id = a.id
		JOIN hnh_data.vacancy v ON
			r.vacancy_id = v.id
		WHERE
			a.user_id = $1;`

	rows, err := p.responseStorage.Query(query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []domain.ApiResponse{}
	for rows.Next() {
		var response domain.ApiResponse

		err := rows.Scan(
			&response.Id,
			&response.VacancyName,
			&response.VacancyID,
			&response.OrganizationName,
			&response.EmployerID,
			&response.CreatedAt,
			&response.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, response)
	}

	return result, nil
}
