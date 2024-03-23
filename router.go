// router.go

package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// InitRouter initializes the HTTP request handlers and starts the HTTP server
func InitRouter() {

	// Create a new Gin router
	router := gin.Default()

	// Create a new cookie store
	store := cookie.NewStore([]byte("93wt4nvp2y34b223r2wet43"))

	// Use the cookie store for the session middleware
	router.Use(sessions.Sessions("session", store))

	// Load HTML files
	router.LoadHTMLGlob("templates/*")

	// Serve static files
	router.Static("/static", "./static")

	// Register the HTTP request handlers
	registerHandlers(router)

	// Start the HTTP server on port 8080
	router.Run(":8080")

}

// HTTP Request Handlers
func registerHandlers(router *gin.Engine) {

	// Get authorized credentials
	authorizedCredentials := getAuthorizedCredentials()

	// Apply basic auth middleware
	authorized := router.Group("/", gin.BasicAuth(authorizedCredentials))

	// Register
	authorized.GET("/register", registerPage)
	authorized.POST("/register", register)
	// Login
	router.GET("/login", loginPage)
	router.POST("/login", login)
	// Logout
	router.GET("/logout", logout)
	// Home
	router.GET("/", indexPage)
	router.GET("/items", items)
	router.POST("/items/create", createItem)
	router.GET("/items/edit/:id", editItem)
	router.POST("/items/update/:id", updateItem)
	router.DELETE("/items/delete/:id", deleteItem)
	// Health check
	router.GET("/ping", ping)
}

func getAuthorizedCredentials() gin.Accounts {
	// Load the .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve credentials
	adminUser := os.Getenv("ADMIN_USER")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Define authorized credentials using the loaded environment variables
	authorizedCredentials := gin.Accounts{
		adminUser: adminPassword,
	}

	return authorizedCredentials
}
