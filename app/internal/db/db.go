package db

import (
	"context"
	"go-tasks-api/app/internal/logging"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDbService() (*mongo.Client, error) {

	err := godotenv.Load()

	if err != nil {
		logging.Logger.Fatal("Unable to open environment variables:", err)
	}

	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		logging.Warn("Mongo URI is not set up.")
		return nil, err
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		logging.Error("Failed to connect to Mongo", err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logging.Error("Failed to ping MongoDB:", err)
		return nil, err
	}

	defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				logging.Warn("Failed to disconnect from MongoDB:", err)
			}
		}()

	return client, nil

}

