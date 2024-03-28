package database_manager

import "database/sql"

func ConnectDB() error {
	pool, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec("CREATE TABLE IF NOT EXISTS users (name TEXT, password TEXT);")
	if err != nil {
		panic(err)
	}
	_, err = pool.Exec("CREATE TABLE IF NOT EXISTS key_values (key TEXT, value TEXT);")
	if err != nil {
		panic(err)
	}

	return err
}
