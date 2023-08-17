package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	var db *sql.DB

	connStr := "postgres://ade:password@localhost:5432/chat?sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	stmt := `CREATE TABLE IF NOT EXISTS message 
	(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	content TEXT NOT NULL,
	sender TEXT NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	// stmt := `DROP TABLE IF EXISTS message;`
	fmt.Println(stmt)
	_, err = db.Exec(stmt)

	if err != nil {
		panic(err)
	}
	fmt.Println("done...")

}
