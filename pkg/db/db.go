package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/turanoo/bitebattle/pkg/config"
)

var db *sql.DB

func Init(cfg *config.Config) error {
	instanceConnName := cfg.DB.InstanceConn
	user := cfg.DB.User
	password := cfg.DB.Pass
	dbName := cfg.DB.Name
	host := cfg.DB.Host
	port := cfg.DB.Port

	var connStr string

	if instanceConnName != "" {
		connStr = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable",
			user, password, dbName, instanceConnName,
		)
		log.Print("[DB DEBUG] Using Cloud SQL Unix socket.")
	} else {
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbName,
		)
		log.Print("[DB DEBUG] Using TCP")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("[DB ERROR] sql.Open failed: %v", err)
		return err
	}

	if err = db.Ping(); err != nil {
		log.Printf("[DB ERROR] db.Ping failed: %v", err)
		return err
	}

	log.Println("Connected to PostgreSQL database.")
	return nil
}

func GetDB() *sql.DB {
	return db
}

// GetPostgresURL returns a postgres:// URL for use with golang-migrate
func GetPostgresURL(cfg *config.Config) string {
	user := cfg.DB.User
	password := cfg.DB.Pass
	dbName := cfg.DB.Name
	instanceConnName := cfg.DB.InstanceConn
	host := cfg.DB.Host
	port := cfg.DB.Port

	if instanceConnName != "" {
		// Cloud SQL Unix socket
		return fmt.Sprintf(
			"postgres://%s:%s@/%s?host=/cloudsql/%s&sslmode=disable",
			user, password, dbName, instanceConnName,
		)
	}

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)
}
