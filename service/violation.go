package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/jptdao"
	"tilank/dao/violationdao"
	"tilank/dto"
	"tilank/enum"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewViolationService(violationDao violationdao.ViolationDaoAssumer, jptDao jptdao.JptDaoAssumer) *ViolationService {
	return &ViolationService{
		vDao: violationDao,
		jDao: jptDao,
	}
}

type ViolationService struct {
	vDao violationdao.ViolationDaoAssumer
	jDao jptdao.JptDaoAssumer
}

func getJPTName(jDao jptdao.JptDaoAssumer, ownerID string) (string, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(ownerID)
	if errT != nil {
		return "", resterr.NewBadRequestError("Owner ID yang dimasukkan salah")
	}

	jpt, err := jDao.GetJptByID(oid, "")
	if err != nil {
		return "", err
	}

	return jpt.Name, nil
}

func (v *ViolationService) InsertViolation(user mjwt.CustomClaim, input dto.ViolationRequest) (*string, resterr.APIError) {
	idGenerated := primitive.NewObjectID()

	// validasi inputan state,
	state := input.State
	// jika bukan 0 draft atau 1 perlu persetujuan set ke draft
	if !(state == enum.StUndefined || state == enum.StNeedApprove) {
		state = enum.StUndefined
	}

	// mendapatkan nama jpt
	jptName, err := getJPTName(v.jDao, input.OwnerID)
	if err != nil {
		return nil, err
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Violation{
		ID:              idGenerated,
		CreatedAt:       timeNow,
		CreatedBy:       user.Name,
		CreatedByID:     user.Identity,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		ApprovedAt:      0,
		ApprovedBy:      "",
		ApprovedByID:    "",
		Branch:          user.Branch,
		State:           state,
		NoIdentity:      input.NoIdentity,
		NoPol:           input.NoPol,
		Mark:            input.Mark,
		Owner:           jptName,
		OwnerID:         input.OwnerID,
		TypeViolation:   input.TypeViolation,
		DetailViolation: input.DetailViolation,
		TimeViolation:   input.TimeViolation,
		Location:        input.Location,
		Images:          []string{},
	}

	// DB
	insertedID, err := v.vDao.InsertViolation(data)
	if err != nil {
		return nil, resterr.NewBadRequestError(err.Message())
	}

	return insertedID, nil
}

func (v *ViolationService) EditViolation(user mjwt.CustomClaim, violationID string, input dto.ViolationEditRequest) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	jptOwnerName, err := getJPTName(v.jDao, input.OwnerID)
	if err != nil {
		return nil, err
	}

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
		Owner:           jptOwnerName,
		OwnerID:         input.OwnerID,
		TypeViolation:   input.TypeViolation,
		DetailViolation: input.DetailViolation,
		TimeViolation:   input.TimeViolation,
		Location:        input.Location,
	}

	// DB
	violationEdited, err := v.vDao.EditViolation(data)
	if err != nil {
		return nil, err
	}

	return violationEdited, nil
}

func (v *ViolationService) SendToDraftViolation(user mjwt.CustomClaim, violationID string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// cek dokumen eksisting
	violation, err := v.vDao.GetViolationByID(oid, "")
	if err != nil {
		return nil, err
	}

	// validasi
	if violation.State != enum.StNeedApprove {
		apiErr := resterr.NewBadRequestError("status dokumen tidak dapat diubah ke draft")
		return nil, apiErr
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.ViolationConfirm{
		ID:           oid,
		FilterBranch: user.Branch,
		UpdatedAt:    timeNow,
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
		ApprovedAt:   0,
		ApprovedBy:   "",
		ApprovedByID: "",
		State:        enum.StDraft,
	}

	// DB
	violationDrafted, err := v.vDao.ChangeStateViolation(data)
	if err != nil {
		return nil, err
	}

	return violationDrafted, nil
}

func (v *ViolationService) SendToConfirmationViolation(user mjwt.CustomClaim, violationID string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// cek dokumen eksisting
	violation, err := v.vDao.GetViolationByID(oid, "")
	if err != nil {
		return nil, err
	}

	// validasi
	if violation.State != enum.StDraft {
		apiErr := resterr.NewBadRequestError("status dokumen tidak dapat diubah ke NeedConfirm")
		return nil, apiErr
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.ViolationConfirm{
		ID:           oid,
		FilterBranch: user.Branch,
		UpdatedAt:    timeNow,
		UpdatedBy:    user.Name,
		UpdatedByID:  user.Identity,
		ApprovedAt:   0,
		ApprovedBy:   "",
		ApprovedByID: "",
		State:        enum.StNeedApprove,
	}

	// DB
	violationDrafted, err := v.vDao.ChangeStateViolation(data)
	if err != nil {
		return nil, err
	}

	return violationDrafted, nil
}

func (v *ViolationService) ApproveViolation(user mjwt.CustomClaim, violationID string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// cek dokumen eksisting
	violation, err := v.vDao.GetViolationByID(oid, "")
	if err != nil {
		return nil, err
	}

	// validasi
	if violation.State != enum.StNeedApprove {
		apiErr := resterr.NewBadRequestError("status dokumen tidak dapat di approve")
		return nil, apiErr
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
		State:        enum.StApproved,
	}

	// DB
	violationApproved, err := v.vDao.ChangeStateViolation(data)
	if err != nil {
		return nil, err
	}

	return violationApproved, nil
}

func (v *ViolationService) DeleteViolation(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := v.vDao.DeleteViolation(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	})
	if err != nil {
		return err
	}

	return nil
}

// PutImage memasukkan lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (v *ViolationService) PutImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := v.vDao.UploadImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return violation, nil
}

// DeleteImage menghapus lokasi file (path) ke dalam database violation dengan mengecek kesesuaian branch
func (v *ViolationService) DeleteImage(user mjwt.CustomClaim, id string, imagePath string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := v.vDao.DeleteImage(oid, imagePath, user.Branch)
	if err != nil {
		return nil, err
	}
	return violation, nil
}

func (v *ViolationService) GetViolationByID(violationID string, branchIfSpecific string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	violation, err := v.vDao.GetViolationByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return violation, nil
}

func (v *ViolationService) FindViolation(filter dto.FilterViolation) (dto.ViolationResponseMinList, resterr.APIError) {
	// jika filter nomer identitas ada maka filter nomor polisinya dihilangkan
	if filter.FilterNoIdentity != "" {
		filter.FilterNoPol = ""
	}
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	violationList, err := v.vDao.FindViolation(filter)
	if err != nil {
		return nil, err
	}

	return violationList, nil
}
