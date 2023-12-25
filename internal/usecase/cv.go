package usecase

import (
	"HnH/internal/domain"
	"HnH/internal/repository/grpc"
	"HnH/internal/repository/psql"
	"HnH/pkg/castUtils"
	"HnH/pkg/contextUtils"
	"HnH/pkg/pdf"
	"HnH/pkg/utils"
	"HnH/services/searchEngineService/searchEnginePB"
	"context"
	"errors"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

type ICVUsecase interface {
	GetCVById(ctx context.Context, cvID int) (*domain.ApiCV, error)
	GetCVList(ctx context.Context) ([]domain.ApiCV, error)
	AddNewCV(ctx context.Context, cv *domain.ApiCV) (int, error)
	GetCVOfUserById(ctx context.Context, cvID int) (*domain.ApiCV, error)
	GetApplicantInfo(ctx context.Context, applicantID int) (*domain.ApplicantInfo, error)
	UpdateCVOfUserById(ctx context.Context, cvID int, cv *domain.ApiCV) error
	DeleteCVOfUserById(ctx context.Context, cvID int) error
	SearchCVs(ctx context.Context, options *searchEnginePB.SearchOptions) (domain.ApiMetaCV, error)
	GenerateCVsPDF(ctx context.Context, cvID int) (*gofpdf.Fpdf, error)
}

type CVUsecase struct {
	cvRepo           psql.ICVRepository
	expRepo          psql.IExperienceRepository
	instRepo         psql.IEducationInstitutionRepository
	sessionRepo      grpc.IAuthRepository
	userRepo         psql.IUserRepository
	responseRepo     psql.IResponseRepository
	vacancyRepo      psql.IVacancyRepository
	searchEngineRepo grpc.ISearchEngineRepository
	skillRepo        psql.ISkillRepository
	statRepo         psql.IStatRepository
}

func NewCVUsecase(
	cvRepository psql.ICVRepository,
	expRepository psql.IExperienceRepository,
	instRepository psql.IEducationInstitutionRepository,
	sessionRepository grpc.IAuthRepository,
	userRepository psql.IUserRepository,
	responseRepository psql.IResponseRepository,
	vacancyRepository psql.IVacancyRepository,
	searchEngineRepository grpc.ISearchEngineRepository,
	skillRepository psql.ISkillRepository,
	statRepo psql.IStatRepository,
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
		skillRepo:        skillRepository,
		statRepo:         statRepo,
	}
}

func (cvUsecase *CVUsecase) validateRole(ctx context.Context, userID int, requiredRole domain.Role) error {
	userRole, err := cvUsecase.userRepo.GetRoleById(ctx, userID)
	if err != nil {
		return err
	} else if userRole != requiredRole {
		return ErrInapropriateRole
	}

	return nil
}

func (cvUsecase *CVUsecase) constructApiCV(cv *domain.DbCV, exps []domain.DbExperience, edInsts []domain.DbEducationInstitution) *domain.ApiCV {
	apiCV := cv.ToAPI()

	apiInsts := make([]*domain.ApiEducationInstitution, len(edInsts))
	for i := range edInsts {
		apiInsts[i] = edInsts[i].ToAPI()
	}
	apiCV.EducationInstitutions = apiInsts

	apiExps := make([]*domain.ApiExperience, len(exps))
	for i := range exps {
		apiExps[i] = exps[i].ToAPI()
	}
	apiCV.Experience = apiExps

	return apiCV
}

func (cvUsecase *CVUsecase) setAvatarPath(ctx context.Context, cvs ...domain.ApiCV) ([]domain.ApiCV, error) {
	cvIDs := []int{}

	for _, cv := range cvs {
		cvIDs = append(cvIDs, cv.ID)
	}

	cvIDToPath, err := cvUsecase.userRepo.GetAvatarPathesByCVIDList(ctx, cvIDs...)
	if err != nil {
		return nil, err
	}

	cvsToReturn := []domain.ApiCV{}
	for _, cv := range cvs {
		path, found := cvIDToPath[cv.ID]
		if !found || path == "" {
			cv.AvatarURL = ""
		} else {
			cv.AvatarURL = "/image" + path
		}

		cvsToReturn = append(cvsToReturn, cv)
	}

	return cvsToReturn, nil
}

// Finds cv that responded to one of the current user's vacancy
func (cvUsecase *CVUsecase) GetCVById(ctx context.Context, cvID int) (*domain.ApiCV, error) {
	sesssionID, err := contextUtils.GetSessionIDFromCtx(ctx)
	if err == nil {
		userID, err := cvUsecase.sessionRepo.GetUserIdBySession(ctx, sesssionID)
		if err != nil {
			return nil, err
		}
		employerID, err := cvUsecase.userRepo.GetUserEmpId(ctx, userID)
		if !errors.Is(err, psql.ErrEntityNotFound) {
			if err != nil {
				return nil, err
			} else {
				err = cvUsecase.statRepo.AddCvView(ctx, cvID, employerID)
				if err != nil && !errors.Is(err, psql.ErrRecordAlredyExists) {
					return nil, err
				}
			}
		}
	}

	cv, exps, edInsts, err := cvUsecase.cvRepo.GetCVById(ctx, cvID)
	if err != nil {
		return nil, err
	}

	apiCV := cvUsecase.constructApiCV(cv, exps, edInsts)

	viewsCount, err := cvUsecase.statRepo.CountCvViews(ctx, cvID)
	if err != nil {
		return nil, err
	}
	apiCV.ViewsCount = viewsCount

	responsesCount, err := cvUsecase.statRepo.CountCvResponses(ctx, cvID)
	if err != nil {
		return nil, err
	}
	apiCV.ResponsesCount = responsesCount

	skills, err := cvUsecase.skillRepo.GetSkillsByCvID(ctx, apiCV.ID)
	if err != nil {
		return nil, err
	}
	apiCV.Skills = skills

	cvToReturn, err := cvUsecase.setAvatarPath(ctx, *apiCV)
	if err != nil {
		return nil, err
	}

	return &cvToReturn[0], nil
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

func (cvUsecase *CVUsecase) GetCVList(ctx context.Context) ([]domain.ApiCV, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)
	validStatus := cvUsecase.validateRole(ctx, userID, domain.Applicant)
	if validStatus != nil {
		return nil, validStatus
	}

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

	apiCvs := cvUsecase.combineDbCVs(cvs, exps, insts)

	for i := range apiCvs {
		viewsCount, err := cvUsecase.statRepo.CountCvViews(ctx, apiCvs[i].ID)
		if err != nil {
			return nil, err
		}
		apiCvs[i].ViewsCount = viewsCount

		responsesCount, err := cvUsecase.statRepo.CountCvResponses(ctx, apiCvs[i].ID)
		if err != nil {
			return nil, err
		}
		apiCvs[i].ResponsesCount = responsesCount
	}

	// TODO: optimize
	for i := range apiCvs {
		skills, err := cvUsecase.skillRepo.GetSkillsByCvID(ctx, apiCvs[i].ID)
		if err != nil {
			return nil, err
		}
		apiCvs[i].Skills = skills
	}

	apiCvs, err = cvUsecase.setAvatarPath(ctx, apiCvs...)
	if err != nil {
		return nil, err
	}

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

func (cvUsecase *CVUsecase) AddNewCV(ctx context.Context, cv *domain.ApiCV) (int, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	dbExperiences, dbEducationInstitutions, dbCV := cvUsecase.getDataFromApiCV(cv)

	cvID, addErr := cvUsecase.cvRepo.AddCV(ctx, userID, dbCV, dbExperiences, dbEducationInstitutions)
	if addErr != nil {
		return 0, addErr
	}

	err := cvUsecase.skillRepo.AddSkillsByCvID(ctx, cvID, cv.Skills)
	if err != nil {
		return 0, err
	}

	return cvID, nil
}

func (cvUsecase *CVUsecase) GetCVOfUserById(ctx context.Context, cvID int) (*domain.ApiCV, error) {
	userID := contextUtils.GetUserIDFromCtx(ctx)

	cv, exps, insts, err := cvUsecase.cvRepo.GetOneOfUsersCV(ctx, userID, cvID)
	if err != nil {
		return nil, err
	}

	apiCv := cvUsecase.constructApiCV(cv, exps, insts)

	viewsCount, err := cvUsecase.statRepo.CountCvViews(ctx, cvID)
	if err != nil {
		return nil, err
	}
	apiCv.ViewsCount = viewsCount

	responsesCount, err := cvUsecase.statRepo.CountCvResponses(ctx, cvID)
	if err != nil {
		return nil, err
	}
	apiCv.ResponsesCount = responsesCount

	skills, err := cvUsecase.skillRepo.GetSkillsByCvID(ctx, apiCv.ID)
	if err != nil {
		return nil, err
	}
	apiCv.Skills = skills

	cvToReturn, err := cvUsecase.setAvatarPath(ctx, *apiCv)
	if err != nil {
		return nil, err
	}

	return &cvToReturn[0], nil
}

func (cvUsecase *CVUsecase) GetApplicantInfo(ctx context.Context, applicantID int) (*domain.ApplicantInfo, error) {
	first_name, last_name, email, cvs, exp, edu, err := cvUsecase.cvRepo.GetApplicantInfo(ctx, applicantID)
	if err != nil {
		return nil, err
	}

	cvsToReturn := cvUsecase.combineDbCVs(cvs, exp, edu)

	infoCVs, err := cvUsecase.setAvatarPath(ctx, cvsToReturn...)
	if err != nil {
		return nil, err
	}

	info := &domain.ApplicantInfo{
		FirstName: first_name,
		LastName:  last_name,
		Email:     email,
		CVs:       infoCVs,
	}

	return info, nil
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

func (cvUsecase *CVUsecase) UpdateCVOfUserById(ctx context.Context, cvID int, cv *domain.ApiCV) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

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
	updSkillsErr := cvUsecase.skillRepo.UpdateSkillsByCvID(ctx, cvID, cv.Skills)
	if updSkillsErr != nil && updSkillsErr != psql.ErrNoRowsUpdated {
		return updSkillsErr
	}

	return nil
}

func (cvUsecase *CVUsecase) DeleteCVOfUserById(ctx context.Context, cvID int) error {
	userID := contextUtils.GetUserIDFromCtx(ctx)

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
	options *searchEnginePB.SearchOptions,
) (domain.ApiMetaCV, error) {
	cvSearchResponse, err := cvUsecase.searchEngineRepo.SearchCVsIDs(ctx, options)
	if err != nil {
		return domain.ApiMetaCV{
			Filters: nil,
			CVs:     domain.ApiCVCount{},
		}, err
	}

	// dbCvs, cvErr := cvUsecase.cvRepo.GetCVsByIds(ctx, castUtils.Int64SliceToIntSlice(cvIDs))
	dbCvs, dbExps, dbInsts, cvErr := cvUsecase.cvRepo.GetCVsByIds(ctx, castUtils.Int64SliceToIntSlice(cvSearchResponse.Ids))
	if cvErr == psql.ErrEntityNotFound {
		return domain.ApiMetaCV{
			Filters: nil,
			CVs:     domain.ApiCVCount{},
		}, nil
	}
	if cvErr != nil {
		return domain.ApiMetaCV{}, cvErr
	}

	cvsToReturn := cvUsecase.combineDbCVs(dbCvs, dbExps, dbInsts)

	for i := range cvsToReturn {
		viewsCount, err := cvUsecase.statRepo.CountCvViews(ctx, cvsToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaCV{}, err
		}
		cvsToReturn[i].ViewsCount = viewsCount

		responsesCount, err := cvUsecase.statRepo.CountCvResponses(ctx, cvsToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaCV{}, err
		}
		cvsToReturn[i].ResponsesCount = responsesCount
	}

	// TODO: optimize
	for i := range cvsToReturn {
		skills, err := cvUsecase.skillRepo.GetSkillsByCvID(ctx, cvsToReturn[i].ID)
		if err != nil {
			return domain.ApiMetaCV{}, err
		}
		cvsToReturn[i].Skills = skills
	}

	cvsToReturn, err = cvUsecase.setAvatarPath(ctx, cvsToReturn...)
	if err != nil {
		return domain.ApiMetaCV{}, err
	}

	result := domain.ApiMetaCV{
		Filters: cvSearchResponse.Filters,
		CVs: domain.ApiCVCount{
			Count: cvSearchResponse.Count,
			CVs:   cvsToReturn,
		},
	}

	return result, nil
}

func (cvUsecase *CVUsecase) GenerateCVsPDF(ctx context.Context, cvID int) (*gofpdf.Fpdf, error) {
	cv, err := cvUsecase.GetCVById(ctx, cvID)
	if err != nil {
		return nil, err
	}

	pdfFile, err := pdf.MarshalPDF(&pdf.CVConfig, cv)
	if err != nil {
		return nil, err
	}

	return pdfFile, nil
}
