package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Rules struct penuh dari domain pelanggaran
type Rules struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`
	Branch      string             `json:"branch" bson:"branch"`
	Score       int                `json:"score" bson:"score"`
	BlockTime   int64              `json:"block_time" bson:"block_time"`
	Description string             `json:"description" bson:"description"`
}

type RulesEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64
	UpdatedAt       int64  `json:"updated_at" bson:"updated_at"`
	UpdatedBy       string `json:"updated_by" bson:"updated_by"`
	UpdatedByID     string `json:"updated_by_id" bson:"updated_by_id"`
	Score           int    `json:"score" bson:"score"`
	BlockTime       int64  `json:"block_time" bson:"block_time"`
	Description     string `json:"description" bson:"description"`
}

type RulesRequest struct {
	Score       int    `json:"score" bson:"score"`
	BlockTime   int64  `json:"block_time" bson:"block_time"`
	Description string `json:"description" bson:"description"`
}

type RulesEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	Score       int    `json:"score" bson:"score"`
	BlockTime   int64  `json:"block_time" bson:"block_time"`
	Description string `json:"description" bson:"description"`
}
