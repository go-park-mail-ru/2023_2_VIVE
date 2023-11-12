package psql

import (
	"database/sql"
)

type IExperienceRepository interface {
	// AddExperience(cvID int, experience domain.DbExperience)
}

type psqlExperienceRepository struct {
	DB *sql.DB
}

func NewPsqlExperienceRepository(db *sql.DB) IExperienceRepository {
	return &psqlExperienceRepository{
		DB: db,
	}
}

// func (repo *psqlExperienceRepository) AddExperience(cvID int, experience domain.DbExperience) (int, error) {
// 	Query := `INSERT
// 		INTO
// 		hnh_data.experience (
// 			cv_id,
// 			organization_name,
// 			"position",
// 			description,
// 			start_date,
// 			end_date
// 		)
// 	VALUES 
// 		($1, $2, $3, $4, $5, $6)`
// }
