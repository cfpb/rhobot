package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func GetPGConnection(uri string) *sql.DB {

	if uri == "" {
		// Database connection info
		host := os.Getenv("PGHOST")
		db := os.Getenv("PGDATABASE")
		user := os.Getenv("PGUSER")
		pass := os.Getenv("PGPASSWORD")
		uri = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full", user, pass, host, db)

	}
	fmt.Println("DATABASE URI ", uri)
	cxn, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}
	return cxn
}
