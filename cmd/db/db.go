package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vickon16/third-backend-tutorial/cmd/utils"
)

func NewMySQLStorage() (*sql.DB, error) {
	db, err := sql.Open("mysql", utils.Configs.DB_URL)
	if err != nil {
		return nil, err
	}

	// initialize database
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("MySQL-DB: Connected Successfully")
	return db, nil
}
