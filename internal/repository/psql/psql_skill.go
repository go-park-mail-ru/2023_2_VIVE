package psql

import (
	"HnH/pkg/contextUtils"
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

type ISkillRepository interface {
	AddSkillsByVacID(ctx context.Context, vacID int, skills []string) error
	GetSkillsByVacID(ctx context.Context, vacID int) ([]string, error)
	AddSkillsByCvID(ctx context.Context, cvsID int, skills []string) error
	GetSkillsByCvID(ctx context.Context, cvID int) ([]string, error)
}

type psqlSkillRepository struct {
	DB *sql.DB
}

func NewPsqlSkillRepository(db *sql.DB) ISkillRepository {
	return &psqlSkillRepository{
		DB: db,
	}
}

// TODO: optimize
func (repo *psqlSkillRepository) AddSkillsByVacID(ctx context.Context, vacID int, skills []string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacID,
		"skills":     skills,
	}).
		Info("adding 'skills' by 'vacancy_id'")

	for _, skill := range skills {
		skillQuery := `INSERT INTO hnh_data.skill ("name")
		VALUES ($1)
		ON CONFLICT ("name") DO NOTHING
		RETURNING id`

		var skillID int
		err := repo.DB.QueryRow(skillQuery, skill).Scan(&skillID)
		if err == sql.ErrNoRows {
			selectErr := repo.DB.QueryRow(
				`SELECT s.id FROM hnh_data.skill s WHERE s.name = $1`,
				skill,
			).
				Scan(&skillID)
			if selectErr == sql.ErrNoRows {
				return ErrNotInserted
			}
			if selectErr != nil {
				return selectErr
			}
		} else if err != nil {
			return err
		}

		vacSkillQuery := `INSERT INTO hnh_data.vacancy_skill_assign (vacancy_id, skill_id)
						VALUES ($1, $2)`

		result, err := repo.DB.Exec(vacSkillQuery, vacID, skillID)
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
	}
	return nil
}

func (repo *psqlSkillRepository) GetSkillsByVacID(ctx context.Context, vacID int) ([]string, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"vacancy_id": vacID,
	}).
		Info("getting skills by 'vacancy_id'")

	query := `SELECT
			s."name"
		FROM
			hnh_data.vacancy_skill_assign vsa
		JOIN hnh_data.skill s ON
			vsa.skill_id = s.id
		WHERE
			vsa.vacancy_id = $1`

	rows, err := repo.DB.Query(query, vacID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skillsToReturn := []string{}

	for rows.Next() {
		var skill string

		err := rows.Scan(&skill)
		if err != nil {
			return nil, err
		}
		skillsToReturn = append(skillsToReturn, skill)
	}

	return skillsToReturn, nil
}

// TODO: optimize
func (repo *psqlSkillRepository) AddSkillsByCvID(ctx context.Context, cvID int, skills []string) error {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id":  cvID,
		"skills": skills,
	}).
		Info("adding 'skills' by 'cv_id'")

	for _, skill := range skills {
		skillQuery := `INSERT INTO hnh_data.skill ("name")
		VALUES ($1)
		ON CONFLICT ("name") DO NOTHING
		RETURNING id`

		var skillID int
		err := repo.DB.QueryRow(skillQuery, skill).Scan(&skillID)
		if err == sql.ErrNoRows {
			contextLogger.WithFields(logrus.Fields{
				"err":   err,
				"skill": skill,
			}).
				Debug("inserting skill already exists")
			selectErr := repo.DB.QueryRow(
				`SELECT s.id FROM hnh_data.skill s WHERE s.name = $1`,
				skill,
			).
				Scan(&skillID)
			if selectErr == sql.ErrNoRows {
				return ErrNotInserted
			}
			if selectErr != nil {
				return selectErr
			}
		} else if err != nil {
			return err
		}

		contextLogger.WithFields(logrus.Fields{
			"cv_id":    cvID,
			"skill_id": skillID,
		}).
			Debug("inserting data into hnh_data.cv_skill_assign")

		cvSkillQuery := `INSERT INTO hnh_data.cv_skill_assign (cv_id, skill_id)
						VALUES ($1, $2)`

		result, err := repo.DB.Exec(cvSkillQuery, cvID, skillID)
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
	}
	return nil
}

func (repo *psqlSkillRepository) GetSkillsByCvID(ctx context.Context, cvID int) ([]string, error) {
	contextLogger := contextUtils.GetContextLogger(ctx)

	contextLogger.WithFields(logrus.Fields{
		"cv_id": cvID,
	}).
		Info("getting skills by 'cv_id'")

	query := `SELECT
			s."name"
		FROM
			hnh_data.cv_skill_assign vsa
		JOIN hnh_data.skill s ON
			vsa.skill_id = s.id
		WHERE
			vsa.cv_id = $1`

	rows, err := repo.DB.Query(query, cvID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skillsToReturn := []string{}

	for rows.Next() {
		var skill string

		err := rows.Scan(&skill)
		if err != nil {
			return nil, err
		}
		skillsToReturn = append(skillsToReturn, skill)
	}

	return skillsToReturn, nil
}
