// Description: This is the main file for the server. It initializes the server and the handlers for the routes.
package main

import (
	"database/sql"
	"log"
	"net/http"
	psql "wout/internal/conn"
	initalizer "wout/internal/init"
	"wout/internal/utils"

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
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.DeleteUserByID).Methods("DELETE")

	r.HandleFunc("/notebooks", notebookHandler.CreateNotebook).Methods("POST")
	r.HandleFunc("/notebooks", notebookHandler.UpdateNotebook).Methods("PUT")
	r.HandleFunc("/notebooks/{id}", notebookHandler.GetNotebookByID).Methods("GET")
	r.HandleFunc("/notebooks/by_owner/{owner_id}", notebookHandler.GetNotebooksByOwnerId).Methods("GET")
	r.HandleFunc("/notebooks/{id}", notebookHandler.DeleteNotebookByID).Methods("DELETE")

	r.HandleFunc("/sessions", sessionHTTPHandler.CreateSession).Methods("POST")
	r.HandleFunc("/sessions", sessionHTTPHandler.GetActiveSessions).Methods("GET")
	r.HandleFunc("/sessions/{id}", sessionHTTPHandler.GetSessionByID).Methods("GET")
	r.HandleFunc("/sessions/{id}", sessionHTTPHandler.DeleteSessionByID).Methods("DELETE")
	r.HandleFunc("/sessions/{id}/end", sessionHTTPHandler.EndSessionByID).Methods("PUT")

	r.HandleFunc("/content", contentHandler.CreateContent).Methods("POST")
	r.HandleFunc("/content/by_session/{session_id}", contentHandler.GetLatestContentBySessionId).Methods("GET")

	// `ws://localhost:8080/ws/${notebookID}?owner_id=${ownerID}`
	r.HandleFunc("/ws/{notebook_id}", sessionWSHandler.WebSocketEndpoint).Methods("GET")

	log.Print("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(r)))
	defer db.Close()
}

// withCORS function sets up the CORS headers for the server.
// It allows requests from the localhost.
// It also sets up the allowed methods and headers for the server.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
