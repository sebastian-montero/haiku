package http_handlers

import (
	"encoding/json"
	"net/http"
	"wip/internal/models"
	"wip/internal/services"
	"wip/internal/utils/logger"
)

type UserHTTPHandler struct {
	UserService *services.UserService
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func UserHTTPHandlerManager(userService *services.UserService) *UserHTTPHandler {
	logger.Info("Creating user HTTP handler...")
	return &UserHTTPHandler{UserService: userService}
}
