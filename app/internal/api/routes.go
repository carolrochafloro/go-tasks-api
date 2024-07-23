package api

import (
	"go-tasks-api/app/internal/user"
	"net/http"
)

var Routes = map[string]func(http.ResponseWriter, *http.Request){
"/user/new": user.CreateUser,
}