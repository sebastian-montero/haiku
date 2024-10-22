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

type SessionHTTPHandler struct {
	Service *services.SessionService
}

func (h *SessionHTTPHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var session models.Session

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateSession(&session); err != nil {
		logger.Error(fmt.Sprintf("Failed to create session: %v", err))
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

func (h *SessionHTTPHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.Service.GetActiveSessions()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get active sessions: %v", err))
		http.Error(w, "Failed to get active sessions", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessions)
}

func (h *SessionHTTPHandler) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	session, err := h.Service.GetSessionByID(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get session: %v", err))
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func (h *SessionHTTPHandler) DeleteSessionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.Service.DeleteSessionByID(id); err != nil {
		logger.Error(fmt.Sprintf("Failed to delete session: %v", err))
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHTTPHandler) EndSessionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["id"]

	var sessionOwner models.SessionOwner

	if err := json.NewDecoder(r.Body).Decode(&sessionOwner); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.EndSessionByID(sessionID, sessionOwner.OwnerID); err != nil {
		logger.Error(fmt.Sprintf("Failed to end session: %v", err))
		http.Error(w, "Failed to end session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
