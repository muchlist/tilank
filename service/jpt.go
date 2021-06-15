package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/jptdao"
	"tilank/dto"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewJptService(jptDao jptdao.JptDaoAssumer) *JptService {
	return &JptService{
		daoC: jptDao,
	}
}

type JptService struct {
	daoC jptdao.JptDaoAssumer
}

func (j *JptService) InsertJpt(user mjwt.CustomClaim, input dto.JptRequest) (*string, resterr.APIError) {
	idGenerated := primitive.NewObjectID()

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Jpt{
		ID:          idGenerated,
		CreatedAt:   timeNow,
		CreatedBy:   user.Name,
		CreatedByID: user.Identity,
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		Branch:      user.Branch,
		Name:        input.Name,
		OwnerName:   input.OwnerName,
		IDPelindo:   input.IDPelindo,
		Hp:          input.Hp,
		Email:       input.Email,
		Deleted:     false,
	}

	// DB
	insertedID, err := j.daoC.InsertJpt(data)
	if err != nil {
		return nil, resterr.NewBadRequestError(err.Message())
	}

	return insertedID, nil
}

func (j *JptService) EditJpt(user mjwt.CustomClaim, jptID string, input dto.JptEditRequest) (*dto.Jpt, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(jptID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// Filling data
	timeNow := time.Now().Unix()
	data := dto.JptEdit{
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Name:            input.Name,
		OwnerName:       input.OwnerName,
		IDPelindo:       input.IDPelindo,
		Hp:              input.Hp,
		Email:           input.Email,
	}

	// DB
	jptEdited, err := j.daoC.EditJpt(data)
	if err != nil {
		return nil, err
	}

	return jptEdited, nil
}

func (j *JptService) DeleteJpt(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := j.daoC.DeleteJpt(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (j *JptService) ActivateJpt(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := j.daoC.DeleteJpt(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}, false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JptService) GetJptByID(jptID string, branchIfSpecific string) (*dto.Jpt, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(jptID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	jpt, err := j.daoC.GetJptByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return jpt, nil
}

func (j *JptService) FindJpt(filter dto.FilterJpt) (dto.JptResponseMinList, resterr.APIError) {

	jptList, err := j.daoC.FindJpt(filter)
	if err != nil {
		return nil, err
	}

	return jptList, nil
}
