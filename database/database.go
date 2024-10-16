package database

import (
	"core_two_go/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to db: %v", err)
	}

	log.Println("Connected to the PostgreSQL database!")
	return db, nil
}
