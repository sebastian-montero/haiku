package http_handlers

import (
	"encoding/json"
	"net/http"
	"wip/internal/models"
	"wip/internal/services"

	"github.com/gorilla/mux"
)

type UserHTTPHandler struct {
	Service *services.UserService
}

func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUser(&user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User = h.Service.GetUserByID(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.Service.DeleteUserByID(id)

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	h.Service.UpdateUser(&user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
