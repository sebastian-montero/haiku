package http_handlers

import (
	"encoding/json"
	"net/http"
	"wip/internal/models"
	"wip/internal/services"
	"wip/internal/utils/logger"

	"github.com/gorilla/mux"
)

type UserHTTPHandler struct {
	UserService *services.UserService
}


func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return 
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}


func (h *UserHTTPHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User = h.UserService.GetUserByID(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.UserService.DeleteUserByID(id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	h.UserService.UpdateUser(&user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UserHTTPHandlerManager(userService *services.UserService) *UserHTTPHandler {
	logger.Info("Creating user HTTP handler...")
	return &UserHTTPHandler{UserService: userService}
}
