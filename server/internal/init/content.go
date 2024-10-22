package initalizer

import (
	"database/sql"
	http_handlers "wout/internal/handlers/http"
	"wout/internal/repositories"
	"wout/internal/services"
	"wout/internal/utils/logger"
)

func contentRepositoryManager(db *sql.DB) *repositories.ContentRepository {
	logger.Info("Creating content repository...")
	return &repositories.ContentRepository{DB: db}
}

func contentServiceManager(repository *repositories.ContentRepository) *services.ContentService {
	logger.Info("Creating content service...")
	sesessionRepository := repositories.SessionRepository{DB: repository.DB}
	notebookRepository := repositories.NotebookRepository{DB: repository.DB}
	return &services.ContentService{ContentRepository: repository, SessionRepository: &sesessionRepository, NotebookRepository: &notebookRepository}
}

func contentHTTPHandlerManager(service *services.ContentService) *http_handlers.ContentHTTPHandler {
	logger.Info("Creating content HTTP handler...")
	return &http_handlers.ContentHTTPHandler{Service: service}
}

func InitContent(db *sql.DB) *http_handlers.ContentHTTPHandler {
	repository := contentRepositoryManager(db)
	service := contentServiceManager(repository)
	handler := contentHTTPHandlerManager(service)
	return handler
}
