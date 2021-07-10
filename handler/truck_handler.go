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

func NewTruckHandler(truckService *service.TruckService) *truckHandler {
	return &truckHandler{
		service: truckService,
	}
}

type truckHandler struct {
	service *service.TruckService
}

func (th *truckHandler) Insert(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)

	var req dto.TruckRequest
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

	insertID, apiErr := th.service.InsertTruck(*claims, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	res := fmt.Sprintf("Menambahkan truck berhasil, ID: %s", *insertID)
	return c.JSON(fiber.Map{"error": nil, "data": res})
}

func (th *truckHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	id := c.Params("id")

	apiErr := th.service.DeleteTruck(*claims, id)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": fmt.Sprintf("truck %s berhasil dihapus", id)})
}

func (th *truckHandler) Edit(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	truckID := c.Params("id")

	var req dto.TruckEditRequest
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

	truckEdited, apiErr := th.service.EditTruck(*claims, truckID, req)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}
	return c.JSON(fiber.Map{"error": nil, "data": truckEdited})
}

// Get menampilkan truckDetail
func (th *truckHandler) Get(c *fiber.Ctx) error {
	truckID := c.Params("id")

	truck, apiErr := th.service.GetTruckByID(truckID, "")
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": truck})
}

// GetByNopol menampilkan truckDetail berdasarkan nopol
func (th *truckHandler) GetByNoLambung(c *fiber.Ctx) error {
	claims := c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim)
	noLambung := c.Params("id")

	truck, apiErr := th.service.GetTruckByNoLambung(noLambung, claims.Branch)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": truck})
}

// Find menampilkan list truck
// Query [branch, identity, owner, active, block ]
func (th *truckHandler) Find(c *fiber.Ctx) error {
	branch := strings.ToUpper(c.Query("branch"))
	noIdentity := strings.ToUpper(c.Query("identity"))
	owner := strings.ToUpper(c.Query("owner"))
	tempActive := sfunc.StrToInt(c.Query("active"), 1)
	tempBlocked := sfunc.StrToInt(c.Query("block"), 0)

	if branch == "" {
		branch = c.Locals(mjwt.CLAIMS).(*mjwt.CustomClaim).Branch
	}

	active := true
	if tempActive == 0 {
		active = false
	}

	blocked := false
	if tempBlocked == 1 {
		blocked = true
	}

	filterA := dto.FilterTruck{
		FilterBranch:     branch,
		FilterNoIdentity: noIdentity,
		FilterOwner:      owner,
		Active:           active,
		Blocked:          blocked,
	}

	truckList, apiErr := th.service.FindTruck(filterA)
	if apiErr != nil {
		return c.Status(apiErr.Status()).JSON(fiber.Map{"error": apiErr, "data": nil})
	}

	return c.JSON(fiber.Map{"error": nil, "data": truckList})
}
