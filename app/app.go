package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"tilank/db"
	"tilank/utils/logger"
)

func RunApp() {
	// inisiasi database
	client, ctx, cancel := db.Init()

	defer func(client *mongo.Client, ctx context.Context) {
		_ = client.Disconnect(ctx)
	}(client, ctx)

	defer cancel()

	app := fiber.New()
	mapUrls(app)
	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		return
	}
}
