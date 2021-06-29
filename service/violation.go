package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/rulesdao"
	"tilank/dao/truckdao"
	"tilank/dao/violationdao"
	"tilank/dto"
	"tilank/enum"
	"tilank/utils/logger"
	"tilank/utils/mjwt"
	"tilank/utils/pdfgen"
	"tilank/utils/rest_err"
	"time"
)

func NewViolationService(violationDao violationdao.ViolationDaoAssumer,
	truckDao truckdao.TruckDaoAssumer,
	rulesDao rulesdao.RulesDaoAssumer) *ViolationService {
	return &ViolationService{
		vDao: violationDao,
		tDao: truckDao,
		rDao: rulesDao,
	}
}

type ViolationService struct {
	vDao violationdao.ViolationDaoAssumer
	tDao truckdao.TruckDaoAssumer
	rDao rulesdao.RulesDaoAssumer
}

func (v *ViolationService) InsertViolation(user mjwt.CustomClaim, input dto.ViolationRequest) (*string, resterr.APIError) {
	idGenerated := primitive.NewObjectID()

	// validasi inputan state,
	state := input.State
	// jika bukan 0 draft atau 1 perlu persetujuan set ke draft
	if !(state == enum.StDraft || state == enum.StNeedApprove) {
		state = enum.StDraft
	}

	if input.TimeViolation == 0 {
		input.TimeViolation = time.Now().Unix()
	}

	// mendapatkan truck
	truck, err := v.tDao.GetTruckByIdentity(input.NoIdentity, user.Branch)
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
		NViol:           truck.Score,
		NoIdentity:      truck.NoIdentity,
		NoPol:           truck.NoPol,
		Mark:            truck.Mark,
		Owner:           truck.Owner,
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

	truck, err := v.tDao.GetTruckByIdentity(input.NoIdentity, user.Branch)
	if err != nil {
		return nil, err
	}

	if input.TimeViolation == 0 {
		input.TimeViolation = time.Now().Unix()
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
		NoIdentity:      truck.NoIdentity,
		NoPol:           truck.NoPol,
		Mark:            truck.Mark,
		Owner:           truck.Owner,
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
		NViol:        violation.NViol,
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
	if !(violation.State == enum.StDraft || violation.State == enum.StUndefined) {
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
		NViol:        violation.NViol,
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

	// 1 cek dokumen eksisting
	violation, err := v.vDao.GetViolationByID(oid, "")
	if err != nil {
		return nil, err
	}

	// 2 validasi
	if violation.State != enum.StNeedApprove {
		apiErr := resterr.NewBadRequestError("status dokumen tidak dapat di approve")
		return nil, apiErr
	}

	// 3 mendapatkan data truck eksisting
	truck, err := v.tDao.GetTruckByIdentity(violation.NoIdentity, violation.Branch)
	if err != nil {
		logger.
			Error(fmt.Sprintf(
				"error di ApproveViolation 5 mendapatkan data truck eksisting, need roleback. id : %s",
				violationID),
				err)
		return nil, err
	}

	// 4 Filling data
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
		NViol:        truck.Score + 1,
	}

	// 5 DB
	violationApproved, err := v.vDao.ChangeStateViolation(data)
	if err != nil {
		return nil, err
	}

	payloadTruck := dto.TruckScoreEdit{
		ID:             truck.ID,
		Score:          truck.Score + 1,
		ResetScoreDate: 0,
		Blocked:        false,
		BlockStart:     0,
		BlockEnd:       0,
	}

	// 6 mendapatkan rules block truck
	rules, _ := v.rDao.GetRulesByScore(payloadTruck.Score, truck.Branch)
	if rules != nil {
		timeNow = time.Now().Unix()
		// jika ditemukan role pemblokiran
		if rules.BlockTime != 0 {
			payloadTruck.Blocked = true
			payloadTruck.BlockStart = timeNow
			payloadTruck.BlockEnd = timeNow + rules.BlockTime
		}
	}

	// 7 menambahkan status di truck
	truckUpdated, err := v.tDao.ChangeScore(payloadTruck)
	if err != nil {
		logger.
			Error(fmt.Sprintf(
				"error di ApproveViolation 7 menambahkan status di truck, need roleback. id : %s",
				violationID),
				err)
		return nil, err
	}

	// 8 membuat pdf
	errPdf := pdfgen.GeneratePDF(violationApproved, truckUpdated, rules)
	if errPdf != nil {
		logger.Error(fmt.Sprintf("membuat pdf gagal. id : %s", violationID), errPdf)
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

func (v *ViolationService) GeneratePDFViolation(violationID string) (*dto.Violation, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(violationID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	// 1 mendapatkan dokumen pelanggaran
	violation, err := v.vDao.GetViolationByID(oid, "")
	if err != nil {
		return nil, err
	}

	// 2 mendapatkan data pelanggaran ini adalah pelanggaran keberapa
	violations, err := v.vDao.FindViolation(dto.FilterViolation{
		FilterBranch:     violation.Branch,
		FilterNoIdentity: violation.NoIdentity,
		FilterState:      2,
	})
	if err != nil {
		return nil, err
	}

	violationScore := 0
	for i, v := range violations {
		if v.ID == violation.ID {
			// karena sort violation terbalik yang terbaru (index 0) adalah
			// pelanggaran dengan skor lebih tinggi
			violationScore = len(violations) - i
		}
	}
	// 3 mendapatkan rules block truck
	rules, _ := v.rDao.GetRulesByScore(violationScore, violation.Branch)

	// 4 dummy truck, pdf hanya melihat score saja
	truck := &dto.Truck{
		Score: violationScore,
	}

	errPdf := pdfgen.GeneratePDF(violation, truck, rules)
	if errPdf != nil {
		logger.Error("membuat pdf gagal", errPdf)
	}

	return violation, nil
}
