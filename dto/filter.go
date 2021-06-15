package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilterIDBranch struct {
	FilterID     primitive.ObjectID
	FilterBranch string
}

type FilterViolation struct {
	FilterBranch     string
	FilterNoIdentity string
	FilterNoPol      string
	FilterState      int
	FilterStart      int64
	FilterEnd        int64
	Limit            int64
}

type FilterJpt struct {
	FilterBranch string
	FilterName   string
	Active       bool
}
