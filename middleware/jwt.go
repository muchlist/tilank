package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
	"tilank/utils/mjwt"
	"tilank/utils/rest_err"
	"tilank/utils/sfunc"
)

var (
	jwt = mjwt.NewJwt()
)

const (
	headerKey = "Authorization"
	bearerKey = "Bearer"
)

func NormalAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authMustHaveRoleValidator(authHeader, false, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}
		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

func FreshAuth(rolesReq ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(headerKey)
		claims, err := authMustHaveRoleValidator(authHeader, true, rolesReq)
		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{"error": err, "data": nil})
		}

		c.Locals(mjwt.CLAIMS, claims)
		return c.Next()
	}
}

func authMustHaveRoleValidator(authHeader string, mustFresh bool, rolesRequired []string) (*mjwt.CustomClaim, resterr.APIError) {
	if !strings.Contains(authHeader, bearerKey) {
		apiErr := resterr.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}

	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 {
		apiErr := resterr.NewUnauthorizedError("Unauthorized")
		return nil, apiErr
	}

	token, apiErr := jwt.ValidateToken(tokenString[1])
	if apiErr != nil {
		return nil, apiErr
	}

	claims, apiErr := jwt.ReadToken(token)
	if apiErr != nil {
		return nil, apiErr
	}

	if mustFresh {
		if !claims.Fresh {
			apiErr := resterr.NewUnauthorizedError("Memerlukan token yang baru untuk mengakses halaman ini")
			return nil, apiErr
		}
	}

	if len(rolesRequired) != 0 {
		for _, roleReq := range rolesRequired {
			if !sfunc.InSlice(roleReq, claims.Roles) {
				apiErr := resterr.NewUnauthorizedError(fmt.Sprintf("Unauthorized, memerlukan hak akses %s", roleReq))
				return nil, apiErr
			}
		}
	}
	return claims, nil
}