package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"tilank/dto"
	"tilank/service"
	"tilank/utils/logger"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"time"
)

func NewUserHandler(userService service.UserServiceAssumer) *userHandler {
	return &userHandler{
		service: userService,
	}
}

type userHandler struct {
	service service.UserServiceAssumer
}

// Get menampilkan user berdasarkan ID (bukan email)
func (usr *userHandler) Get(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	user, apiErr := usr.service.GetUser(userID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// GetProfile mengembalikan user yang sedang login
func (usr *userHandler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	user, apiErr := usr.service.GetUserByID(claims.Identity)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": user})
}

// Register menambahkan user
func (usr *userHandler) Register(c *fiber.Ctx) error {
	var user dto.UserRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := usr.service.InsertUser(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Register berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

// Find menampilkan list user
func (usr *userHandler) Find(c *fiber.Ctx) error {
	userList, apiErr := usr.service.FindUsers()
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userList})
}

// Edit mengedit user oleh admin
func (usr *userHandler) Edit(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var user dto.UserEditRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := usr.service.EditUser(userID, user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// UpdateFcmToken mengupdateFCM token
func (usr *userHandler) UpdateFcmToken(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var fcmPayload dto.UserUpdateFcmRequest
	if err := c.BodyParser(&fcmPayload); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := fcmPayload.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	userEdited, apiErr := usr.service.EditFcm(claims.Identity, fcmPayload.FcmToken)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": userEdited})
}

// DeleteFcmToken menghapus FCM token
func (usr *userHandler) DeleteFcmToken(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	_, apiErr := usr.service.EditFcm(claims.Identity, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": "fcm di hapus"})
}

// Delete menghapus user, idealnya melalui middleware is_admin
func (usr *userHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	userIDParams := c.Params("user_id")

	if claims.Identity == userIDParams {
		apiErr := resterr.NewBadRequestError("Tidak dapat menghapus akun terkait (diri sendiri)!")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	apiErr := usr.service.DeleteUser(userIDParams)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("user %s berhasil dihapus", userIDParams)})
}

// ChangePassword mengganti password pada user sendiri
func (usr *userHandler) ChangePassword(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var user dto.UserChangePasswordRequest
	if err := c.BodyParser(&user); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := user.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	//mengganti user id dengan user aktif
	user.ID = claims.Identity

	apiErr := usr.service.ChangePassword(user)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(apiErr)
	}

	return c.JSON(fiber.Map{"error": apiErr, "data": "Password berhasil diubah!"})
}

// ResetPassword mengganti password oleh admin pada user tertentu
func (usr *userHandler) ResetPassword(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	data := dto.UserChangePasswordRequest{
		ID:          userID,
		NewPassword: "Password",
	}

	apiErr := usr.service.ResetPassword(data)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("Password user %s berhasil di reset!", c.Params("user_id"))})
}

// Login login
func (usr *userHandler) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := login.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | validate | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := usr.service.Login(login)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// RefreshToken
func (usr *userHandler) RefreshToken(c *fiber.Ctx) error {
	var payload dto.UserRefreshTokenRequest
	if err := c.BodyParser(&payload); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := payload.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("usr: - | parse | %s", err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	response, apiErr := usr.service.Refresh(payload)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": response})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
// sesuai authorisasi aktif. File disimpan di folder static/images dengan nama file == jwt.identity alias username
func (usr *userHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	randomName := fmt.Sprintf("%s%v", claims.Identity, time.Now().Unix())
	pathInDB, apiErr := saveImage(c, *claims, "avatar", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	usersResult, apiErr := usr.service.PutAvatar(claims.Identity, pathInDB)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": usersResult})
}
