package user

import (
	"encoding/json"
	"fmt"
	"go-tasks-api/app/internal/logging"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

	_, userExists := getUser(user.Email, "email")

	if userExists {
		w.WriteHeader(http.StatusConflict) // Status 409 Conflict, pois o usuário já existe
		w.Write([]byte("User already exists"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		logging.Error("There was an error on hashing the password.", err)
	}

	user.Password = string(hashedPassword)

	addUserToDB(user)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created."))
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
    id := vars["id"]

	user, found := getUser(id, "_id")
	
	if !found {
		w.WriteHeader(http.StatusNotFound) // Status 409 Conflict, pois o usuário já existe
		w.Write([]byte("This user doesn't exist."))
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