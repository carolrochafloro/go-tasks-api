package user

import (
	"encoding/json"
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

	var existingUser UserT
	err = utils.GetByKey(user.Email, "email", "users", &existingUser)
	if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        logging.Error("Error querying database: ", err)
        return
    }
	
	if existingUser.Email != "" {
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

	var user UserT

	err:= utils.GetByKey(id, "_id", "users", &user)

	if err != nil {
		logging.Error("Unable to get user", err)
		utils.RespondWithError(w, http.StatusNotFound, "User not found.")
		return
	}


	// Codificar o usuário em JSON
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logging.Error("Unable to encode user", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to encode user.")
		return
	}

	utils.RespondWithJSON(w, http.StatusFound, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]

	var user UserT

	err := utils.GetByKey(id, "_id", "users", &user)

	if err != nil {
		logging.Error("This user doesn't exist.")
		utils.RespondWithError(w, http.StatusNotFound, "User not found.")
		return
	}

	_, err = deleteUserService(id)

	if err != nil {
		logging.Error("Failed to delete user.", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to delete user.")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "User was deleted.")
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

	var existingUser UserT

	err = utils.GetByKey(id, "_id", "users", &existingUser)

	if err != nil {
		logging.Error("User not found.")
		utils.RespondWithError(w, http.StatusNotFound, "User not found.")
		return
	}

	result := updateUser(id, user)

	if (result.ModifiedCount < 1) {
		logging.Error("No changes were performed.")
		utils.RespondWithError(w, http.StatusBadGateway, "Unable to update user.")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "The user was updated.")

}