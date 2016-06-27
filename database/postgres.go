package database

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"

	_ "github.com/lib/pq" // Blank import required
)

// GetPGConnection returns a connection to the postgres database from the URI
func GetPGConnection(uri string) *sql.DB {
	cxn, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}
	return cxn
}
