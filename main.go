package main

import (
	"lab1/container"
	_ "lab1/docs"
	"lab1/handlers"
	"lab1/middleware"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
// @title Library API
// @version 1.0
// @description REST API for library management with SQLite database
// @host localhost:8080
// @BasePath /
func main() {
	c, err := container.NewContainer("library.db", "config.json")
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}
	defer c.Close()

	booksHandler := handlers.NewBooksHandler(c.BookRepository, c.Validator, c.Config)
	readersHandler := handlers.NewReadersHandler(c.ReaderRepository, c.Validator, c.Config)
	authHandler := handlers.NewAuthHandler(c.UserRepository, c.Validator)

	r := gin.Default()

	// CORS middleware for frontend
	r.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})

	// serving static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// Public auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(c.UserRepository), authHandler.GetProfile)
	}

	// Protected book routes
	books := r.Group("/books")
	books.Use(middleware.AuthMiddleware(c.UserRepository))
	{
		books.GET("/", booksHandler.GetAll)
		books.POST("/", booksHandler.Create)
		books.DELETE("/", booksHandler.DeleteAll)
		books.GET("/:id", booksHandler.GetByID)
		books.PUT("/:id", booksHandler.Update)
		books.DELETE("/:id", booksHandler.Delete)
	}

	// Protected reader routes
	readers := r.Group("/readers")
	readers.Use(middleware.AuthMiddleware(c.UserRepository))
	{
		readers.GET("/", readersHandler.GetAll)
		readers.POST("/", readersHandler.Create)
		readers.DELETE("/", readersHandler.DeleteAll)
		readers.GET("/:id", readersHandler.GetByID)
		readers.PUT("/:id", readersHandler.Update)
		readers.DELETE("/:id", readersHandler.Delete)
		readers.POST("/:id/books/:bookId", readersHandler.AddCurrentlyReading)
		readers.DELETE("/:id/books/:bookId", readersHandler.RemoveCurrentlyReading)
	}

	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
