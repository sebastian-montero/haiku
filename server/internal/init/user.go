package initalizer

import (
	"database/sql"
	http_handlers "wip/internal/handlers/http"
	"wip/internal/repositories"
	"wip/internal/services"
)

func InitUser(db *sql.DB) *http_handlers.UserHTTPHandler {
	userRepository := repositories.UserRepositoryManager(db)
	userService := services.UserServiceManager(userRepository)
	userHandler := http_handlers.UserHTTPHandlerManager(userService)
	return userHandler
}
