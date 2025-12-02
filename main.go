package main

import (
	"lab1/container"
	_ "lab1/docs"
	"lab1/handlers"
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

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	books := r.Group("/books")
	{
		books.GET("/", booksHandler.GetAll)
		books.POST("/", booksHandler.Create)
		books.DELETE("/", booksHandler.DeleteAll)
		books.GET("/:id", booksHandler.GetByID)
		books.PUT("/:id", booksHandler.Update)
		books.DELETE("/:id", booksHandler.Delete)
	}

	readers := r.Group("/readers")
	{
		readers.GET("/", readersHandler.GetAll)
		readers.POST("/", readersHandler.Create)
		readers.DELETE("/", readersHandler.DeleteAll)
		readers.GET("/:id", readersHandler.GetByID)
		readers.PUT("/:id", readersHandler.Update)
		readers.DELETE("/:id", readersHandler.Delete)
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
