package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/violationdao"
	"tilank/dto"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewViolationService(violationDao violationdao.ViolationDaoAssumer) *violationService {
	return &violationService{
		daoC: violationDao,
	}
}

type violationService struct {
	daoC violationdao.ViolationDaoAssumer
}

func (c *violationService) InsertViolation(user mjwt.CustomClaim, input dto.ViolationRequest) (*string, resterr.APIError) {

	idGenerated := primitive.NewObjectID()

	// todo get owner

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Violation{
		ID:           idGenerated,
		CreatedAt:    timeNow,
		CreatedBy:    user.Name,
		CreatedByID:  user.Identity,
		UpdatedAt:    timeNow,
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
		ApprovedAt:   0,
		ApprovedBy:   "",
		ApprovedByID: "",
		Branch:       user.Branch,
		State:        input.State,
		NoIdentity:   input.NoIdentity,
		NoPol:        input.NoPol,
		Mark:         input.Mark,
		//Owner:           input.Owner,
		//OwnerID:         input.OwnerID,
		TypeViolation:   input.TypeViolation,
		DetailViolation: input.DetailViolation,
		TimeViolation:   input.TimeViolation,
		Location:        input.Location,
		Images:          []string{},
	}

	// DB
	insertedID, err := c.daoC.InsertViolation(data)
	if err != nil {
		return nil, resterr.NewBadRequestError(err.Message())
	}

	return insertedID, nil
}

func (c *violationService) EditViolation(user mjwt.CustomClaim, violationID string, input dto.ViolationEditRequest) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// todo get owner

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.ViolationEdit{
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		ApprovedAt:      0,
		ApprovedBy:      "",
		ApprovedByID:    "",
		NoIdentity:      input.NoIdentity,
		NoPol:           input.NoPol,
		Mark:            input.Mark,
		//Owner:           input.Owner,
		//OwnerID:         input.OwnerID,
		TypeViolation:   input.TypeViolation,
		DetailViolation: input.DetailViolation,
		TimeViolation:   input.TimeViolation,
		Location:        input.Location,
	}

	// DB
	violationEdited, err := c.daoC.EditViolation(data)
	if err != nil {
		return nil, err
	}

	return violationEdited, nil
}

func (c *violationService) ApproveViolation(user mjwt.CustomClaim, violationID string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.ViolationConfirm{
		ID:           oid,
		FilterBranch: user.Branch,
		UpdatedAt:    timeNow,
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
		ApprovedAt:   timeNow,
		ApprovedBy:   user.Name,
		ApprovedByID: user.Identity,
		State:        2,
	}

	// DB
	violationApproved, err := c.daoC.ConfirmViolation(data)
	if err != nil {
		return nil, err
	}

	return violationApproved, nil
}

func (c *violationService) DeleteViolation(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := c.daoC.DeleteViolation(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	})
	if err != nil {
		return err
	}

	return nil
}

// PutImage memasukkan lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (c *violationService) PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := c.daoC.UploadImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return violation, nil
}

// DeleteImage menghapus lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (c *violationService) DeleteImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := c.daoC.DeleteImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return violation, nil
}

func (c *violationService) GetViolationByID(violationID string, branchIfSpecific string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := c.daoC.GetViolationByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return violation, nil
}

func (c *violationService) FindViolation(filter dto.FilterViolation) (dto.ViolationResponseMinList, resterr.APIError) {
	// jika filter nomer identitas ada maka filter nomor polisinya dihilangkan
	if filter.FilterNoIdentity != "" {
		filter.FilterNoPol = ""
	}
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	violationList, err := c.daoC.FindViolation(filter)
	if err != nil {
		return nil, err
	}

	return violationList, nil
}
