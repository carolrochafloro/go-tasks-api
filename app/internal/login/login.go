package login

import (
	"encoding/json"
	"go-tasks-api/app/internal/logging"
	"go-tasks-api/app/internal/middleware"
	"go-tasks-api/app/internal/utils"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	
	if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Content type is not application/json", http.StatusUnsupportedMediaType)
        return
    }

	if r.Body == nil {
        http.Error(w, "Request body is missing", http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	login := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	jwt, err := middleware.Authenticate(login)

	if err != nil {
		logging.Error("Unable to authenticate", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to authenticate.")
		return
	}

	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := 
	map[string]string{"message:":"Success",
						"jwt:": jwt,}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)

}