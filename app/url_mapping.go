package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"tilank/config"
	"tilank/middleware"
)

func mapUrls(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.LimitRequest())

	app.Static("/image/avatar", "./static/image/avatar")
	app.Static("/image/avatar", "./static/image/violation")

	api := app.Group("/api/v1")

	// PING
	api.Get("/ping", pingHandler.Ping)

	// USER
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/users", middleware.NormalAuth(), userHandler.Find)
	api.Get("/profile", middleware.NormalAuth(), userHandler.GetProfile)
	api.Post("/avatar", middleware.NormalAuth(), userHandler.UploadImage)
	api.Post("/change-password", middleware.FreshAuth(), userHandler.ChangePassword)

	// USER ADMIN
	apiAuthAdmin := app.Group("/api/v1/admin")
	apiAuthAdmin.Use(middleware.NormalAuth(config.RoleAdmin))
	apiAuthAdmin.Post("/users", userHandler.Register)
	apiAuthAdmin.Put("/users/:user_id", userHandler.Edit)
	apiAuthAdmin.Delete("/users/:user_id", userHandler.Delete)
	apiAuthAdmin.Get("/users/:user_id/reset-password", userHandler.ResetPassword)

	// VIOLATION
	api.Post("/violation", middleware.NormalAuth(), violationHandler.Insert)
	api.Get("/violation/:id", middleware.NormalAuth(), violationHandler.Get)
	api.Put("/violation/:id", middleware.NormalAuth(), violationHandler.Edit)
	api.Get("/violation-draft/:id", middleware.NormalAuth(), violationHandler.SendToDraft)
	api.Get("/violation-confirm/:id", middleware.NormalAuth(), violationHandler.SendToConfirmation)
	api.Get("/violation-approve/:id", middleware.NormalAuth(config.RoleHSSE), violationHandler.SendToApproved)
	api.Delete("/violation/:id", middleware.NormalAuth(), violationHandler.Delete)
	// Query [branch, lambung, nopol, state, limit, start, end]
	api.Get("/violation", middleware.NormalAuth(), violationHandler.Find)
	api.Post("/violation-upload-image/:id", middleware.NormalAuth(), violationHandler.UploadImage)
	api.Get("/violation-delete-image/:id/:image", middleware.NormalAuth(), violationHandler.DeleteImage)

	// JPT
	api.Post("/jpt", middleware.NormalAuth(config.RoleHSSE), jptHandler.Insert)
	api.Get("/jpt/:id", middleware.NormalAuth(), jptHandler.Get)
	api.Put("/jpt/:id", middleware.NormalAuth(config.RoleHSSE), jptHandler.Edit)
	api.Delete("/jpt/:id", middleware.NormalAuth(config.RoleHSSE), jptHandler.Delete)
	// Query [branch, name, active ]
	api.Get("/jpt", middleware.NormalAuth(), jptHandler.Find)
}
