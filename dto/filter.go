package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tilank/enum"
)

type FilterIDBranch struct {
	FilterID     primitive.ObjectID
	FilterBranch string
}

// FilterViolation filter
// FilterState -1 menampilkan semuanya
type FilterViolation struct {
	FilterBranch     string
	FilterNoIdentity string
	FilterNoPol      string
	FilterState      enum.State
	FilterStart      int64
	FilterEnd        int64
	Limit            int64
}

type FilterJpt struct {
	FilterBranch string
	FilterName   string
	Active       bool
}

type FilterTruck struct {
	FilterBranch     string
	FilterNoIdentity string
	FilterOwner      string
	Active           bool
	Blocked          bool
}
