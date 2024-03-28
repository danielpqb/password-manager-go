package database_manager

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) ConnectDB() {
	instance, err := sql.Open("sqlite", "internal/database/db.sqlite")
	if err != nil {
		panic(errors.New("ConnectDB (15)" + err.Error()))
	}

	_, err = instance.Exec("CREATE TABLE IF NOT EXISTS users (name TEXT, password TEXT);")
	if err != nil {
		panic(err)
	}
	_, err = instance.Exec("CREATE TABLE IF NOT EXISTS key_values (key TEXT, value TEXT);")
	if err != nil {
		panic(err)
	}

	db.Db = instance
}
