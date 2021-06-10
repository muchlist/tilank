package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type FilterBranchNameCatDisable struct {
	FilterBranch   string
	FilterName     string
	FilterCategory string
	FilterDisable  bool
}

// FilterIDBranchCreateGte
// field FilterCreateGTE digunakan untuk memberikan batas waktu, misalnya saat akan menghapus dokumen
// maka FilterCreateGTE di isi dengan tanggal sekarang kurang sekian waktu (misalnya 1 hari)
type FilterIDBranchCreateGte struct {
	FilterID        primitive.ObjectID
	FilterBranch    string
	FilterCreateGTE int64
}

// FilterIDBranchCategoryCreateGte
// field FilterCreateGTE digunakan untuk memberikan batas waktu, misalnya saat akan menghapus dokumen
// maka FilterCreateGTE di isi dengan tanggal sekarang kurang sekian waktu (misalnya 1 hari)
type FilterIDBranchCategoryCreateGte struct {
	FilterID          primitive.ObjectID
	FilterBranch      string
	FilterSubCategory string
	FilterCreateGTE   int64
}

type FilterIDBranchTimestamp struct {
	FilterID        primitive.ObjectID
	FilterBranch    string
	FilterTimestamp int64
}

type FilterIDBranch struct {
	FilterID     primitive.ObjectID
	FilterBranch string
}

type FilterIDBranchAuthor struct {
	FilterID       primitive.ObjectID
	FilterBranch   string
	FilterAuthorID string
}

type FilterParentIDChildIDAuthor struct {
	FilterParentID primitive.ObjectID
	FilterChildID  string
	FilterAuthorID string
}

type FilterTimeRangeLimit struct {
	FilterStart int64
	FilterEnd   int64
	Limit       int64
}

type FilterBranchCompleteTimeRangeLimit struct {
	FilterBranch         string
	FilterCompleteStatus int
	FilterStart          int64
	FilterEnd            int64
	Limit                int64
}

type FilterBranchCategory struct {
	FilterBranch   string
	FilterCategory string
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
