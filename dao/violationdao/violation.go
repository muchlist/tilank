package violationdao

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
	connectTimeout    = 3
	keyViolCollection = "violation"

	keyViolID              = "_id"
	keyViolUpdatedAt       = "updated_at"
	keyViolUpdatedBy       = "updated_by"
	keyViolUpdatedByID     = "updated_by_id"
	keyViolApprovedAt      = "approved_at"
	keyViolApprovedBy      = "approved_by"
	keyViolApprovedByID    = "approved_by_id"
	keyViolBranch          = "branch"
	keyViolState           = "state"
	keyViolNoIdentity      = "no_identity"
	keyViolNoPol           = "no_pol"
	keyViolMark            = "mark"
	keyViolOwner           = "owner"
	keyViolOwnerID         = "owner_id"
	keyViolTypeViolation   = "type_violation"
	keyViolDetailViolation = "detail_violation"
	keyViolTimeViolation   = "time_violation"
	keyViolLocation        = "location"
	keyViolImages          = "images"
)

func NewViolationDao() ViolationDaoAssumer {
	return &violationDao{}
}

type violationDao struct {
}

type ViolationDaoAssumer interface {
	InsertViolation(input dto.Violation) (*string, resterr.APIError)
	EditViolation(input dto.ViolationEdit) (*dto.Violation, resterr.APIError)
	DeleteViolation(input dto.FilterIDBranch) (*dto.Violation, resterr.APIError)
	UploadImage(violationID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Violation, resterr.APIError)
	DeleteImage(violationID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Violation, resterr.APIError)
	ConfirmViolation(input dto.ViolationConfirm) (*dto.Violation, resterr.APIError)

	GetViolationByID(violationID primitive.ObjectID, branchIfSpecific string) (*dto.Violation, resterr.APIError)
	FindViolation(filter dto.FilterViolation) (dto.ViolationResponseMinList, resterr.APIError)
}

func (c *violationDao) InsertViolation(input dto.Violation) (*string, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.NoIdentity = strings.ToUpper(input.NoIdentity)
	input.NoPol = strings.ToUpper(input.NoPol)
	input.Branch = strings.ToUpper(input.Branch)
	if input.Images == nil {
		input.Images = []string{}
	}

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := resterr.NewInternalServerError("Gagal menyimpan pelanggaran ke database", err)
		logger.Error("Gagal menyimpan violation ke database, (InsertViolation)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *violationDao) EditViolation(input dto.ViolationEdit) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.NoIdentity = strings.ToUpper(input.NoIdentity)
	input.NoPol = strings.ToUpper(input.NoPol)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyViolID:        input.ID,
		keyViolBranch:    input.FilterBranch,
		keyViolUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyViolUpdatedAt:   input.UpdatedAt,
			keyViolUpdatedBy:   input.UpdatedBy,
			keyViolUpdatedByID: input.UpdatedByID,

			keyViolNoIdentity:      input.NoIdentity,
			keyViolNoPol:           input.NoPol,
			keyViolMark:            input.Mark,
			keyViolOwner:           input.Owner,
			keyViolOwnerID:         input.OwnerID,
			keyViolTypeViolation:   input.TypeViolation,
			keyViolDetailViolation: input.DetailViolation,
			keyViolTimeViolation:   input.TimeViolation,
			keyViolLocation:        input.Location,
		},
	}

	var violation dto.Violation
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("pelanggaran tidak diupdate : validasi id timestamp")
		}

		logger.Error("Gagal mendapatkan violation dari database (EditViolation)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan violation dari database", err)
		return nil, apiErr
	}

	return &violation, nil
}

func (c *violationDao) ConfirmViolation(input dto.ViolationConfirm) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyViolID:     input.ID,
		keyViolBranch: input.FilterBranch,
	}

	update := bson.M{
		"$set": bson.M{
			keyViolUpdatedAt:   input.UpdatedAt,
			keyViolUpdatedBy:   input.UpdatedBy,
			keyViolUpdatedByID: input.UpdatedByID,

			keyViolApprovedAt:   input.ApprovedAt,
			keyViolApprovedBy:   input.ApprovedBy,
			keyViolApprovedByID: input.ApprovedByID,

			keyViolState: input.State,
		},
	}

	var violation dto.Violation
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("pelanggaran tidak diupdate : validasi id timestamp")
		}

		logger.Error("Gagal mendapatkan violation dari database (EditViolation)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan violation dari database", err)
		return nil, apiErr
	}

	return &violation, nil
}

func (c *violationDao) DeleteViolation(input dto.FilterIDBranch) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyViolID:     input.FilterID,
		keyViolBranch: input.FilterBranch,
		keyViolState:  bson.M{"$lte": 1}, // yang dapat diedit 0 draft dan 1 need approved
	}

	var violation dto.Violation
	err := coll.FindOneAndDelete(ctx, filter).Decode(&violation)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("pelanggaran tidak dihapus : validasi id branch status")
		}

		logger.Error("Gagal menghapus violation dari database (DeleteViolation)", err)
		apiErr := resterr.NewInternalServerError("Gagal menghapus violation dari database", err)
		return nil, apiErr
	}

	return &violation, nil
}

// UploadImage menambahan slice image
func (c *violationDao) UploadImage(violationID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyViolID:     violationID,
		keyViolBranch: strings.ToUpper(filterBranch),
	}
	update := bson.M{
		"$push": bson.M{
			keyViolImages: imagePath,
		},
	}

	var violation dto.Violation
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, violation dengan id %s tidak ditemukan", violationID.Hex()))
		}

		logger.Error("Memasukkan path image violation ke db gagal, (UploadImage)", err)
		apiErr := resterr.NewInternalServerError("Memasukkan path image violation ke db gagal", err)
		return nil, apiErr
	}

	return &violation, nil
}

// DeleteImage mengurangi image dari slice image
func (c *violationDao) DeleteImage(violationID primitive.ObjectID, imagePath string, filterBranch string) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyViolID:     violationID,
		keyViolBranch: strings.ToUpper(filterBranch),
	}
	var violation dto.Violation
	if err := coll.FindOne(ctx, filter).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("delete image gagal, violation dengan id %s tidak ditemukan", violationID.Hex()))
		}

		logger.Error("Delete image violation dari db gagal, (DeleteImage)", err)
		apiErr := resterr.NewInternalServerError("Delete image violation dari db gagal", err)
		return nil, apiErr
	}

	// mendelete dari data yang sudah ditemukan
	var finalImages []string
	for _, image := range violation.Images {
		if image != imagePath {
			finalImages = append(finalImages, image)
		}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	update := bson.M{
		"$set": bson.M{
			keyViolImages: finalImages,
		},
	}

	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("Memasukkan path image gagal, violation dengan id %s tidak ditemukan", violationID.Hex()))
		}

		logger.Error("Memasukkan path image violation ke db gagal, (UploadImage)", err)
		apiErr := resterr.NewInternalServerError("Memasukkan path image violation ke db gagal", err)
		return nil, apiErr
	}

	return &violation, nil
}

func (c *violationDao) GetViolationByID(violationID primitive.ObjectID, branchIfSpecific string) (*dto.Violation, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyViolID: violationID}
	if branchIfSpecific != "" {
		filter[keyViolBranch] = strings.ToUpper(branchIfSpecific)
	}

	var violation dto.Violation
	if err := coll.FindOne(ctx, filter).Decode(&violation); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("Pelanggaran dengan ID %s tidak ditemukan", violationID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan pelanggaran dari database (GetViolationByID)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan pelanggaran dari database", err)
		return nil, apiErr
	}

	return &violation, nil
}

func (c *violationDao) FindViolation(filterA dto.FilterViolation) (dto.ViolationResponseMinList, resterr.APIError) {
	coll := db.DB.Collection(keyViolCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filterA.FilterBranch = strings.ToUpper(filterA.FilterBranch)
	filterA.FilterNoIdentity = strings.ToUpper(filterA.FilterNoIdentity)
	filterA.FilterNoPol = strings.ToUpper(filterA.FilterNoPol)

	// filter
	filter := bson.M{}

	// filter condition
	if filterA.FilterBranch != "" {
		filter[keyViolBranch] = filterA.FilterBranch
	}
	if filterA.FilterNoIdentity != "" {
		filter[keyViolNoIdentity] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterNoIdentity),
		}
	}
	if filterA.FilterNoPol != "" {
		filter[keyViolNoPol] = bson.M{
			"$regex": fmt.Sprintf(".*%s", filterA.FilterNoPol),
		}
	}
	if filterA.FilterState != -1 {
		filter[keyViolState] = filterA.FilterState
	}
	if filterA.FilterStart != 0 && filterA.FilterEnd != 0 {
		filter[keyViolTimeViolation] = bson.M{"$lte": filterA.FilterEnd, "$gte": filterA.FilterStart}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{keyViolID, -1}}) //nolint:govet
	opts.SetLimit(filterA.Limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar pelanggaran dari database (FindViolation)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.ViolationResponseMinList{}, apiErr
	}

	violationList := dto.ViolationResponseMinList{}
	if err = cursor.All(ctx, &violationList); err != nil {
		logger.Error("Gagal decode violationList cursor ke objek slice (FindViolation)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.ViolationResponseMinList{}, apiErr
	}

	return violationList, nil
}
