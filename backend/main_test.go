package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestDatabaseConnection(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/todo_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to create database connection object: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v. Please ensure the database 'todo_test' exists and the connection string is correct.", err)
	}
}
