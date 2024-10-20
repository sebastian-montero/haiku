package http_handlers

import (
	"encoding/json"
	"net/http"
	"wip/internal/models"
	"wip/internal/services"

	"github.com/gorilla/mux"
)

type NotebookHTTPHandler struct {
	Service *services.NotebookService
}

func (h *NotebookHTTPHandler) CreateNotebook(w http.ResponseWriter, r *http.Request) {
	var notebook models.Notebook

	if err := json.NewDecoder(r.Body).Decode(&notebook); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateNotebook(&notebook); err != nil {
		http.Error(w, "Failed to register notebook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notebook)
}

func (h *NotebookHTTPHandler) GetNotebookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var notebook models.Notebook = h.Service.GetNotebookByID(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notebook)
}

func (h *NotebookHTTPHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.Service.DeleteUserByID(id)

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
