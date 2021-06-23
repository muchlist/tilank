package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/rulesdao"
	"tilank/dto"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewRulesService(rulesDao rulesdao.RulesDaoAssumer) *RulesService {
	return &RulesService{
		daoC: rulesDao,
	}
}

type RulesService struct {
	daoC rulesdao.RulesDaoAssumer
}

func (j *RulesService) InsertRules(user mjwt.CustomClaim, input dto.RulesRequest) (*string, resterr.APIError) {
	idGenerated := primitive.NewObjectID()

	// DB 2
	rules, _ := j.daoC.GetRulesByScore(input.Score, user.Branch)
	if rules != nil {
		return nil, resterr.NewBadRequestError("Score tersebut sudah ada, silahkan lakukan perubahan di menu edit!")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Rules{
		ID:          idGenerated,
		UpdatedAt:   timeNow,
		UpdatedBy:   user.Name,
		UpdatedByID: user.Identity,
		Branch:      user.Branch,
		Score:       input.Score,
		BlockTime:   input.BlockTime,
		Description: input.Description,
	}

	// DB 2
	insertedID, err := j.daoC.InsertRules(data)
	if err != nil {
		return nil, resterr.NewBadRequestError(err.Message())
	}

	return insertedID, nil
}

func (j *RulesService) EditRules(user mjwt.CustomClaim, rulesID string, input dto.RulesEditRequest) (*dto.Rules, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(rulesID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// Filling data
	timeNow := time.Now().Unix()
	data := dto.RulesEdit{
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		Score:           input.Score,
		BlockTime:       input.BlockTime,
		Description:     input.Description,
	}

	// DB
	rulesEdited, err := j.daoC.EditRules(data)
	if err != nil {
		return nil, err
	}

	return rulesEdited, nil
}

func (j *RulesService) DeleteRules(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := j.daoC.DeleteRules(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (j *RulesService) GetRulesByID(rulesID string, branchIfSpecific string) (*dto.Rules, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(rulesID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	rules, err := j.daoC.GetRulesByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (j *RulesService) FindRules() ([]dto.Rules, resterr.APIError) {

	rulesList, err := j.daoC.FindRules()
	if err != nil {
		return nil, err
	}

	return rulesList, nil
}
