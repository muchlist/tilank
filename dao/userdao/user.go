package userdao

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"tilank/config"
	"tilank/db"
	"tilank/dto"
	"tilank/utils/logger"
	"tilank/utils/rest_err"
	"time"

	"context"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	connectTimeout = 3

	keyUserColl = "user"

	keyUserID        = "_id"
	keyUserEmail     = "email"
	keyUserHashPw    = "hash_pw"
	keyUserName      = "name"
	keyUserRoles     = "roles"
	keyUserBranch    = "branch"
	keyUserAvatar    = "avatar"
	keyUserFcmToken  = "fcm_token"
	keyUserTimeStamp = "timestamp"
)

func NewUserDao() UserDaoAssumer {
	return &userDao{}
}

type userDao struct {
}

// InsertUser menambahkan user dan mengembalikan insertedID, err
func (u *userDao) InsertUser(user dto.UserRequest) (*string, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	user.Name = strings.ToUpper(user.Name)
	user.ID = strings.ToUpper(user.ID)
	user.Email = strings.ToLower(user.Email)
	if user.Roles == nil {
		user.Roles = []string{}
	}

	//nolint:govet
	insertDoc := bson.D{
		{keyUserID, user.ID},
		{keyUserName, user.Name},
		{keyUserEmail, user.Email},
		{keyUserRoles, user.Roles},
		{keyUserBranch, user.Branch},
		{keyUserAvatar, user.Avatar},
		{keyUserHashPw, user.Password},
		{keyUserTimeStamp, user.Timestamp},
	}

	result, err := coll.InsertOne(ctx, insertDoc)
	if err != nil {
		apiErr := resterr.NewInternalServerError("Gagal menyimpan user ke database", err)
		logger.Error("Gagal menyimpan user ke database", err)
		return nil, apiErr
	}

	// insertID := result.InsertedID.(primitive.ObjectID).Hex()
	insertID := result.InsertedID.(string)

	return &insertID, nil
}

// EditUser mengubah user, memerlukan timestamp int64 agar lebih safety pada saat pengeditan oleh dua user
func (u *userDao) EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	userRequest.Name = strings.ToUpper(userRequest.Name)
	userID = strings.ToUpper(userID)
	if userRequest.Roles == nil {
		userRequest.Roles = []string{}
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyUserID:        userID,
		keyUserTimeStamp: userRequest.TimestampFilter,
	}
	update := bson.M{
		"$set": bson.M{
			keyUserName:      userRequest.Name,
			keyUserRoles:     userRequest.Roles,
			keyUserBranch:    userRequest.Branch,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	var user dto.UserResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("User tidak diupdate karena ID atau timestamp tidak valid")
		}

		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

func (u *userDao) EditFcm(userID string, fcmToken string) (*dto.UserResponse, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	userID = strings.ToUpper(userID)

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyUserID: userID,
	}
	update := bson.M{
		"$set": bson.M{
			keyUserFcmToken: fcmToken,
		},
	}

	var user dto.UserResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError("User tidak diupdate karena ID tidak valid")
		}

		logger.Error("Gagal mendapatkan user dari database (UpdateFCM)", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

// DeleteUser menghapus user
func (u *userDao) DeleteUser(userID string) resterr.APIError {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyUserID: strings.ToUpper(userID),
	}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return resterr.NewBadRequestError("User gagal dihapus, dokumen tidak ditemukan")
		}

		logger.Error("Gagal menghapus user dari database", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return apiErr
	}

	if result.DeletedCount == 0 {
		return resterr.NewBadRequestError("User gagal dihapus, dokumen tidak ditemukan")
	}

	return nil
}

// PutAvatar hanya mengubah avatar berdasarkan filter email
func (u *userDao) PutAvatar(userID string, avatar string) (*dto.UserResponse, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(1)

	filter := bson.M{
		keyUserID: strings.ToUpper(userID),
	}
	update := bson.M{
		"$set": bson.M{
			keyUserAvatar:    avatar,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	var user dto.UserResponse
	if err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, resterr.NewBadRequestError(fmt.Sprintf("User avatar gagal diupload, user dengan id %s tidak ditemukan", userID))
		}

		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

// ChangePassword merubah hash_pw dengan password baru sesuai masukan
func (u *userDao) ChangePassword(data dto.UserChangePasswordRequest) resterr.APIError {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	filter := bson.M{
		keyUserID: strings.ToUpper(data.ID),
	}

	update := bson.M{
		"$set": bson.M{
			keyUserHashPw:    data.NewPassword,
			keyUserTimeStamp: time.Now().Unix(),
		},
	}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return resterr.NewBadRequestError("Penggantian password gagal, ID salah")
		}

		logger.Error("Gagal mendapatkan user dari database (ChangePassword)", err)
		apiErr := resterr.NewInternalServerError("Gagal mengganti password user", err)
		return apiErr
	}

	if result.ModifiedCount == 0 {
		return resterr.NewBadRequestError("Penggantian password gagal, kemungkinan ID salah")
	}

	return nil
}

// GetUser mendapatkan user dari database berdasarkan userID, jarang digunakan
// pada case ini biasanya menggunakan email karena user yang digunakan adalah email
func (u *userDao) GetUserByID(userID string) (*dto.UserResponse, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var user dto.UserResponse
	opts := options.FindOne()
	opts.SetProjection(bson.M{keyUserHashPw: 0})

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}, opts).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// apiErr := rest_err.NewNotFoundError(fmt.Sprintf("User dengan FilterID %v tidak ditemukan", userID.Hex()))
			apiErr := resterr.NewNotFoundError(fmt.Sprintf("User dengan ID %s tidak ditemukan", userID))
			return nil, apiErr
		}
		logger.Error("gagal mendapatkan user (GetUserByID) dari database", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return nil, apiErr
	}

	return &user, nil
}

// GetUserByIDWithPassword mendapatkan user dari database berdasarkan id dengan memunculkan passwordhash
// password hash digunakan pada endpoint login dan change password
func (u *userDao) GetUserByIDWithPassword(userID string) (*dto.User, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	var user dto.User

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// karena sudah pasti untuk keperluan login maka error yang dikembalikan unauthorized
			apiErr := resterr.NewUnauthorizedError("Username atau password tidak valid")
			return nil, apiErr
		}

		logger.Error("Gagal mendapatkan user dari database (GetUserByIDWithPassword)", err)
		apiErr := resterr.NewInternalServerError("Error pada database", errors.New("database error"))
		return nil, apiErr
	}

	return &user, nil
}

// FindUser mendapatkan daftar semua user dari database
func (u *userDao) FindUser() (dto.UserResponseList, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	users := dto.UserResponseList{}
	opts := options.Find()
	opts.SetSort(bson.D{{keyUserID, -1}}) //nolint:govet
	sortCursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	if err = sortCursor.All(ctx, &users); err != nil {
		logger.Error("Gagal decode usersCursor ke objek slice", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	return users, nil
}

// FindUser mendapatkan daftar semua user dari database
func (u *userDao) FindUserHSSE(branch string) (dto.UserResponseList, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// filter
	filter := bson.M{
		keyUserBranch: strings.ToUpper(branch),
		keyUserRoles:  config.RoleHSSE,
	}

	users := dto.UserResponseList{}
	opts := options.Find()
	opts.SetSort(bson.D{{keyUserID, -1}}) //nolint:govet
	sortCursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Gagal mendapatkan user dari database", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	if err = sortCursor.All(ctx, &users); err != nil {
		logger.Error("Gagal decode usersCursor ke objek slice", err)
		apiErr := resterr.NewInternalServerError("Database error", err)
		return dto.UserResponseList{}, apiErr
	}

	return users, nil
}

// CheckEmailAvailable melakukan pengecekan apakah alamat email sdh terdaftar di database
// jika ada akan return false ,yang artinya email tidak available
func (u *userDao) CheckIDAvailable(userID string) (bool, resterr.APIError) {
	coll := db.DB.Collection(keyUserColl)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	opts := options.FindOne()
	opts.SetProjection(bson.M{keyUserID: 1})
	var user dto.UserResponse

	if err := coll.FindOne(ctx, bson.M{keyUserID: strings.ToUpper(userID)}, opts).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}

		logger.Error("Gagal mendapatkan user dari database,CheckIDAvailable", err)
		apiErr := resterr.NewInternalServerError("Gagal mendapatkan user dari database", err)
		return false, apiErr
	}

	apiErr := resterr.NewBadRequestError("ID tidak tersedia")
	return false, apiErr
}
