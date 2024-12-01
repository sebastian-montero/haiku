// Description: This is the main file for the server. It initializes the server and the handlers for the routes.
package main

import (
	"database/sql"
	psql "haiku/internal/conn"
	initalizer "haiku/internal/init"
	"haiku/internal/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// main function initializes the server and the handlers for the routes.
// It also initializes the database connection.
// It also sets up the CORS headers for the server.
func main() {
	var config utils.Config = utils.LoadConfig("./cfg/dev.yaml")

	var db *sql.DB = psql.DBManager(&config)

	userHandler := initalizer.InitUser(db)
	notebookHandler := initalizer.InitNotebook(db)
	contentHandler := initalizer.InitContent(db)
	sessionHTTPHandler, sessionWSHandler := initalizer.InitSession(db)

	r := mux.NewRouter()

	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/signup", userHandler.SignUp).Methods("POST")

	api := r.PathPrefix("/").Subrouter()
	api.Use(utils.AuthMiddleware)

	api.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.DeleteUserByID).Methods("DELETE")

	api.HandleFunc("/notebooks", notebookHandler.CreateNotebook).Methods("POST")
	api.HandleFunc("/notebooks", notebookHandler.UpdateNotebook).Methods("PUT")
	api.HandleFunc("/notebooks/{id}", notebookHandler.GetNotebookByID).Methods("GET")
	api.HandleFunc("/notebooks/by_owner/{owner_id}", notebookHandler.GetNotebooksByOwnerId).Methods("GET")
	api.HandleFunc("/notebooks/{id}", notebookHandler.DeleteNotebookByID).Methods("DELETE")

	api.HandleFunc("/sessions", sessionHTTPHandler.CreateSession).Methods("POST")
	api.HandleFunc("/sessions", sessionHTTPHandler.GetActiveSessions).Methods("GET")
	api.HandleFunc("/sessions/by_notebook/{notebook_id}", sessionHTTPHandler.GetSessionByNotebookId).Methods("GET")
	api.HandleFunc("/sessions/{id}", sessionHTTPHandler.GetSessionByID).Methods("GET")
	api.HandleFunc("/sessions/{id}", sessionHTTPHandler.DeleteSessionByID).Methods("DELETE")
	api.HandleFunc("/sessions/{id}/end", sessionHTTPHandler.EndSessionByID).Methods("PUT")
	api.HandleFunc("/sessions/{id}", sessionHTTPHandler.UpdateSession).Methods("PUT")

	api.HandleFunc("/content", contentHandler.CreateContent).Methods("POST")
	api.HandleFunc("/content/by_session/{session_id}", contentHandler.GetLatestContentBySessionId).Methods("GET")

	r.HandleFunc("/ws/{conn_type}/{notebook_id}", sessionWSHandler.WebSocketEndpoint).Methods("GET")

	log.Print("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(r)))
	defer db.Close()
}

// withCORS function sets up the CORS headers for the server.
// It allows requests from the localhost.
// It also sets up the allowed methods and headers for the server.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigin := "https://haiku.incendia.dev"

		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
