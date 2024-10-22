package http_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wout/internal/models"
	"wout/internal/services"
	"wout/internal/utils/logger"
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

// func (h *NotebookHTTPHandler) GetNotebookByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	notebook, err := h.Service.GetNotebookByID(id)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed to get notebook: %v", err))
// 		http.Error(w, "Failed to get notebook", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(notebook)
// }

// func (h *NotebookHTTPHandler) DeleteNotebookByID(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	err := h.Service.DeleteNotebookByID(id)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed to delete notebook: %v", err))
// 		http.Error(w, "Failed to delete notebook", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *NotebookHTTPHandler) UpdateNotebook(w http.ResponseWriter, r *http.Request) {

// 	var notebook models.Notebook
// 	err := json.NewDecoder(r.Body).Decode(&notebook)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	err = h.Service.UpdateNotebook(&notebook)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed to update notebook: %v", err))
// 		http.Error(w, "Failed to update notebook", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(notebook)
// }
