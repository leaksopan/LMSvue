package main

import (
	"lms-vue-go/backend/config"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Initialize database connection
	dbConfig := config.DefaultConfig()
	err := config.InitDB(dbConfig)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer config.CloseDB()

	// Test if connection is successful
	if config.DB == nil {
		t.Fatal("Database connection is nil")
	}

	// Test if connection is working
	err = config.DB.Ping()
	if err != nil {
		t.Fatalf("Error pinging database: %v", err)
	}

	t.Log("Database connection successful")
}
