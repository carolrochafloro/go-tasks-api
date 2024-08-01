package db

import (
	"context"
	"go-tasks-api/app/internal/logging"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func NewDbService() (*mongo.Client, error) {

	uri:= os.Getenv("MONGO_URI")

	if uri == "" {
		logging.Warn("Mongo URI is not set up.")
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

	Client = client
 
	return client, nil
}

// criar função para desconectar
