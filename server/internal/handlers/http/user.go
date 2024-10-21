package http_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wout/internal/models"
	"wout/internal/services"
	"wout/internal/utils/logger"

	"github.com/gorilla/mux"
)

type UserHTTPHandler struct {
	Service *services.UserService
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUser(&user); err != nil {
		logger.Error(fmt.Sprintf("Failed to create user: %v", err))
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get user: %v", err))
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.Service.DeleteUserByID(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete user: %v", err))
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateUser(&user)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update user: %v", err))
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
