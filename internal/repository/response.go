package repository

type IResponseRepository interface {
	RespondToVacancy(vacancyID, cvID int) error
	GetVacanciesIdsByCVId(cvID int) ([]int, error)
	GetAttachedCVs(vacancyID int) ([]int, error)
}

type psqlResponseRepository struct {
}

func NewPsqlResponseRepository() IResponseRepository {
	return &psqlResponseRepository{}
}

func (p *psqlResponseRepository) RespondToVacancy(vacancyID, cvID int) error {
	return nil
}

func (p *psqlResponseRepository) GetVacanciesIdsByCVId(cvID int) ([]int, error) {
	return []int{}, nil
}

func (p *psqlResponseRepository) GetAttachedCVs(vacancyID int) ([]int, error) {
	return []int{}, nil
}
