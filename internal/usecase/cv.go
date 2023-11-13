package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/serverErrors"
	"fmt"
)

type ICVUsecase interface {
	GetCVById(sessionID string, cvID int) (*domain.ApiCV, error)
	GetCVList(sessionID string) ([]domain.ApiCV, error)
	AddNewCV(sessionID string, cv *domain.ApiCV) (int, error)
	GetCVOfUserById(sessionID string, cvID int) (*domain.ApiCV, error)
	UpdateCVOfUserById(sessionID string, cvID int, cv *domain.ApiCV) error
	DeleteCVOfUserById(sessionID string, cvID int) error
}

type CVUsecase struct {
	cvRepo       psql.ICVRepository
	sessionRepo  redisRepo.ISessionRepository
	userRepo     psql.IUserRepository
	responseRepo psql.IResponseRepository
	vacancyRepo  psql.IVacancyRepository
}

func NewCVUsecase(cvRepository psql.ICVRepository,
	sessionRepository redisRepo.ISessionRepository,
	userRepository psql.IUserRepository,
	responseRepository psql.IResponseRepository,
	vacancyRepository psql.IVacancyRepository) ICVUsecase {
	return &CVUsecase{
		cvRepo:       cvRepository,
		sessionRepo:  sessionRepository,
		userRepo:     userRepository,
		responseRepo: responseRepository,
		vacancyRepo:  vacancyRepository,
	}
}

func (cvUsecase *CVUsecase) validateSessionAndGetUserId(sessionID string) (int, error) {
	validStatus := cvUsecase.sessionRepo.ValidateSession(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userID, err := cvUsecase.sessionRepo.GetUserIdBySession(sessionID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (cvUsecase *CVUsecase) validateRoleAndGetUserId(sessionID string, requiredRole domain.Role) (int, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userRole, err := cvUsecase.userRepo.GetRoleById(userID)
	if err != nil {
		return 0, err
	} else if userRole != requiredRole {
		return 0, INAPPROPRIATE_ROLE
	}

	return userID, nil
}

func (cvUsecase *CVUsecase) constructApiCV(cv *domain.DbCV, exps []domain.DbExperience, edInsts []domain.DbEducationInstitution) *domain.ApiCV {
	apiCV := cv.ToAPI()

	apiInsts := make([]domain.ApiEducationInstitution, len(edInsts))
	for i := range edInsts {
		apiInsts[i] = *edInsts[i].ToAPI()
	}
	apiCV.EducationInstitutions = apiInsts

	apiExps := make([]domain.ApiExperience, len(exps))
	for i := range exps {
		apiExps[i] = *exps[i].ToAPI()
	}
	apiCV.Experience = apiExps

	return apiCV
}

// Finds cv that responded to one of the current user's vacancy
func (cvUsecase *CVUsecase) GetCVById(sessionID string, cvID int) (*domain.ApiCV, error) {
	userID, validStatus := cvUsecase.validateRoleAndGetUserId(sessionID, domain.Employer)
	if validStatus != nil {
		return nil, validStatus
	}

	vacIdsList, err := cvUsecase.responseRepo.GetVacanciesIdsByCVId(cvID)
	if err != nil {
		return nil, err
	}

	userOrgID, err := cvUsecase.userRepo.GetUserOrgId(userID)
	if err != nil {
		return nil, err
	}

	_, err = cvUsecase.vacancyRepo.GetVacanciesByIds(userOrgID, vacIdsList)
	if err == psql.ErrEntityNotFound {
		return nil, serverErrors.FORBIDDEN
	}
	if err != nil {
		return nil, err
	}

	cv, exps, edInsts, err := cvUsecase.cvRepo.GetCVById(cvID)
	if err != nil {
		return nil, err
	}

	apiCV := cvUsecase.constructApiCV(cv, exps, edInsts)

	return apiCV, nil
}

func (cvUsecase *CVUsecase) combineDbCVs(cvs []domain.DbCV, exps []domain.DbExperience, insts []domain.DbEducationInstitution) []domain.ApiCV {
	res := []domain.ApiCV{}
	for _, cv := range cvs {
		cvID := cv.ID
		cvExps := []domain.DbExperience{}
		for _, exp := range exps {
			if exp.CvID == cvID {
				cvExps = append(cvExps, exp)
			}
		}
		cvInsts := []domain.DbEducationInstitution{}
		for _, inst := range insts {
			if inst.CvID == cvID {
				cvInsts = append(cvInsts, inst)
			}
		}
		res = append(res, *cvUsecase.constructApiCV(&cv, cvExps, cvInsts))
	}
	return res
}

func (cvUsecase *CVUsecase) GetCVList(sessionID string) ([]domain.ApiCV, error) {

	userID, validStatus := cvUsecase.validateRoleAndGetUserId(sessionID, domain.Applicant)
	if validStatus != nil {
		return nil, validStatus
	}
	fmt.Printf("userID: %v\n", userID)
	fmt.Println("before getting cvs")

	cvs, exps, insts, err := cvUsecase.cvRepo.GetCVsByUserId(userID)
	if err != nil {
		return nil, err
	}
	fmt.Println("after getting cvs")

	apiCvs := cvUsecase.combineDbCVs(cvs, exps, insts)

	return apiCvs, nil
}

func (cvUsecase *CVUsecase) getExperiences(apiCV *domain.ApiCV) []domain.DbExperience {
	res := []domain.DbExperience{}
	for _, experience := range apiCV.Experience {
		res = append(res, *experience.ToDb())
	}
	return res
}

func (cvUsecase *CVUsecase) getEducationInstitutions(apiCV *domain.ApiCV) []domain.DbEducationInstitution {
	res := []domain.DbEducationInstitution{}
	for _, institution := range apiCV.EducationInstitutions {
		res = append(res, *institution.ToDb())
	}
	return res
}

func (cvUsecase *CVUsecase) getDataFromApiCV(apiCV *domain.ApiCV) ([]domain.DbExperience, []domain.DbEducationInstitution, *domain.DbCV) {
	dbExperiences := cvUsecase.getExperiences(apiCV)
	dbEducationInstitutions := cvUsecase.getEducationInstitutions(apiCV)
	dbCV := apiCV.ToDb()

	return dbExperiences, dbEducationInstitutions, dbCV
}

func (cvUsecase *CVUsecase) AddNewCV(sessionID string, cv *domain.ApiCV) (int, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	// fmt.Println(userID)
	if validStatus != nil {
		return 0, validStatus
	}

	dbExperiences, dbEducationInstitutions, dbCV := cvUsecase.getDataFromApiCV(cv)

	cvID, addErr := cvUsecase.cvRepo.AddCV(userID, dbCV, dbExperiences, dbEducationInstitutions)
	if addErr != nil {
		return 0, addErr
	}

	return cvID, nil
}

func (cvUsecase *CVUsecase) GetCVOfUserById(sessionID string, cvID int) (*domain.ApiCV, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	cv, exps, insts, err := cvUsecase.cvRepo.GetOneOfUsersCV(userID, cvID)
	if err != nil {
		return nil, err
	}

	apiCv := cvUsecase.constructApiCV(cv, exps, insts)

	return apiCv, nil
}

func (cvUsecase *CVUsecase) UpdateCVOfUserById(sessionID string, cvID int, cv *domain.ApiCV) error {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	dbExperiences, dbEducationInstitutions, dbCV := cvUsecase.getDataFromApiCV(cv)
	// fmt.Printf("before update db\n")
	updStatus := cvUsecase.cvRepo.UpdateOneOfUsersCV(userID, cvID, dbCV, dbExperiences, dbEducationInstitutions)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (cvUsecase *CVUsecase) DeleteCVOfUserById(sessionID string, cvID int) error {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := cvUsecase.cvRepo.DeleteOneOfUsersCV(userID, cvID)
	if delStatus != nil {
		return delStatus
	}

	return nil
}
