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

type NotebookHTTPHandler struct {
	Service *services.NotebookService
}

func (h *NotebookHTTPHandler) CreateNotebook(w http.ResponseWriter, r *http.Request) {
	var notebook models.Notebook

	if err := json.NewDecoder(r.Body).Decode(&notebook); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateNotebook(&notebook); err != nil {
		logger.Error(fmt.Sprintf("Failed to create notebook: %v", err))
		http.Error(w, "Failed to create notebook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notebook)
}

func (h *NotebookHTTPHandler) GetNotebookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	notebook, err := h.Service.GetNotebookByID(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get notebook: %v", err))
		http.Error(w, "Failed to get notebook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notebook)
}

func (h *NotebookHTTPHandler) DeleteNotebookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.Service.DeleteNotebookByID(id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete notebook: %v", err))
		http.Error(w, "Failed to delete notebook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func (h *NotebookHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

// 	var user models.User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	h.Service.UpdateUser(&user)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(user)
// }
