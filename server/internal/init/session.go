package initalizer

import (
	"database/sql"
	http_handlers "wout/internal/handlers/http"
	ws_handler "wout/internal/handlers/ws"
	"wout/internal/repositories"
	"wout/internal/services"
	"wout/internal/utils/logger"

	"github.com/gorilla/websocket"
)

// Creates and returns a new session repository
func sessionRepositoryManager(db *sql.DB) *repositories.SessionRepository {
	logger.Info("Creating session repository...")
	return &repositories.SessionRepository{DB: db}
}

// Creates and returns a new session service
func sessionServiceManager(repository *repositories.SessionRepository) *services.SessionService {
	logger.Info("Creating session service...")
	return &services.SessionService{
		SessionRepository:  repository,
		NotebookRepository: notebookRepositoryManager(repository.DB),
		ContentRepository:  contentRepositoryManager(repository.DB),
	}
}

// Creates and returns a new HTTP session handler
func sessionHTTPHandlerManager(service *services.SessionService) *http_handlers.SessionHTTPHandler {
	logger.Info("Creating session HTTP handler...")
	return &http_handlers.SessionHTTPHandler{Service: service}
}

// Creates and returns a new WebSocket session handler
func sessionWSHandlerManager(service *services.SessionService) *ws_handler.WebSocketHandler {
	logger.Info("Creating session WebSocket handler...")
	return &ws_handler.WebSocketHandler{
		Service:        service,
		Clients:        make(map[int]map[*websocket.Conn]bool),
		SessionContent: make(map[int]string),
	}
}

// Initializes both HTTP and WebSocket session handlers
func InitSession(db *sql.DB) (*http_handlers.SessionHTTPHandler, *ws_handler.WebSocketHandler) {
	repository := sessionRepositoryManager(db)
	service := sessionServiceManager(repository)
	handler := sessionHTTPHandlerManager(service)
	wsHandler := sessionWSHandlerManager(service)
	return handler, wsHandler
}
