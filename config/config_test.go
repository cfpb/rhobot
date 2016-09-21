package config

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestGetDBURI(t *testing.T) {
	config := NewDefaultConfig()

	if config.DBURI() != "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" {
		log.Errorf("Testing GetDBURI failed: %s", config.DBURI())
		t.Fail()
	}
}

func TestSetDBURI(t *testing.T) {
	config := NewDefaultConfig()
	config.SetDBURI("postgres://test_user:password@localhost:5432/postgres?sslmode=require")

	if config.PgUser != "test_user" {
		log.Errorf("Testing SetDBURI failed, user: %s", config.PgUser)
		t.Fail()
	}
}
