package main

import (
	"core_two_go/config"
	"core_two_go/controllers"
	"core_two_go/database"
	"core_two_go/services"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	err = createUsersTable(db)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("Users table created or already exists")

	userService := services.NewUserService(db)
	userHandler := controllers.NewUserHandler(userService)

	http.HandleFunc("/users" /* middleware.Auth( */, userHandler.HandleUsers)

	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY,
		name TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}
