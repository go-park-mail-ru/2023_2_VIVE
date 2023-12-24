package psql

import (
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IStatRepository interface {
	CountVacancyViews(ctx context.Context, vacancyID int) (int, error)
	AddVacancyView(ctx context.Context, vacancyID, applicantID int) error
	// CountVacancyViewsByIDs(ctx context.Context, vacancyIDs ...int) (int, error)
	CountCvViews(ctx context.Context, cvID int) (int, error)
	AddCvView(ctx context.Context, cvID, employerID int) error
	CountVacancyResponses(ctx context.Context, vacancyID int) (int, error)
	CountCvResponses(ctx context.Context, cvID int) (int, error)
}

type psqlStatRepository struct {
	DB *sql.DB
}

func NewPsqlStatRepository(db *sql.DB) IStatRepository {
	return &psqlStatRepository{
		DB: db,
	}
}

func (repo *psqlStatRepository) CountVacancyViews(ctx context.Context, vacancyID int) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("counting vacancy views in postgres")

	query := `SELECT
			count(*)
		FROM
			hnh_data.vacancy_view vv
		WHERE
			vv.vacancy_id = $1`

	var count int
	err := repo.DB.QueryRow(query, vacancyID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *psqlStatRepository) AddVacancyView(ctx context.Context, vacancyID, applicantID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id":   vacancyID,
		"applicant_id": applicantID,
	}).
		Info("adding vacancy view in postgres")

	query := `INSERT INTO hnh_data.vacancy_view (vacancy_id, applicant_id, created_at)
	VALUES ($1, $2, now())`

	result, err := repo.DB.Exec(query, vacancyID, applicantID)
	if errors.Is(err, sql.ErrNoRows) {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("could not insert data")
		return ErrNotInserted
	}
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err_msg": err,
		}).
			Error("error while adding view to vacancy")
		if pgErr, ok := err.(pgx.PgError); ok {
			contextLogger.WithFields(logrus.Fields{
				"pg_err": pgErr,
				"code": pgErr.Code,
			}).
				Error("pgx error")
			switch pgErr.Code {
			case "23505":
				contextLogger.WithFields(logrus.Fields{
					"err": pgErr,
				}).
					Error("could not duplicate data")
				return ErrRecordAlredyExists
			}
		}
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("could not insert data")
		return err
	}
	return nil
}

// func (repo *psqlStatRepository) CountVacancyViewsByIDs(ctx context.Context, vacancyIDs ...int) ([]int, error) {
// 	contextLogger := contextUtils.GetContextLogger(ctx)

// 	contextLogger.WithFields(logrus.Fields{
// 		"vacancy_ids": vacancyIDs,
// 	}).
// 		Info("counting vacancy views by vacancies ids in postgres")

// 	query := `SELECT
// 			count(*)
// 		FROM
// 			hnh_data.vacancy_view vv
// 		WHERE
// 			vv.vacancy_id = ANY($1)`

// 	var count int
// 	err := repo.DB.QueryRow(query, pq.Array(vacancyIDs))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

func (repo *psqlStatRepository) CountCvViews(ctx context.Context, cvID int) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id": cvID,
	}).
		Info("counting cv views in postgres")

	query := `SELECT
			count(*)
		FROM
			hnh_data.cv_view cv
		WHERE
			cv.cv_id = $1`

	var count int
	err := repo.DB.QueryRow(query, cvID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *psqlStatRepository) AddCvView(ctx context.Context, cvID, employerID int) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id":       cvID,
		"employer_id": employerID,
	}).
		Info("adding vacancy view in postgres")

	query := `INSERT INTO hnh_data.cv_view (cv_id, employer_id, created_at)
	VALUES ($1, $2, now())`

	result, err := repo.DB.Exec(query, cvID, employerID)
	if err == sql.ErrNoRows {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("could not insert data")
		return ErrNotInserted
	}
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			switch pgErr.Code {
			case "23505":
				contextLogger.WithFields(logrus.Fields{
					"err": pgErr,
				}).
					Error("could not duplicate data")
				return ErrRecordAlredyExists
			}
		}
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		contextLogger.WithFields(logrus.Fields{
			"err": err,
		}).
			Error("could not insert data")
		return err
	}
	return nil
}

func (repo *psqlStatRepository) CountVacancyResponses(ctx context.Context, vacancyID int) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacancyID,
	}).
		Info("counting vacancy responses in postgres")

	query := `SELECT
			count(*)
		FROM
			hnh_data.response r
		WHERE
			r.vacancy_id = $1`

	var count int
	err := repo.DB.QueryRow(query, vacancyID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *psqlStatRepository) CountCvResponses(ctx context.Context, cvID int) (int, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id": cvID,
	}).
		Info("counting cv responses in postgres")

	query := `SELECT
			count(*)
		FROM
			hnh_data.response r
		WHERE
			r.cv_id = $1`

	var count int
	err := repo.DB.QueryRow(query, cvID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
