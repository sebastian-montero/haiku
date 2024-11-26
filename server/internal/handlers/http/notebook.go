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

func (h *NotebookHTTPHandler) GetNotebooksByOwnerId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerId := vars["owner_id"]

	notebooks, err := h.Service.GetNotebooksByOwnerId(ownerId)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get notebooks: %v", err))
		http.Error(w, "Failed to get notebooks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notebooks)
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

func (h *NotebookHTTPHandler) UpdateNotebook(w http.ResponseWriter, r *http.Request) {

	var notebook models.Notebook
	err := json.NewDecoder(r.Body).Decode(&notebook)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to decode request body: %v", err))
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateNotebook(&notebook)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update notebook: %v", err))
		http.Error(w, "Failed to update notebook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notebook)
}
