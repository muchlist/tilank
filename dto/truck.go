package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Truck struct penuh dari domain truck
type Truck struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedByID string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`

	Branch     string `json:"branch" bson:"branch"`
	NoIdentity string `json:"no_identity" bson:"no_identity"`
	NoPol      string `json:"no_pol" bson:"no_pol"`
	Mark       string `json:"mark" bson:"mark"`
	Owner      string `json:"owner" bson:"owner"`
	Email      string `json:"email" bson:"email"`
	Hp         string `json:"hp" bson:"hp"`
	Deleted    bool   `json:"deleted" bson:"deleted"`

	Score          int   `json:"score" bson:"score"`
	ResetScoreDate int64 `json:"reset_score_date" bson:"reset_score_date"`
	Blocked        bool  `json:"blocked" bson:"blocked"`
	BlockStart     int64 `json:"block_start" bson:"block_start"`
	BlockEnd       int64 `json:"block_end" bson:"block_end"`
}

type TruckScoreEdit struct {
	ID             primitive.ObjectID
	Score          int
	ResetScoreDate int64
	Blocked        bool
	BlockStart     int64
	BlockEnd       int64
}

// TruckRequest user input, id tidak diinput oleh user
type TruckRequest struct {
	NoIdentity string `json:"no_identity" bson:"no_identity"`
	NoPol      string `json:"no_pol" bson:"no_pol"`
	Mark       string `json:"mark" bson:"mark"`
	Owner      string `json:"owner" bson:"owner"`
	Email      string `json:"email" bson:"email"`
	Hp         string `json:"hp" bson:"hp"`
}

type TruckEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	NoIdentity string `json:"no_identity" bson:"no_identity"`
	NoPol      string `json:"no_pol" bson:"no_pol"`
	Mark       string `json:"mark" bson:"mark"`
	Owner      string `json:"owner" bson:"owner"`
	Email      string `json:"email" bson:"email"`
	Hp         string `json:"hp" bson:"hp"`
}

// TruckEditRequest user input
type TruckEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	NoIdentity string `json:"no_identity" bson:"no_identity"`
	NoPol      string `json:"no_pol" bson:"no_pol"`
	Mark       string `json:"mark" bson:"mark"`
	Owner      string `json:"owner" bson:"owner"`
	Email      string `json:"email" bson:"email"`
	Hp         string `json:"hp" bson:"hp"`
}

type TruckResponseMinList []TruckResponseMin

type TruckResponseMin struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch string             `json:"branch" bson:"branch"`

	NoIdentity string `json:"no_identity" bson:"no_identity"`
	NoPol      string `json:"no_pol" bson:"no_pol"`
	Mark       string `json:"mark" bson:"mark"`
	Owner      string `json:"owner" bson:"owner"`
	Email      string `json:"email" bson:"email"`
	Hp         string `json:"hp" bson:"hp"`
	Deleted    bool   `json:"deleted" bson:"deleted"`

	Score          int   `json:"score" bson:"score"`
	ResetScoreDate int64 `json:"reset_score_date" bson:"reset_score_date"`
	Blocked        bool  `json:"blocked" bson:"blocked"`
	BlockStart     int64 `json:"block_start" bson:"block_start"`
	BlockEnd       int64 `json:"block_end" bson:"block_end"`
}
