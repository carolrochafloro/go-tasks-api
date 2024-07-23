package user

import (
	"context"
	"encoding/json"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Content type is not application/json", http.StatusUnsupportedMediaType)
        return
    }

	if r.Body == nil {
        http.Error(w, "Request body is missing", http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	user := UserT{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logging.Warn("Wrong content.")
		return
	}

	err = user.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logging.Warn("Validation failed", err)
		return
	}

	client := db.Client

	collection := client.Database("go-tasks").Collection("users")

	result, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logging.Warn("Unable to insert user.")
		println(result)
		return
	}

	w.WriteHeader(http.StatusCreated)
}