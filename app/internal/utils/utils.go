package utils

import (
	"context"
	"encoding/json"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func GetByKey(s, key, col string, result interface{}) error {

	client := db.Client

	if client == nil {
		logging.Error("MongoDB client is nil")
		return nil
	}

	collection := client.Database("go-tasks").Collection(col)

	var filter interface{}

	if key == "_id" {

		objId, err := primitive.ObjectIDFromHex(s)

		if err != nil {
			logging.Error("Error on converting string to ObjectId", err)
		}

		filter = bson.M{key: objId}

	} else {

		filter = bson.M{key: s}

	}

	err := collection.FindOne(context.TODO(), filter).Decode(result)

	if err != nil {
		if (err == mongo.ErrNoDocuments) {
			return nil
		}
		logging.Error("Error finding document.", err)
		return nil
	}

	return nil
	
}

// resposta de erro em JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
    RespondWithJSON(w, code, map[string]string{"error": message})
}

// resposta em JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}