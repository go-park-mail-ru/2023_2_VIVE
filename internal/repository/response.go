package repository

type IResponseRepository interface {
	RespondToVacancy(vacancyID, cvID int) error
}

type psqlResponseRepository struct {
}

func NewPsqlResponseRepository() IResponseRepository {
	return &psqlResponseRepository{}
}

func (p *psqlResponseRepository) RespondToVacancy(vacancyID, cvID int) error {
	return nil
}
