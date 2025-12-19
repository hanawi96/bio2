package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"linkbio/api/internal/modules/pages"
	"linkbio/api/internal/modules/public"
)

func Setup(app *fiber.App, pool *pgxpool.Pool) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public routes (no /api prefix)
	publicRepo := public.NewRepo(pool)
	publicService := public.NewService(publicRepo)
	publicHandler := public.NewHandler(publicService)

	app.Get("/r", publicHandler.GetPage)
	app.Post("/r/password", publicHandler.VerifyPassword)

	// API routes
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "LinkBio API v1"})
	})

	// Pages module
	pagesRepo := pages.NewRepo(pool)
	pagesService := pages.NewService(pagesRepo)
	pagesHandler := pages.NewHandler(pagesService)

	api.Post("/pages", pagesHandler.CreatePage)
	api.Get("/pages", pagesHandler.ListPages)
	api.Get("/pages/:id/draft", pagesHandler.GetDraft)
	api.Get("/pages/:id/published", pagesHandler.GetPublished)
	api.Post("/pages/:id/save", pagesHandler.SaveDraft)
	api.Post("/pages/:id/publish", pagesHandler.Publish)
}
