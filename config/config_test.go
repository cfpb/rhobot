package config

import (
	"fmt"
	"testing"
)

func TestGetDBURI(t *testing.T) {
	config := NewDefaultConfig()

	if config.DBURI() != "postgres://postgres:password@localhost:5432/postgres?sslmode=require" {
		fmt.Printf("Testing GetDBURI failed: %s", config.DBURI())
		t.Fail()
	}
}

func TestSetDBURI(t *testing.T) {
	config := NewDefaultConfig()
	config.SetDBURI("postgres://test_user:password@localhost:5432/postgres?sslmode=require")

	if config.pgUser != "test_user" {
		fmt.Printf("Testing SetDBURI failed, user: %s", config.pgUser)
		t.Fail()
	}
}
