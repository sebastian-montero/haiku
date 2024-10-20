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


	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.DeleteUserByID).Methods("DELETE")
	

	log.Print("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
	defer db.Close()
}
