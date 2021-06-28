package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/dao/truckdao"
	"tilank/dto"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewTruckService(truckDao truckdao.TruckDaoAssumer) *TruckService {
	return &TruckService{
		daoC: truckDao,
	}
}

type TruckService struct {
	daoC truckdao.TruckDaoAssumer
}

func (j *TruckService) InsertTruck(user mjwt.CustomClaim, input dto.TruckRequest) (*string, resterr.APIError) {
	idGenerated := primitive.NewObjectID()

	truckExisting, _ := j.daoC.GetTruckByIdentity(input.NoIdentity, user.Branch)
	if truckExisting != nil {
		return nil, resterr.NewBadRequestError("Nomor lambung tidak tersedia! ")
	}

	// Filling data
	timeNow := time.Now().Unix()
	data := dto.Truck{
		ID:             idGenerated,
		CreatedAt:      timeNow,
		CreatedBy:      user.Name,
		CreatedByID:    user.Identity,
		UpdatedAt:      timeNow,
		UpdatedBy:      user.Name,
		UpdatedByID:    user.Identity,
		Branch:         user.Branch,
		NoIdentity:     input.NoIdentity,
		NoPol:          input.NoPol,
		Mark:           input.Mark,
		Owner:          input.Owner,
		Email:          input.Email,
		Hp:             input.Hp,
		Deleted:        false,
		Score:          0,
		ResetScoreDate: 0,
		Blocked:        false,
		BlockStart:     0,
		BlockEnd:       0,
	}

	// DB
	insertedID, err := j.daoC.InsertTruck(data)
	if err != nil {
		return nil, resterr.NewBadRequestError(err.Message())
	}

	return insertedID, nil
}

func (j *TruckService) EditTruck(user mjwt.CustomClaim, truckID string, input dto.TruckEditRequest) (*dto.Truck, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(truckID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// Filling data
	timeNow := time.Now().Unix()
	data := dto.TruckEdit{
		ID:              oid,
		FilterBranch:    user.Branch,
		FilterTimestamp: input.FilterTimestamp,
		UpdatedAt:       timeNow,
		UpdatedBy:       user.Name,
		UpdatedByID:     user.Identity,
		NoIdentity:      input.NoIdentity,
		NoPol:           input.NoPol,
		Mark:            input.Mark,
		Owner:           input.Owner,
		Email:           input.Email,
		Hp:              input.Hp,
	}

	// DB
	truckEdited, err := j.daoC.EditTruck(data)
	if err != nil {
		return nil, err
	}

	return truckEdited, nil
}

func (j *TruckService) DeleteTruck(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := j.daoC.DeleteTruck(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func (j *TruckService) ActivateTruck(user mjwt.CustomClaim, id string) resterr.APIError {
	oid, errT := primitive.ObjectIDFromHex(id)
	if errT != nil {
		return resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}
	// DB
	_, err := j.daoC.DeleteTruck(dto.FilterIDBranch{
		FilterID:     oid,
		FilterBranch: user.Branch,
	}, false)
	if err != nil {
		return err
	}

	return nil
}

func (j *TruckService) GetTruckByID(truckID string, branchIfSpecific string) (*dto.Truck, resterr.APIError) {
	oid, errT := primitive.ObjectIDFromHex(truckID)
	if errT != nil {
		return nil, resterr.NewBadRequestError("ObjectID yang dimasukkan salah")
	}

	truck, err := j.daoC.GetTruckByID(oid, branchIfSpecific)
	if err != nil {
		return nil, err
	}

	return truck, nil
}

func (j *TruckService) FindTruck(filter dto.FilterTruck) (dto.TruckResponseMinList, resterr.APIError) {
	truckList, err := j.daoC.FindTruck(filter)
	if err != nil {
		return nil, err
	}

	return truckList, nil
}
