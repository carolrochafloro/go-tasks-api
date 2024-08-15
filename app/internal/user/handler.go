package user

import (
	"encoding/json"
	"fmt"
	"go-tasks-api/app/internal/logging"
	"go-tasks-api/app/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "Content type is not application/json")
        logging.Warn("Unsupported Media Type")
        return
    }

	if r.Body == nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Request body is missing")
        logging.Warn("Request body is missing")
        return
    }

	defer r.Body.Close()

	var user UserT

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
        logging.Warn("Failed to decode request body: ", err)
        return
	}

	err = user.Validate()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
        logging.Warn("Validation failed: ", err)
        return
	}

	res, err := utils.GetByKey(user.Email, "email", "users")
	if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        logging.Error("Error querying database: ", err)
        return
    }
	if res != nil {
        utils.RespondWithError(w, http.StatusConflict, "User already exists")
        logging.Warn("User already exists: ", user.Email)
        return
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
        logging.Error("Error hashing password: ", err)
        return
	}

	user.Password = string(hashedPassword)

	err = addUserToDB(user)
	if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error adding user to database")
        logging.Error("Error adding user to database: ", err)
        return
    }

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
    logging.Info("User created successfully: ", user.Email)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
    id := vars["id"]

	res, err:= utils.GetByKey(id, "_id", "users")

	if res == nil || err != nil {
		logging.Error("Unable to get user", err)
		w.WriteHeader(http.StatusNotFound) // Status 409 Conflict, pois o usuário já existe
		w.Write([]byte("This user doesn't exist."))
		return
	}

	user, ok := res.(UserT)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read User."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Status 200 OK

	// Codificar o usuário em JSON
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]

	_, found := getUser(id, "_id")

	if !found {
		http.Error(w, "This user doesn't exist.", http.StatusNotFound)
		return
	}

	result, err := deleteUserService(id)

	if err != nil {
		http.Error(w, "Failed to delete user.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Deleted %d user(s)", result.DeletedCount)))
}

func EditProfile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	
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

	_, found := getUser(id, "_id")

	if !found {
		http.Error(w, "This user doesn't exist.", http.StatusNotFound)
		return
	}

	result := updateUser(id, user)

	if (result.ModifiedCount < 1) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Updated %d user(s)", result.ModifiedCount)))

}