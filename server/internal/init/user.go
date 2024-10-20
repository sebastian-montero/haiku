package initalizer

import (
	"database/sql"
	http_handlers "wip/internal/handlers/http"
	"wip/internal/repositories"
	"wip/internal/services"
	"wip/internal/utils/logger"
)

func userRepositoryManager(db *sql.DB) *repositories.UserRepository {
	logger.Info("Creating user repository...")
	return &repositories.UserRepository{DB: db}
}

func userServiceManager(repo *repositories.UserRepository) *services.UserService {
	logger.Info("Creating user service...")
	return &services.UserService{Repository: repo}
}

func userHTTPHandlerManager(service *services.UserService) *http_handlers.UserHTTPHandler {
	logger.Info("Creating user HTTP handler...")
	return &http_handlers.UserHTTPHandler{Service: service}
}

func InitUser(db *sql.DB) *http_handlers.UserHTTPHandler {
	repo := userRepositoryManager(db)
	service := userServiceManager(repo)
	handler := userHTTPHandlerManager(service)
	return handler
}
