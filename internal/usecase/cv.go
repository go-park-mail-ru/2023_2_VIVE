package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/internal/repository/redisRepo"
	"HnH/pkg/castUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/serverErrors"
	"HnH/pkg/utils"
	"context"

	"github.com/sirupsen/logrus"
)

type ICVUsecase interface {
	GetCVById(ctx context.Context, sessionID string, cvID int) (*domain.ApiCV, error)
	GetCVList(ctx context.Context, sessionID string) ([]domain.ApiCV, error)
	AddNewCV(ctx context.Context, sessionID string, cv *domain.ApiCV) (int, error)
	GetCVOfUserById(ctx context.Context, sessionID string, cvID int) (*domain.ApiCV, error)
	UpdateCVOfUserById(ctx context.Context, sessionID string, cvID int, cv *domain.ApiCV) error
	DeleteCVOfUserById(ctx context.Context, sessionID string, cvID int) error
	SearchCVs(ctx context.Context, query string, pageNumber, resultsPerPage int64) (domain.ApiMetaCV, error)
}

type CVUsecase struct {
	cvRepo           psql.ICVRepository
	expRepo          psql.IExperienceRepository
	instRepo         psql.IEducationInstitutionRepository
	sessionRepo      redisRepo.ISessionRepository
	userRepo         psql.IUserRepository
	responseRepo     psql.IResponseRepository
	vacancyRepo      psql.IVacancyRepository
	searchEngineRepo grpc.ISearchEngineRepository
}

func NewCVUsecase(
	cvRepository psql.ICVRepository,
	expRepository psql.IExperienceRepository,
	instRepository psql.IEducationInstitutionRepository,
	sessionRepository redisRepo.ISessionRepository,
	userRepository psql.IUserRepository,
	responseRepository psql.IResponseRepository,
	vacancyRepository psql.IVacancyRepository,
	searchEngineRepository grpc.ISearchEngineRepository,
) ICVUsecase {
	return &CVUsecase{
		cvRepo:           cvRepository,
		expRepo:          expRepository,
		instRepo:         instRepository,
		sessionRepo:      sessionRepository,
		userRepo:         userRepository,
		responseRepo:     responseRepository,
		vacancyRepo:      vacancyRepository,
		searchEngineRepo: searchEngineRepository,
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

func (cvUsecase *CVUsecase) validateRoleAndGetUserId(ctx context.Context, sessionID string, requiredRole domain.Role) (int, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return 0, validStatus
	}

	userRole, err := cvUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return 0, err
	} else if userRole != requiredRole {
		return 0, ErrInapropriateRole
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
func (cvUsecase *CVUsecase) GetCVById(ctx context.Context, sessionID string, cvID int) (*domain.ApiCV, error) {
	userID, validStatus := cvUsecase.validateRoleAndGetUserId(ctx, sessionID, domain.Employer)
	if validStatus != nil {
		return nil, validStatus
	}

	vacIdsList, err := cvUsecase.responseRepo.GetVacanciesIdsByCVId(ctx, cvID)
	if err != nil {
		return nil, err
	}

	userEmpID, err := cvUsecase.userRepo.GetUserEmpId(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = cvUsecase.vacancyRepo.GetEmpVacanciesByIds(ctx, userEmpID, vacIdsList)
	if err == psql.ErrEntityNotFound {
		return nil, serverErrors.FORBIDDEN
	}
	if err != nil {
		return nil, err
	}

	cv, exps, edInsts, err := cvUsecase.cvRepo.GetCVById(ctx, cvID)
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

func (cvUsecase *CVUsecase) GetCVList(ctx context.Context, sessionID string) ([]domain.ApiCV, error) {

	userID, validStatus := cvUsecase.validateRoleAndGetUserId(ctx, sessionID, domain.Applicant)
	if validStatus != nil {
		return nil, validStatus
	}
	// fmt.Printf("userID: %v\n", userID)
	// fmt.Println("before getting cvs")

	cvs, exps, insts, err := cvUsecase.cvRepo.GetCVsByUserId(ctx, userID)
	if err != nil && err != psql.ErrEntityNotFound {
		contextLogger := contextUtils.GetContextLogger(ctx)
		contextLogger.WithFields(logrus.Fields{
			"err":     err,
			"user_id": userID,
		}).
			Debug("got 'err' while trying to get all user's cvs by 'user_id'")
		return nil, err
	}
	// fmt.Println("after getting cvs")

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

func (cvUsecase *CVUsecase) AddNewCV(ctx context.Context, sessionID string, cv *domain.ApiCV) (int, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	// fmt.Println(userID)
	if validStatus != nil {
		return 0, validStatus
	}

	dbExperiences, dbEducationInstitutions, dbCV := cvUsecase.getDataFromApiCV(cv)

	cvID, addErr := cvUsecase.cvRepo.AddCV(ctx, userID, dbCV, dbExperiences, dbEducationInstitutions)
	if addErr != nil {
		return 0, addErr
	}

	return cvID, nil
}

func (cvUsecase *CVUsecase) GetCVOfUserById(ctx context.Context, sessionID string, cvID int) (*domain.ApiCV, error) {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return nil, validStatus
	}

	cv, exps, insts, err := cvUsecase.cvRepo.GetOneOfUsersCV(ctx, userID, cvID)
	if err != nil {
		return nil, err
	}

	apiCv := cvUsecase.constructApiCV(cv, exps, insts)

	return apiCv, nil
}

func (cvUsecase *CVUsecase) getExpsBatches(
	expsFromApi []domain.DbExperience,
	expIDsFromDB []int,
) (toDeleteIDs []int, toUpdate, toInsert []domain.DbExperience) {
	var idsToUpdate []int
	for _, exp := range expsFromApi {
		if !utils.Contains(exp.ID, expIDsFromDB) {
			toInsert = append(toInsert, exp)
		} else {
			toUpdate = append(toUpdate, exp)
			idsToUpdate = append(idsToUpdate, exp.ID)
		}
	}
	toDeleteIDs = utils.Difference(expIDsFromDB, idsToUpdate)
	return
}

func (cvUsecase *CVUsecase) getInstsBatches(
	instsFromApi []domain.DbEducationInstitution,
	instIDsFromDB []int,
) (toDeleteIDs []int, toUpdate, toInsert []domain.DbEducationInstitution) {
	var idsToUpdate []int
	for _, inst := range instsFromApi {
		if !utils.Contains(inst.ID, instIDsFromDB) {
			toInsert = append(toInsert, inst)
		} else {
			toUpdate = append(toUpdate, inst)
			idsToUpdate = append(idsToUpdate, inst.ID)
		}
	}
	toDeleteIDs = utils.Difference(instIDsFromDB, idsToUpdate)
	return
}

func (cvUsecase *CVUsecase) UpdateCVOfUserById(ctx context.Context, sessionID string, cvID int, cv *domain.ApiCV) error {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	expIDs, expErr := cvUsecase.expRepo.GetCVExperiencesIDs(ctx, cvID)
	if expErr != nil && expErr != psql.ErrEntityNotFound {
		return expErr
	}
	instIDs, instErr := cvUsecase.instRepo.GetCVInstitutionsIDs(ctx, cvID)
	if instErr != nil && instErr != psql.ErrEntityNotFound {
		return instErr
	}
	dbExperiences, dbEducationInstitutions, dbCV := cvUsecase.getDataFromApiCV(cv)

	expsIDsToDelete, expsToUpdate, expsToInsert := cvUsecase.getExpsBatches(dbExperiences, expIDs)

	instsIDsToDelete, instsToUpdate, instsToInsert := cvUsecase.getInstsBatches(dbEducationInstitutions, instIDs)

	updStatus := cvUsecase.cvRepo.UpdateOneOfUsersCV(ctx, userID, cvID, dbCV,
		expsIDsToDelete, expsToUpdate, expsToInsert,
		instsIDsToDelete, instsToUpdate, instsToInsert,
	)
	if updStatus != nil {
		return updStatus
	}

	return nil
}

func (cvUsecase *CVUsecase) DeleteCVOfUserById(ctx context.Context, sessionID string, cvID int) error {
	userID, validStatus := cvUsecase.validateSessionAndGetUserId(sessionID)
	if validStatus != nil {
		return validStatus
	}

	delStatus := cvUsecase.cvRepo.DeleteOneOfUsersCV(ctx, userID, cvID)
	if delStatus != nil {
		return delStatus
	}

	return nil
}

// func (cvUsecase *CVUsecase) collectApiCVs(
// 	cvs []domain.DbCV,
// 	exps []domain.DbExperience,
// 	insts []domain.DbEducationInstitution,
// ) []domain.ApiCV {
// 	res := []domain.ApiCV{}
// 	for _, cv := range cvs {
// 		cv := cvUsecase.combineDbCVs()
// 		cvID := cv.ID
// 		cvExps := []domain.DbExperience{}
// 		for _, exp := range exps {
// 			if exp.CvID == cvID {
// 				cvExps = append(cvExps, exp)
// 			}
// 		}
// 		cvInsts := []domain.DbEducationInstitution{}
// 		for _, inst := range insts {
// 			if inst.CvID == cvID {
// 				cvInsts = append(cvInsts, inst)
// 			}
// 		}
// 		res = append(res, *cv.ToAPI())
// 	}
// 	return res
// }

func (cvUsecase *CVUsecase) SearchCVs(
	ctx context.Context,
	query string,
	pageNumber, resultsPerPage int64,
) (domain.ApiMetaCV, error) {
	cvIDs, count, err := cvUsecase.searchEngineRepo.SearchCVsIDs(ctx, query, pageNumber, resultsPerPage)
	if err != nil {
		return domain.ApiMetaCV{
			Count: 0,
			CVs:   nil,
		}, nil
	}

	// dbCvs, cvErr := cvUsecase.cvRepo.GetCVsByIds(ctx, castUtils.Int64SliceToIntSlice(cvIDs))
	dbCvs, dbExps, dbInsts, cvErr := cvUsecase.cvRepo.GetCVsByIds(ctx, castUtils.Int64SliceToIntSlice(cvIDs))
	if cvErr == psql.ErrEntityNotFound {
		return domain.ApiMetaCV{
			Count: 0,
			CVs:   nil,
		}, nil
	}
	if cvErr != nil {
		return domain.ApiMetaCV{}, cvErr
	}

	cvsToReturn := cvUsecase.combineDbCVs(dbCvs, dbExps, dbInsts)
	result := domain.ApiMetaCV{
		Count:     count,
		CVs: cvsToReturn,
	}

	return result, nil
}
