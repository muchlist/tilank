package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// Jpt struct penuh dari domain pelanggaran
type Jpt struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   int64              `json:"created_at" bson:"created_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedByID string             `json:"created_by_id" bson:"created_by_id"`
	UpdatedAt   int64              `json:"updated_at" bson:"updated_at"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	UpdatedByID string             `json:"updated_by_id" bson:"updated_by_id"`

	Branch    string `json:"branch" bson:"branch"`
	Name      string `json:"name" bson:"name"`
	OwnerName string `json:"owner_name" bson:"owner_name"`
	IDPelindo string `json:"id_pelindo" bson:"id_pelindo"`
	Hp        string `json:"hp" bson:"hp"`
	Email     string `json:"email" bson:"email"`
	Deleted   bool   `json:"deleted" bson:"deleted"`
}

// JptRequest user input, id tidak diinput oleh user
type JptRequest struct {
	ID string `json:"-" bson:"-"`

	Name      string `json:"name" bson:"name"`
	OwnerName string `json:"owner_name" bson:"owner_name"`
	IDPelindo string `json:"id_pelindo" bson:"id_pelindo"`
	Hp        string `json:"hp" bson:"hp"`
	Email     string `json:"email" bson:"email"`
}

type JptEdit struct {
	ID              primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64

	UpdatedAt   int64
	UpdatedBy   string
	UpdatedByID string

	Name      string
	OwnerName string
	IDPelindo string
	Hp        string
	Email     string
}

// JptEditRequest user input
type JptEditRequest struct {
	FilterTimestamp int64 `json:"filter_timestamp"`

	Name      string `json:"name" bson:"name"`
	OwnerName string `json:"owner_name" bson:"owner_name"`
	IDPelindo string `json:"id_pelindo" bson:"id_pelindo"`
	Hp        string `json:"hp" bson:"hp"`
	Email     string `json:"email" bson:"email"`
}

type JptResponseMinList []JptResponseMin

type JptResponseMin struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Branch string             `json:"branch" bson:"branch"`

	Name      string `json:"name" bson:"name"`
	OwnerName string `json:"owner_name" bson:"owner_name"`
	IDPelindo string `json:"id_pelindo" bson:"id_pelindo"`
	Hp        string `json:"hp" bson:"hp"`
	Email     string `json:"email" bson:"email"`
	Deleted   bool   `json:"deleted" bson:"deleted"`
}
