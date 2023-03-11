package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func DbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./wingo.db")
	if err != nil {
		zap.S().Fatal("Database connection failed: ", err)
	}

	return db
}
