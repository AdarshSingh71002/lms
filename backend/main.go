package main

import (
	"log"
	"os"

	"lms/backend/config"
	"lms/backend/handlers"
	"lms/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database.
	config.ConnectDatabase()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/signup", handlers.Signup)
	router.POST("/library/create", handlers.CreateLibrary)
	// Global authentication middleware (checks for X-User-Email header).
	router.Use(middlewares.AuthMiddleware)

	// Library Owner Flow: Create library (and Owner).

	// Owner-only endpoint to onboard a LibraryAdmin.
	router.POST("/owner/admin/create", handlers.CreateAdmin)

	// Admin routes: accessible by Owner or LibraryAdmin.
	adminGroup := router.Group("/admin")
	adminGroup.Use(middlewares.AdminMiddleware)
	{
		adminGroup.POST("/books", handlers.AddBook)
		adminGroup.DELETE("/books/:isbn", handlers.RemoveBook)
		adminGroup.PUT("/books/:isbn", handlers.UpdateBook)
		adminGroup.GET("/requests", handlers.ListIssueRequests)
		adminGroup.POST("/requests/:reqid/approve", handlers.ApproveIssueRequest)
		adminGroup.POST("/requests/:reqid/reject", handlers.RejectIssueRequest)
	}

	// Reader routes.
	readerGroup := router.Group("/reader")
	readerGroup.Use(middlewares.ReaderMiddleware)
	{
		readerGroup.GET("/books", handlers.SearchBooks)
		readerGroup.POST("/request", handlers.RaiseIssueRequest)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	router.Run(":" + port)
}
