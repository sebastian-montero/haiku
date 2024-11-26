package http_handlers

import (
	"encoding/json"
	"fmt"
	"haiku/internal/models"
	"haiku/internal/services"
	"haiku/internal/utils/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type ContentHTTPHandler struct {
	Service *services.ContentService
}

func (h *ContentHTTPHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	var content models.Content

	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateContent(&content); err != nil {
		logger.Error(fmt.Sprintf("Failed to create content: %v", err))
		http.Error(w, "Failed to create content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(content)
}

func (h *ContentHTTPHandler) GetLatestContentBySessionId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId := vars["session_id"]

	content, err := h.Service.GetLatestContentBySessionId(sessionId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get content by session id: %v", err))
		http.Error(w, "Failed to get content by session id", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)
}
