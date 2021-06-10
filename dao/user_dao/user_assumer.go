package userdao

import (
	"tilank/dto"
	"tilank/utils/rest_err"
)

type UserDaoAssumer interface {
	InsertUser(user dto.UserRequest) (*string, resterr.APIError)
	EditUser(userID string, userRequest dto.UserEditRequest) (*dto.UserResponse, resterr.APIError)
	DeleteUser(userID string) resterr.APIError
	PutAvatar(userID string, avatar string) (*dto.UserResponse, resterr.APIError)
	ChangePassword(data dto.UserChangePasswordRequest) resterr.APIError

	GetUserByID(userID string) (*dto.UserResponse, resterr.APIError)
	GetUserByIDWithPassword(userID string) (*dto.User, resterr.APIError)
	FindUser() (dto.UserResponseList, resterr.APIError)
	CheckIDAvailable(email string) (bool, resterr.APIError)
}