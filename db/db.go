package db

import (
	"database/sql"
	"os"

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

func DbInit() {
	dbName := "wingo.db"
	_, err := os.Stat(dbName)
	if os.IsNotExist(err) {
		zap.S().Warn("Wingo.db does not exist. Creating new database...")
		os.Create("wingo.db")
	} else {
		zap.S().Info("wingo.db exists. Continuing...")
	}

	conn := DbConnect()
	defer conn.Close()

	sqlStmt := `
	create table if not exists currentusers (username text not null primary key, fullname text, enabled integer, locked integer, admin integer, passwdexpired integer, cantchangepasswd integer, passwdage integer, lastlogon text, badpasswdattempts numeric, numoflogons numeric);
	`
	_, err1 := conn.Exec(sqlStmt)
	if err1 != nil {
		zap.S().Fatal("Database initialization failed: ", err)
	}

}

func CountTableRecords(tableName string) int {
	conn := DbConnect()
	defer conn.Close()

	var count int
	query := "SELECT COUNT(*) FROM " + tableName
	err := conn.QueryRow(query).Scan(&count)
	if err != nil {
		zap.S().Error("Error counting records in database: ", err)
	}

	return count
}
