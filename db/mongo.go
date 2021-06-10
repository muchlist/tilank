package db

import (
	"context"
	"os"
	"tilank/utils/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout = 10
	mongoURLGetKey = "MONGO_DB_URL_TILANK"
	databaseName   = "tilank"
)

var (
	// DB objek sebagai database objek
	DB       *mongo.Database
	mongoURL = "mongodb://localhost:27017"
)

// Init menginisiasi database
// responsenya digunakan untuk memutus koneksi apabila main program dihentikan
func Init() (*mongo.Client, context.Context, context.CancelFunc) {
	if os.Getenv(mongoURLGetKey) != "" {
		mongoURL = os.Getenv(mongoURLGetKey)
	}

	Client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("gagal membuat client mongodb", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = Client.Connect(ctx)
	if err != nil {
		logger.Error("gagal menghubungkan koneksi mongodb", err)
		panic(err)
	}

	logger.Info("database berhasil terkoneksi")
	DB = Client.Database(databaseName)

	return Client, ctx, cancel
}
