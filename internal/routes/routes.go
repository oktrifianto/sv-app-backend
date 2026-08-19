package routes

import (
	"sv-app/backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// Setup registers all routes and middleware on the Fiber app.
func Setup(app *fiber.App, db *gorm.DB) {
	// Allow the React dev server to call the API during local development.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // React dev server
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	h := &handlers.ArticleHandler{DB: db}

	article := app.Group("/article")
	article.Post("/", h.CreateArticle)

	// List with pagination: /article/:limit/:offset
	// This must be registered before /:id so Fiber doesn't swallow "10" as an id.
	article.Get("/:limit/:offset", h.ListArticles)

	article.Get("/:id", h.GetArticle)
	article.Put("/:id", h.UpdateArticle)
	article.Delete("/:id", h.DeleteArticle)
}
