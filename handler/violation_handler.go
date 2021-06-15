package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"tilank/dto"
	"tilank/enum"
	"tilank/service"
	"tilank/utils/logger"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"tilank/utils/sfunc"
	"time"
)

func NewViolationHandler(violationService *service.ViolationService) *violationHandler {
	return &violationHandler{
		service: violationService,
	}
}

type violationHandler struct {
	service *service.ViolationService
}

func (vh *violationHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.ViolationRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))

		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	insertID, apiErr := vh.service.InsertViolation(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan violation berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (vh *violationHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := vh.service.DeleteViolation(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("violation %s berhasil dihapus", id)})
}

func (vh *violationHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	violationID := c.Params("id")

	var req dto.ViolationEditRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | parse | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	if err := req.Validate(); err != nil {
		apiErr := resterr.NewBadRequestError(err.Error())
		logger.Info(fmt.Sprintf("u: %s | validate | %s", claims.Name, err.Error()))
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	violationEdited, apiErr := vh.service.EditViolation(*claims, violationID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": violationEdited})
}

func (vh *violationHandler) SendToDraft(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	violationID := c.Params("id")

	violationEdited, apiErr := vh.service.SendToDraftViolation(*claims, violationID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": violationEdited})
}

func (vh *violationHandler) SendToConfirmation(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	violationID := c.Params("id")

	violationEdited, apiErr := vh.service.SendToConfirmationViolation(*claims, violationID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": violationEdited})
}

func (vh *violationHandler) SendToApproved(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	violationID := c.Params("id")

	violationEdited, apiErr := vh.service.ApproveViolation(*claims, violationID)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": violationEdited})
}

// UploadImage melakukan pengambilan file menggunakan form "image" mengecek ekstensi dan memasukkannya ke database
func (vh *violationHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	// cek apakah ID violation && branch ada
	_, apiErr := vh.service.GetViolationByID(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	randomName := fmt.Sprintf("%s%v", id, time.Now().Unix())
	// simpan image
	pathInDb, apiErr := saveImage(c, *claims, "violation", randomName, false)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	violationResult, apiErr := vh.service.PutImage(*claims, id, pathInDb)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": violationResult})
}

// DeleteImage penghapusan gambar di database
func (vh *violationHandler) DeleteImage(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")
	pathImage := c.Params("image")

	// cek apakah ID violation && branch ada
	violation, apiErr := vh.service.GetViolationByID(id, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// untuk sementara gambar tidak dihapus dari system
	imageExist := sfunc.InSlice(pathImage, violation.Images)
	if !imageExist {
		apiErr := resterr.NewBadRequestError("gambar tidak ditemukan")
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	// update path image di database
	violationResult, apiErr := vh.service.DeleteImage(*claims, id, pathImage)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": violationResult})
}

// Get menampilkan violationDetail
func (vh *violationHandler) Get(c *fiber.Ctx) error {
	violationID := c.Params("id")

	violation, apiErr := vh.service.GetViolationByID(violationID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": violation})
}

// Find menampilkan list violation
// Query [branch, lambung, nopol, state, limit, start, end]
func (vh *violationHandler) Find(c *fiber.Ctx) error {
	branch := strings.ToUpper(c.Query("branch"))
	lambung := strings.ToUpper(c.Query("lambung"))
	noPol := c.Query("nopol")
	state := sfunc.StrToInt(c.Query("state"), -1)
	limit := sfunc.StrToInt(c.Query("limit"), 100)
	start := sfunc.StrToInt(c.Query("start"), 0)
	end := sfunc.StrToInt(c.Query("end"), 0)

	if branch == "" {
		branch = c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim).Branch
	}

	filterA := dto.FilterViolation{
		FilterBranch:     branch,
		FilterNoIdentity: lambung,
		FilterNoPol:      noPol,
		FilterState:      enum.IntToState(state),
		FilterStart:      int64(start),
		FilterEnd:        int64(end),
		Limit:            int64(limit),
	}

	violationList, apiErr := vh.service.FindViolation(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": violationList})
}
