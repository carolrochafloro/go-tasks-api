package user

import (
	"context"
	"encoding/json"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"
	"net/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func init() {
	err := godotenv.Load()

	if err != nil {
		logging.Warn("Failed to load .env file.", err)
	}

	client, err = db.NewDbService()

	if err != nil {
		logging.Warn("Failed to connect to MongoDB", err)
	}
}

func CreateUser(w http.ResponseWriter, r http.Request) {
	user := UserT{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := client.Database("go-tasks").Collection("users")

	_, err = collection.InsertOne(context.TODO(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}