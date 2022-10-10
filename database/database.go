package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo_micro_db")
	if err != nil {
		panic(err.Error())
	}

	DB = db
}
