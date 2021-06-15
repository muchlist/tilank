package jptdao

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
	connectTimeout   = 3
	keyJptCollection = "jpt"

	keyJptID          = "_id"
	keyJptUpdatedAt   = "updated_at"
	keyJptUpdatedBy   = "updated_by"
	keyJptUpdatedByID = "updated_by_id"
	keyJptBranch      = "branch"
	keyJptName        = "name"
	keyJptOwnerName   = "owner_name"
	keyJptIDPelindo   = "id_pelindo"
	keyJptHp          = "hp"
	keyJptEmail       = "email"
	keyJptDeleted     = "deleted"
)

func NewJptDao() JptDaoAssumer {
	return &jptDao{}
}

type jptDao struct {
}

type JptDaoAssumer interface {
	InsertJpt(input dto.Jpt) (*string, resterr.APIError)
	EditJpt(input dto.JptEdit) (*dto.Jpt, resterr.APIError)
	DeleteJpt(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Jpt, resterr.APIError)

	GetJptByID(jptID primitive.ObjectID, branchIfSpecific string) (*dto.Jpt, resterr.APIError)
	FindJpt(filter dto.FilterJpt) (dto.JptResponseMinList, resterr.APIError)
}

func (c *jptDao) InsertJpt(input dto.Jpt) (*string, resterr.APIError) {
	coll := db.DB.Collection(keyJptCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.OwnerName = strings.ToUpper(input.OwnerName)
	input.Branch = strings.ToUpper(input.Branch)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := resterr.NewInternalServerError("Gagal menyimpan jpt ke database", err)
		logger.Error("Gagal menyimpan jpt ke database, (InsertJpt)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *jptDao) EditJpt(input dto.JptEdit) (*dto.Jpt, resterr.APIError) {
	coll := db.DB.Collection(keyJptCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Name = strings.ToUpper(input.Name)
	input.OwnerName = strings.ToUpper(input.OwnerName)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyJptID:        input.ID,
		keyJptBranch:    input.FilterBranch,
		keyJptUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyJptUpdatedAt:   input.UpdatedAt,
			keyJptUpdatedBy:   input.UpdatedBy,
			keyJptUpdatedByID: input.UpdatedByID,
			keyJptName:        input.Name,
			keyJptOwnerName:   input.OwnerName,
			keyJptIDPelindo:   input.IDPelindo,
			keyJptHp:          input.Hp,
			keyJptEmail:       input.Email,
		},
	}

	var jpt dto.Jpt
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&jpt); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("jpt tidak diupdate : validasi id timestamp")
		}

		logger.Error("Gagal mendapatkan jpt dari database (EditJpt)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan jpt dari database", err)
		return nil, apiErr
	}

	return &jpt, nil
}

func (c *jptDao) DeleteJpt(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Jpt, resterr.APIError) {
	coll := db.DB.Collection(keyJptCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyJptID:     input.FilterID,
		keyJptBranch: input.FilterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyJptDeleted: isSoftDelete,
		},
	}

	var jpt dto.Jpt
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&jpt)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("Jpt tidak dihapus : validasi id branch status")
		}

		logger.Error("Gagal menghapus jpt dari database (DeleteJpt)", err)
		apiErr := resterr.NewInternalServerError("Gagal menghapus jpt dari database", err)
		return nil, apiErr
	}

	return &jpt, nil
}

func (c *jptDao) GetJptByID(jptID primitive.ObjectID, branchIfSpecific string) (*dto.Jpt, resterr.APIError) {
	coll := db.DB.Collection(keyJptCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyJptID: jptID}
	if branchIfSpecific != "" {
		filter[keyJptBranch] = strings.ToUpper(branchIfSpecific)
	}

	var jpt dto.Jpt
	if err := coll.FindOne(ctx, filter).Decode(&jpt); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("Jpt dengan ID %s tidak ditemukan", jptID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan jpt dari database (GetJptByID)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan jpt dari database", err)
		return nil, apiErr
	}

	return &jpt, nil
}

func (c *jptDao) FindJpt(filterA dto.FilterJpt) (dto.JptResponseMinList, resterr.APIError) {
	coll := db.DB.Collection(keyJptCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterName = strings.ToUpper(filterA.FilterName)

	// filter
	filter := bson.M{}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyJptBranch] = filterA.FilterBranch
	}
	if filterA.FilterName != "" {
		filter[keyJptName] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterName),
		}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyJptName, 1}}) //nolint:govet

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar jpt dari database (FindJpt)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.JptResponseMinList{}, apiErr
	}

	jptList := dto.JptResponseMinList{}
	if err = cursor.All(ctx, &jptList); err != nil {
		logger.Error("Gagal decode jptList cursor ke objek slice (FindJpt)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.JptResponseMinList{}, apiErr
	}

	return jptList, nil
}
