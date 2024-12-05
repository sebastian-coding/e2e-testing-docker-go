package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserEndpoint(t *testing.T) {
	// Set environment variables needed for the test
	err := os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/testdb?sslmode=disable")
	assert.NoError(t, err)

	// Set up the Gin router in test mode
	gin.SetMode(gin.TestMode)

	// Initialize dependencies
	db, err := InitializeDependencies()
	if err != nil {
		t.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer db.Close()

	truncateTable(t, db)

	router := SetupRouter(db)

	t.Run("Responds OK", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(router)
		defer server.Close()

		// Prepare payload
		userPayload := `{"name": "John Doe", "email": "john.doe@example.com"}`

		// Send POST request
		resp, err := http.Post(server.URL+"/users", "application/json", bytes.NewBuffer([]byte(userPayload)))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	truncateTable(t, db)
}

func truncateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE users")
	assert.NoError(t, err)
}
