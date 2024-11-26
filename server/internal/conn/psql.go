package psql

import (
	"database/sql"
	"fmt"
	utilities "haiku/internal/utils"
	"haiku/internal/utils/logger"

	_ "github.com/lib/pq"
)

func DBManager(cfg *utilities.Config) *sql.DB {
	logger.Info("Creating db conn...")

	dbConn := cfg.DBConn
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConn.Host,
		dbConn.Port,
		dbConn.User,
		dbConn.Password,
		dbConn.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Failed to ping db")
		panic(err)
	}

	for _, query := range cfg.DBTables {
		createTables(db, query)
	}

	logger.Info("Created db conn ok...")
	return db
}

func createTables(db *sql.DB, query string) {
	logger.Info(fmt.Sprintf("Executing query %s", query))
	_, err := db.Exec(query)

	if err != nil {
		logger.Error("Failed to execute query")
		panic(err)
	}

	logger.Info("Executed query OK...")
}
