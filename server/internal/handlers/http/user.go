package http_handlers

import (
	"encoding/json"
	"fmt"
	"haiku/internal/models"
	"haiku/internal/services"
	"haiku/internal/utils"
	"haiku/internal/utils/logger"
	"net/http"

	"encoding/base64"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

func (h *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.Service.GetUserByUsername(credentials.Username)
    if err != nil {
        logger.Error(fmt.Sprintf("Failed to get user: %v", err))
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Decode the stored salt
    salt, err := base64.StdEncoding.DecodeString(user.Salt)
    if err != nil {
        logger.Error(fmt.Sprintf("Failed to decode salt: %v", err))
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Hash the provided password with the stored salt
    saltedPassword := append(salt, []byte(credentials.Password)...)
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), saltedPassword); err != nil {
        logger.Error(fmt.Sprintf("Invalid password: %v", err))
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    token, err := utils.GenerateJWT(user)
    if err != nil {
        logger.Error(fmt.Sprintf("Failed to generate token: %v", err))
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHTTPHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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

    token, err := utils.GenerateJWT(user)
    if err != nil {
        logger.Error(fmt.Sprintf("Failed to generate token: %v", err))
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}