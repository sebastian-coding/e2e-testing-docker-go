package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// InitializeDependencies initializes the database and Kafka writer
func InitializeDependencies() (*sql.DB, error) {
	// Database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SetupRouter initializes the Gin router and configures endpoints
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// POST /users endpoint
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		user.ID = uuid.New().String()

		// Store user in database
		_, err := db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", user.ID, user.Name, user.Email)
		if err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created"})
	})

	return router
}

func main() {
	// Initialize dependencies
	db, err := InitializeDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer db.Close()

	// Set up the router with dependencies
	router := SetupRouter(db)

	// Start the server
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
