package app

import (
	"context"
	"go-tasks-api/app/internal/api"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"

	"github.com/joho/godotenv"
)

var (
	MainContext context.Context
	MainCancel  context.CancelFunc
)


func New() {

	MainContext, MainCancel = context.WithCancel(context.Background())

	logging.NewLogger()

	logging.Info("Starting server...")

	err := godotenv.Load()
	if err != nil {
		logging.Warn("Unable to load environment variables.")
	}

	db.NewDbService()
	api.NewHTTPService()
	api.HTTPService.StartServer()

}

// inserir termination listener
