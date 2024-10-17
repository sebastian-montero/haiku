package main

import (
	"database/sql"
	"log"
	"net/http"
	psql "wip/internal/conn"
	initalizer "wip/internal/init"
	"wip/internal/utils"
)

func main() {
	cfg_path := "./cfg/dev.yaml"
	config := utils.LoadConfig(cfg_path)

	var db *sql.DB = psql.DBManager(config)

	userHandler := initalizer.InitUser(db)
	http.HandleFunc("/users", userHandler.CreateUser)

	log.Print("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer db.Close()
}
