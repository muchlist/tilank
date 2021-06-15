package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"tilank/dto"
	"tilank/service"
	"tilank/utils/logger"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"tilank/utils/sfunc"
)

func NewJptHandler(jptService *service.JptService) *jptHandler {
	return &jptHandler{
		service: jptService,
	}
}

type jptHandler struct {
	service *service.JptService
}

func (vj *jptHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.JptRequest
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

	insertID, apiErr := vj.service.InsertJpt(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan jpt berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (vj *jptHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := vj.service.DeleteJpt(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("jpt %s berhasil dihapus", id)})
}

func (vj *jptHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	jptID := c.Params("id")

	var req dto.JptEditRequest
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

	jptEdited, apiErr := vj.service.EditJpt(*claims, jptID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": jptEdited})
}

// Get menampilkan jptDetail
func (vj *jptHandler) Get(c *fiber.Ctx) error {
	jptID := c.Params("id")

	jpt, apiErr := vj.service.GetJptByID(jptID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": jpt})
}

// Find menampilkan list jpt
// Query [branch, lambung, nopol, state, limit, start, end]
func (vj *jptHandler) Find(c *fiber.Ctx) error {
	branch := strings.ToUpper(c.Query("branch"))
	name := strings.ToUpper(c.Query("name"))
	tempActive := sfunc.StrToInt(c.Query("active"), 1)

	if branch == "" {
		branch = c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim).Branch
	}

	active := true
	if tempActive == 0 {
		active = false
	}

	filterA := dto.FilterJpt{
		FilterBranch: branch,
		FilterName:   name,
		Active:       active,
	}

	jptList, apiErr := vj.service.FindJpt(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": jptList})
}
