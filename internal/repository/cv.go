package repository

import (
	"HnH/internal/domain"
)

type ICVRepository interface {
	GetCVById(cvID int) (*domain.CV, error)
	GetCVsByIds(idList []int) ([]domain.CV, error)
	GetCVsByUserId(userID int) ([]domain.CV, error)
	AddCV(cv *domain.CV) (int, error)
	GetOneOfUsersCV(userID, cvID int) (*domain.CV, error)
	UpdateOneOfUsersCV(userID, cvID int) error
	DeleteOneOfUsersCV(userID, cvID int) error
}

type psqlCVRepository struct {
}

func NewPsqlCVRepository() ICVRepository {
	return &psqlCVRepository{}
}

func (p *psqlCVRepository) GetCVById(cvID int) (*domain.CV, error) {
	return nil, nil
}

func (p *psqlCVRepository) GetCVsByIds(idList []int) ([]domain.CV, error) {
	return nil, nil
}

func (p *psqlCVRepository) GetCVsByUserId(userID int) ([]domain.CV, error) {
	return nil, nil
}

func (p *psqlCVRepository) AddCV(cv *domain.CV) (int, error) {
	return 0, nil
}

func (p *psqlCVRepository) GetOneOfUsersCV(userID, cvID int) (*domain.CV, error) {
	return nil, nil
}

func (p *psqlCVRepository) UpdateOneOfUsersCV(userID, cvID int) error {
	return nil
}

func (p *psqlCVRepository) DeleteOneOfUsersCV(userID, cvID int) error {
	return nil
}
