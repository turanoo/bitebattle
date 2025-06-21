package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/turanoo/bitebattle/pkg/config"
	"github.com/turanoo/bitebattle/pkg/logger"
)

var db *sql.DB

func Init(cfg *config.Config) error {
	instanceConnName := cfg.DB.InstanceConn
	user := cfg.DB.User
	password := cfg.DB.Pass
	dbName := cfg.DB.Name
	host := cfg.DB.Host
	port := cfg.DB.Port

	log := logger.Log.WithFields(logrus.Fields{
		"db_user":       user,
		"db_name":       dbName,
		"instance_conn": instanceConnName,
		"db_host":       host,
		"db_port":       port,
	})

	var connStr string

	if instanceConnName != "" {
		connStr = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=/cloudsql/%s sslmode=disable",
			user, password, dbName, instanceConnName,
		)
		log.Info("Using Cloud SQL Unix socket to connect to database.")
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
		log.Info("Using TCP connection to connect to database.")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.WithError(err).Error("sql.Open failed")
		return err
	}

	if err = db.Ping(); err != nil {
		log.WithError(err).Error("db.Ping failed")
		return err
	}

	log.Info("Connected to PostgreSQL database.")
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
