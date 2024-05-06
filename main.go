package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Database connection string
	envFile := filepath.Join("configs", ".env")
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), 5433, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fmt.Println(psqlInfo)
	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}
	fmt.Println("Connected to the database")

	err = createTableIfNotExists(db)
	if err != nil {
		log.Fatal("Error creating user table database:", err)
	}

	//// Define the HTTP handler function for /v1/api/bulkApproval route
	//http.HandleFunc("/v1/api/bulkApproval", bulkApprovalHttp.BulkApproval)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTableIfNotExists(db *sql.DB) error {

	uuidExtensionQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	// Execute the create table query
	_, err := db.Exec(uuidExtensionQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return err
	}

	createUserTableQuery := `
		CREATE TABLE IF NOT EXISTS user_data (
			userID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			online_status VARCHAR(20),
			friends VARCHAR(255)[] DEFAULT '{}',
			sent_request VARCHAR(255)[] DEFAULT '{}',
			received_request VARCHAR(255)[] DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);`

	// Execute the create table query
	_, err = db.Exec(createUserTableQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return err
	}

	return nil

}
