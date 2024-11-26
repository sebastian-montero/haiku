package initalizer

import (
	"database/sql"
	http_handlers "haiku/internal/handlers/http"
	"haiku/internal/repositories"
	"haiku/internal/services"
	"haiku/internal/utils/logger"
)

func notebookRepositoryManager(db *sql.DB) *repositories.NotebookRepository {
	logger.Info("Creating notebook repository...")
	return &repositories.NotebookRepository{DB: db}
}

func notebookServiceManager(repository *repositories.NotebookRepository) *services.NotebookService {
	logger.Info("Creating notebook service...")
	return &services.NotebookService{Repository: repository}
}

func notebookHTTPHandlerManager(service *services.NotebookService) *http_handlers.NotebookHTTPHandler {
	logger.Info("Creating notebook HTTP handler...")
	return &http_handlers.NotebookHTTPHandler{Service: service}
}

func InitNotebook(db *sql.DB) *http_handlers.NotebookHTTPHandler {
	repository := notebookRepositoryManager(db)
	service := notebookServiceManager(repository)
	handler := notebookHTTPHandlerManager(service)
	return handler
}
