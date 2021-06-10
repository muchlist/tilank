package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"net/http"
	"tilank/utils/logger"
	"tilank/utils/rest_err"
	"time"
)

func LimitRequest() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        120,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			logger.Info(fmt.Sprintf("u : %s | limiter | terlalu banyak request", c.IP()))
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": resterr.NewAPIError("terlalu banyak request", http.StatusTooManyRequests, "rate_limiter", []interface{}{"too many requests in a given amount of time"}),
				"data":  nil,
			})
		},
	})
}
