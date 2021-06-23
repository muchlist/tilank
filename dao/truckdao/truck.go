package truckdao

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"tilank/db"
	"tilank/dto"
	"tilank/utils/logger"
	"tilank/utils/rest_err"
	"time"
)

const (
	connectTimeout     = 3
	keyTruckCollection = "truck"

	keyTruckID          = "_id"
	keyTruckUpdatedAt   = "updated_at"
	keyTruckUpdatedBy   = "updated_by"
	keyTruckUpdatedByID = "updated_by_id"
	keyTruckBranch      = "branch"
	keyNoIdentity       = "no_identity"
	keyNoPol            = "no_pol"
	keyMark             = "mark"
	keyOwner            = "owner"
	keyEmail            = "email"
	keyHp               = "hp"
	keyDeleted          = "deleted"
	keyScore            = "score"
	keyResetScoreDate   = "reset_score_date"
	keyBlocked          = "blocked"
	keyBlockStart       = "block_start"
	keyBlockEnd         = "block_end"
)

func NewTruckDao() TruckDaoAssumer {
	return &truckDao{}
}

type truckDao struct {
}

type TruckDaoAssumer interface {
	InsertTruck(input dto.Truck) (*string, resterr.APIError)
	EditTruck(input dto.TruckEdit) (*dto.Truck, resterr.APIError)
	DeleteTruck(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Truck, resterr.APIError)
	ChangeScore(input dto.TruckScoreEdit) (*dto.Truck, resterr.APIError)

	GetTruckByID(truckID primitive.ObjectID, branchIfSpecific string) (*dto.Truck, resterr.APIError)
	FindTruck(filter dto.FilterTruck) (dto.TruckResponseMinList, resterr.APIError)
}

func (c *truckDao) InsertTruck(input dto.Truck) (*string, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Owner = strings.ToUpper(input.Owner)
	input.Branch = strings.ToUpper(input.Branch)
	input.Email = strings.ToLower(input.Email)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := resterr.NewInternalServerError("Gagal menyimpan truck ke database", err)
		logger.Error("Gagal menyimpan truck ke database, (InsertTruck)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *truckDao) EditTruck(input dto.TruckEdit) (*dto.Truck, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.FilterBranch = strings.ToUpper(input.FilterBranch)
	input.Owner = strings.ToUpper(input.Owner)
	input.Email = strings.ToLower(input.Email)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyTruckID:        input.ID,
		keyTruckBranch:    input.FilterBranch,
		keyTruckUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyTruckUpdatedAt:   input.UpdatedAt,
			keyTruckUpdatedBy:   input.UpdatedBy,
			keyTruckUpdatedByID: input.UpdatedByID,
			keyNoIdentity:       input.NoIdentity,
			keyNoPol:            input.NoPol,
			keyMark:             input.Mark,
			keyOwner:            input.Owner,
			keyEmail:            input.Email,
			keyHp:               input.Hp,
		},
	}

	var truck dto.Truck
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&truck); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("truck tidak diupdate : validasi id timestamp")
		}

		logger.Error("Gagal mendapatkan truck dari database (EditTruck)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan truck dari database", err)
		return nil, apiErr
	}

	return &truck, nil
}

func (c *truckDao) ChangeScore(input dto.TruckScoreEdit) (*dto.Truck, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyTruckID: input.ID,
	}

	update := bson.M{
		"$set": bson.M{
			keyScore:          input.Score,
			keyResetScoreDate: input.ResetScoreDate,
			keyBlocked:        input.Blocked,
			keyBlockStart:     input.BlockStart,
			keyBlockEnd:       input.BlockEnd,
		},
	}

	var truck dto.Truck
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&truck); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("truck tidak diupdate : truck dengan id tersebut tidak ditemukan")
		}

		logger.Error("Gagal mendapatkan truck dari database (ChangeScore)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan truck dari database", err)
		return nil, apiErr
	}

	return &truck, nil
}

func (c *truckDao) DeleteTruck(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Truck, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyTruckID:     input.FilterID,
		keyTruckBranch: input.FilterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyDeleted: isSoftDelete,
		},
	}

	var truck dto.Truck
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&truck)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("Truck tidak dihapus : validasi id branch status")
		}

		logger.Error("Gagal menghapus truck dari database (DeleteTruck)", err)
		apiErr := resterr.NewInternalServerError("Gagal menghapus truck dari database", err)
		return nil, apiErr
	}

	return &truck, nil
}

func (c *truckDao) GetTruckByID(truckID primitive.ObjectID, branchIfSpecific string) (*dto.Truck, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyTruckID: truckID}
	if branchIfSpecific != "" {
		filter[keyTruckBranch] = strings.ToUpper(branchIfSpecific)
	}

	var truck dto.Truck
	if err := coll.FindOne(ctx, filter).Decode(&truck); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("Truck dengan ID %s tidak ditemukan", truckID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan truck dari database (GetTruckByID)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan truck dari database", err)
		return nil, apiErr
	}

	return &truck, nil
}

func (c *truckDao) FindTruck(filterA dto.FilterTruck) (dto.TruckResponseMinList, resterr.APIError) {
	coll := db.DB.Collection(keyTruckCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterNoIdentity = strings.ToUpper(filterA.FilterNoIdentity)
	filterA.FilterOwner = strings.ToUpper(filterA.FilterOwner)

	// filter
	filter := bson.M{
		keyDeleted: !filterA.Active,
	}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyTruckBranch] = filterA.FilterBranch
	}
	if filterA.FilterNoIdentity != "" {
		filter[keyNoIdentity] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterNoIdentity),
		}
	}
	if filterA.FilterOwner != "" {
		filter[keyOwner] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterOwner),
		}
	}
	if filterA.Blocked {
		filter[keyBlocked] = true
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyOwner, 1}, {keyNoIdentity, 1}}) //nolint:govet

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar truck dari database (FindTruck)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.TruckResponseMinList{}, apiErr
	}

	truckList := dto.TruckResponseMinList{}
	if err = cursor.All(ctx, &truckList); err != nil {
		logger.Error("Gagal decode truckList cursor ke objek slice (FindTruck)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.TruckResponseMinList{}, apiErr
	}

	return truckList, nil
}
