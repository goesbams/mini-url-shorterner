package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error

	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)

		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)

		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Connected to database")
	return DB, nil
}
