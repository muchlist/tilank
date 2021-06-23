package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/enum"
)

// Violation struct penuh dari domain pelanggaran
type Violation struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt    int64              `json:"created_at" bson:"created_at"`
	CreatedBy    string             `json:"created_by" bson:"created_by"`
	CreatedByID  string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt    int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy    string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID  string             `json:"updated_by_id" bson:"updated_by_id"`
	ApprovedAt   int64              `json:"approved_at" bson:"approved_at"`
	ApprovedBy   string             `json:"approved_by" bson:"approved_by"`
	ApprovedByID string             `json:"approved_by_id" bson:"approved_by_id"`
	Branch       string             `json:"branch" bson:"branch"`
	// State 0 Draft, 1 Need Approve, 2 Approved, 3 sendToJPT
	State           enum.State `json:"state" bson:"state"`
	NoIdentity      string     `json:"no_identity" bson:"no_identity"`
	NoPol           string     `json:"no_pol" bson:"no_pol"`
	Mark            string     `json:"mark" bson:"mark"`
	Owner           string     `json:"owner" bson:"owner"`
	TypeViolation   string     `json:"type_violation" bson:"type_violation"`
	DetailViolation string     `json:"detail_violation" bson:"detail_violation"`
	TimeViolation   int64      `json:"time_violation" bson:"time_violation"`
	Location        string     `json:"location" bson:"location"`
	Images          []string   `json:"images" bson:"images"`
}

// ViolationRequest user input, id tidak diinput oleh user
type ViolationRequest struct {
	ID string `json:"-" bson:"-"`
	// State 0 Draft, 1 Need Approve, 2 Approved, 3 sendToJPT
	State           enum.State `json:"state" bson:"state"`
	NoIdentity      string     `json:"no_identity" bson:"no_identity"`
	TypeViolation   string     `json:"type_violation" bson:"type_violation"`
	DetailViolation string     `json:"detail_violation" bson:"detail_violation"`
	TimeViolation   int64      `json:"time_violation" bson:"time_violation"`
	Location        string     `json:"location" bson:"location"`
}

type ViolationEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	ApprovedAt   int64
	ApprovedBy   string
	ApprovedByID string

	NoIdentity      string
	NoPol           string
	Mark            string
	Owner           string
	TypeViolation   string
	DetailViolation string
	TimeViolation   int64
	Location        string
}

type ViolationConfirm struct {
	ID           primitive.ObjectID
	FilterBranch string

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	ApprovedAt   int64
	ApprovedBy   string
	ApprovedByID string

	State enum.State
}

// ViolationEditRequest user input
type ViolationEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	// State 0 Draft, 1 Need Approve, 2 Approved, 3 sendToJPT
	State           int    `json:"state" bson:"state"`
	NoIdentity      string `json:"no_identity" bson:"no_identity"`
	TypeViolation   string `json:"type_violation" bson:"type_violation"`
	DetailViolation string `json:"detail_violation" bson:"detail_violation"`
	TimeViolation   int64  `json:"time_violation" bson:"time_violation"`
	Location        string `json:"location" bson:"location"`
}

type ViolationResponseMinList []ViolationResponseMin

type ViolationResponseMin struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch string             `json:"branch" bson:"branch"`
	// State 0 Draft, 1 Need Approve, 2 Approved, 3 sendToJPT
	State           enum.State `json:"state" bson:"state"`
	NoIdentity      string     `json:"no_identity" bson:"no_identity"`
	NoPol           string     `json:"no_pol" bson:"no_pol"`
	Owner           string     `json:"owner" bson:"owner"`
	TypeViolation   string     `json:"type_violation" bson:"type_violation"`
	DetailViolation string     `json:"detail_violation" bson:"detail_violation"`
	TimeViolation   int64      `json:"time_violation" bson:"time_violation"`
	Location        string     `json:"location" bson:"location"`
	Images          []string   `json:"images" bson:"images"`
}
