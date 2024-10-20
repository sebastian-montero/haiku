package main

import (
	"database/sql"
	"log"
	"net/http"
	psql "wip/internal/conn"
	initalizer "wip/internal/init"
	"wip/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	cfg_path := "./cfg/dev.yaml"
	config := utils.LoadConfig(cfg_path)

	var db *sql.DB = psql.DBManager(config)

	userHandler := initalizer.InitUser(db)
	notebookHandler := initalizer.InitNotebook(db)

	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.DeleteUserByID).Methods("DELETE")

	r.HandleFunc("/notebooks", notebookHandler.CreateNotebook).Methods("POST")
	// r.HandleFunc("/notebooks", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/notebooks/{id}", notebookHandler.GetNotebookByID).Methods("GET")
	// r.HandleFunc("/notebooks/{id}", userHandler.DeleteNotebookByID).Methods("DELETE")

	log.Print("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(r)))
	defer db.Close()
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
