package service

import (
	"net/http"
	userdao "tilank/dao/userdao"
	"tilank/dto"
	"tilank/utils/crypt"
	"tilank/utils/mjwt"
	resterr "tilank/utils/rest_err"
	"time"
)

func NewUserService(dao userdao.UserDaoAssumer, crypto crypt.BcryptAssumer, jwt mjwt.JWTAssumer) UserServiceAssumer {
	return &userService{
		dao:    dao,
		crypto: crypto,
		jwt:    jwt,
	}
}

type userService struct {
	dao    userdao.UserDaoAssumer
	crypto crypt.BcryptAssumer
	jwt    mjwt.JWTAssumer
}

type UserServiceAssumer interface {
	GetUser(userID string) (*dto.UserResponse, resterr.APIError)
	GetUserByID(email string) (*dto.UserResponse, resterr.APIError)
	InsertUser(dto.UserRequest) (*string, resterr.APIError)
	FindUsers() (dto.UserResponseList, resterr.APIError)
	EditUser(userID string, userEdit dto.UserEditRequest) (*dto.UserResponse, resterr.APIError)
	EditFcm(userID string, fcmToken string) (*dto.UserResponse, resterr.APIError)
	DeleteUser(userID string) resterr.APIError
	Login(dto.UserLoginRequest) (*dto.UserLoginResponse, resterr.APIError)
	Refresh(login dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, resterr.APIError)
	PutAvatar(userID string, fileLocation string) (*dto.UserResponse, resterr.APIError)
	ChangePassword(data dto.UserChangePasswordRequest) resterr.APIError
	ResetPassword(data dto.UserChangePasswordRequest) resterr.APIError
}

// GetUser mendapatkan user dari database
func (u *userService) GetUser(userID string) (*dto.UserResponse, resterr.APIError) {
	user, err := u.dao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail mendapatkan user berdasarkan email
func (u *userService) GetUserByID(userID string) (*dto.UserResponse, resterr.APIError) {
	user, err := u.dao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers
func (u *userService) FindUsers() (dto.UserResponseList, resterr.APIError) {
	userList, err := u.dao.FindUser()
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// InsertUser melakukan register user
func (u *userService) InsertUser(user dto.UserRequest) (*string, resterr.APIError) {
	// cek ketersediaan id
	_, err := u.dao.CheckIDAvailable(user.ID)
	if err != nil {
		return nil, err
	}
	// END cek ketersediaan id

	hashPassword, err := u.crypto.GenerateHash(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	user.Timestamp = time.Now().Unix()

	insertedID, err := u.dao.InsertUser(user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

// EditUser
func (u *userService) EditUser(userID string, request dto.UserEditRequest) (*dto.UserResponse, resterr.APIError) {
	result, err := u.dao.EditUser(userID, request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *userService) EditFcm(userID string, fcmToken string) (*dto.UserResponse, resterr.APIError) {
	result, err := u.dao.EditFcm(userID, fcmToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUser
func (u *userService) DeleteUser(userID string) resterr.APIError {
	err := u.dao.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

// Login
func (u *userService) Login(login dto.UserLoginRequest) (*dto.UserLoginResponse, resterr.APIError) {
	user, err := u.dao.GetUserByIDWithPassword(login.ID)
	if err != nil {
		return nil, err
	}

	if !u.crypto.IsPWAndHashPWMatch(login.Password, user.HashPw) {
		return nil, resterr.NewUnauthorizedError("Username atau password tidak valid")
	}

	if login.Limit == 0 || login.Limit > 60*24*30 { // 30 days
		login.Limit = 60 * 24 * 30
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
		ExtraMinute: time.Duration(login.Limit),
		Type:        mjwt.Access,
		Fresh:       true,
	}

	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
		ExtraMinute: 60 * 24 * 90, // 90 days
		Type:        mjwt.Refresh,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserLoginResponse{
		ID:           user.ID,
		Name:         user.Name,
		Branch:       user.Branch,
		Email:        user.Email,
		Roles:        user.Roles,
		Avatar:       user.Avatar,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expired:      time.Now().Add(time.Minute * time.Duration(login.Limit)).Unix(),
	}

	return &userResponse, nil
}

// Refresh token
func (u *userService) Refresh(payload dto.UserRefreshTokenRequest) (*dto.UserRefreshTokenResponse, resterr.APIError) {
	token, apiErr := u.jwt.ValidateToken(payload.RefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}
	claims, apiErr := u.jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	// cek apakah tipe claims token yang dikirim adalah tipe refresh (1)
	if claims.Type != mjwt.Refresh {
		return nil, resterr.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{"not a refresh token"})
	}

	// mendapatkan data terbaru dari user
	user, apiErr := u.dao.GetUserByID(claims.Identity)
	if apiErr != nil {
		return nil, apiErr
	}

	if payload.Limit == 0 || payload.Limit > 60*24*30 { // 30 day
		payload.Limit = 60 * 24 * 30
	}

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID,
		Name:        user.Name,
		Roles:       user.Roles,
		Branch:      user.Branch,
		ExtraMinute: time.Duration(payload.Limit),
		Type:        mjwt.Access,
		Fresh:       false,
	}

	accessToken, err := u.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return nil, err
	}

	userRefreshTokenResponse := dto.UserRefreshTokenResponse{
		AccessToken: accessToken,
		Expired:     time.Now().Add(time.Minute * time.Duration(payload.Limit)).Unix(),
	}

	return &userRefreshTokenResponse, nil
}

// PutAvatar memasukkan lokasi file (path) ke dalam database user
func (u *userService) PutAvatar(userID string, fileLocation string) (*dto.UserResponse, resterr.APIError) {
	user, err := u.dao.PutAvatar(userID, fileLocation)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword melakukan perbandingan hashpassword lama dan memasukkan hashpassword baru ke database
func (u *userService) ChangePassword(data dto.UserChangePasswordRequest) resterr.APIError {
	if data.Password == data.NewPassword {
		return resterr.NewBadRequestError("Gagal mengganti password, password tidak boleh sama dengan sebelumnya!")
	}

	userResult, err := u.dao.GetUserByIDWithPassword(data.ID)
	if err != nil {
		return err
	}

	if !u.crypto.IsPWAndHashPWMatch(data.Password, userResult.HashPw) {
		return resterr.NewBadRequestError("Gagal mengganti password, password salah!")
	}

	newPasswordHash, err := u.crypto.GenerateHash(data.NewPassword)
	if err != nil {
		return err
	}
	data.NewPassword = newPasswordHash

	_ = u.dao.ChangePassword(data)

	return nil
}

// ResetPassword . inputan password berada di level handler
// hanya memproses field newPassword, mengabaikan field password
func (u *userService) ResetPassword(data dto.UserChangePasswordRequest) resterr.APIError {
	newPasswordHash, err := u.crypto.GenerateHash(data.NewPassword)
	if err != nil {
		return err
	}
	data.NewPassword = newPasswordHash

	err = u.dao.ChangePassword(data)
	if err != nil {
		return err
	}

	return nil
}
