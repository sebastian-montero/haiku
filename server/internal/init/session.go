package initalizer

import (
	"database/sql"
	http_handlers "wout/internal/handlers/http"
	"wout/internal/repositories"
	"wout/internal/services"
	"wout/internal/utils/logger"
)

func sessionRepositoryManager(db *sql.DB) *repositories.SessionRepository {
	logger.Info("Creating session repository...")
	return &repositories.SessionRepository{DB: db}
}

func sessionServiceManager(repository *repositories.SessionRepository) *services.SessionService {
	logger.Info("Creating session service...")
	return &services.SessionService{Repository: repository}
}

func sessionHTTPHandlerManager(service *services.SessionService) *http_handlers.SessionHTTPHandler {
	logger.Info("Creating session HTTP handler...")
	return &http_handlers.SessionHTTPHandler{Service: service}
}

func InitSession(db *sql.DB) *http_handlers.SessionHTTPHandler {
	repository := sessionRepositoryManager(db)
	service := sessionServiceManager(repository)
	handler := sessionHTTPHandlerManager(service)
	return handler
}
