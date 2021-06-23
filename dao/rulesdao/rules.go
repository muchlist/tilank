package rulesdao

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
	keyRulesCollection = "rules"

	keyRulesID          = "_id"
	keyRulesUpdatedAt   = "updated_at"
	keyRulesUpdatedBy   = "updated_by"
	keyRulesUpdatedByID = "updated_by_id"
	keyRulesBranch      = "branch"
	keyRulesScore       = "score"
	keyRulesBlockTime   = "block_time"
	keyRulesDescription = "description"
)

func NewRulesDao() RulesDaoAssumer {
	return &rulesDao{}
}

type rulesDao struct {
}

type RulesDaoAssumer interface {
	InsertRules(input dto.Rules) (*string, resterr.APIError)
	EditRules(input dto.RulesEdit) (*dto.Rules, resterr.APIError)
	DeleteRules(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Rules, resterr.APIError)

	GetRulesByID(rulesID primitive.ObjectID, branchIfSpecific string) (*dto.Rules, resterr.APIError)
	GetRulesByScore(score int, branch string) (*dto.Rules, resterr.APIError)
	FindRules() ([]dto.Rules, resterr.APIError)
}

func (c *rulesDao) InsertRules(input dto.Rules) (*string, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	input.Branch = strings.ToUpper(input.Branch)

	result, err := coll.InsertOne(ctx, input)
	if err != nil {
		apiErr := resterr.NewInternalServerError("Gagal menyimpan rules ke database", err)
		logger.Error("Gagal menyimpan rules ke database, (InsertRules)", err)
		return nil, apiErr
	}

	insertID := result.InsertedID.(primitive.ObjectID).Hex()

	return &insertID, nil
}

func (c *rulesDao) EditRules(input dto.RulesEdit) (*dto.Rules, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyRulesID:        input.ID,
		keyRulesBranch:    input.FilterBranch,
		keyRulesUpdatedAt: input.FilterTimestamp,
	}

	update := bson.M{
		"$set": bson.M{
			keyRulesUpdatedAt:   input.UpdatedAt,
			keyRulesUpdatedBy:   input.UpdatedBy,
			keyRulesUpdatedByID: input.UpdatedByID,
			keyRulesScore:       input.Score,
			keyRulesBlockTime:   input.BlockTime,
			keyRulesDescription: input.Description,
		},
	}

	var rules dto.Rules
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&rules); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("rules tidak diupdate : validasi id timestamp")
		}

		logger.Error("Gagal mendapatkan rules dari database (EditRules)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan rules dari database", err)
		return nil, apiErr
	}

	return &rules, nil
}

func (c *rulesDao) DeleteRules(input dto.FilterIDBranch, isSoftDelete bool) (*dto.Rules, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndDelete()

	filter := bson.M{
		keyRulesID:     input.FilterID,
		keyRulesBranch: input.FilterBranch,
	}

	var rules dto.Rules
	err := coll.FindOneAndDelete(ctx, filter, opts).Decode(&rules)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("Rules tidak dihapus : validasi id branch")
		}

		logger.Error("Gagal menghapus rules dari database (DeleteRules)", err)
		apiErr := resterr.NewInternalServerError("Gagal menghapus rules dari database", err)
		return nil, apiErr
	}

	return &rules, nil
}

func (c *rulesDao) GetRulesByID(rulesID primitive.ObjectID, branchIfSpecific string) (*dto.Rules, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyRulesID: rulesID}
	if branchIfSpecific != "" {
		filter[keyRulesBranch] = strings.ToUpper(branchIfSpecific)
	}

	var rules dto.Rules
	if err := coll.FindOne(ctx, filter).Decode(&rules); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("Rules dengan ID %s tidak ditemukan", rulesID.Hex()))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan rules dari database (GetRulesByID)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan rules dari database", err)
		return nil, apiErr
	}

	return &rules, nil
}

func (c *rulesDao) GetRulesByScore(score int, branch string) (*dto.Rules, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{keyRulesScore: score, keyRulesBranch: strings.ToUpper(branch)}

	var rules dto.Rules
	if err := coll.FindOne(ctx, filter).Decode(&rules); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("Rules dengan ID %d tidak ditemukan", score))
			return nil, apiErr
		}

		logger.Error("gagal mendapatkan rules dari database (GetRulesByID)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan rules dari database", err)
		return nil, apiErr
	}

	return &rules, nil
}

func (c *rulesDao) FindRules() ([]dto.Rules, resterr.APIError) {
	coll := db.DB.Collection(keyRulesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{}

	opts := options.Find()
	opts.SetSort(bson.D{{keyRulesScore, 1}}) //nolint:govet

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan daftar rules dari database (FindRules)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return []dto.Rules{}, apiErr
	}

	var rulesList []dto.Rules
	if err = cursor.All(ctx, &rulesList); err != nil {
		logger.Error("Gagal decode rulesList cursor ke objek slice (FindRules)", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return []dto.Rules{}, apiErr
	}

	return rulesList, nil
}
