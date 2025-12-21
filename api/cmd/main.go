package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"linkbio/internal/config"
	"linkbio/internal/database"
	"linkbio/internal/handler"
	"linkbio/internal/middleware"
	"linkbio/internal/repo"
	"linkbio/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// Repos
	userRepo := repo.NewUserRepo(db)
	pageRepo := repo.NewPageRepo(db)
	blockRepo := repo.NewBlockRepo(db)
	themeRepo := repo.NewThemeRepo(db)
	domainRepo := repo.NewDomainRepo(db)
	bioRepo := repo.NewBioRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	pageService := service.NewPageService(pageRepo, blockRepo)
	themeService := service.NewThemeService(themeRepo)
	compilerService := service.NewCompilerService(pageRepo, blockRepo, themeRepo, userRepo)
	bioService := service.NewBioService(bioRepo, pageRepo, blockRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	pageHandler := handler.NewPageHandler(pageService, compilerService)
	themeHandler := handler.NewThemeHandler(themeService)
	publicHandler := handler.NewPublicHandler(pageRepo, domainRepo)
	bioHandler := handler.NewBioHandler(bioService)

	// Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
	}))

	// Public routes
	app.Get("/r", publicHandler.Render)
	app.Post("/r/password", publicHandler.VerifyPassword)

	// API routes
	api := app.Group("/api")

	// Auth
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)
	api.Get("/auth/me", middleware.Auth(cfg.JWTSecret), authHandler.Me)
	api.Get("/auth/check-username", authHandler.CheckUsername)

	// Protected routes
	protected := api.Group("", middleware.Auth(cfg.JWTSecret))

	// Username setup
	protected.Post("/auth/username", authHandler.SetUsername)

	// Bio (blocks + groups + links)
	protected.Get("/bio", bioHandler.Get)
	protected.Post("/bio/blocks", bioHandler.AddBlock)
	protected.Put("/bio/blocks/:id", bioHandler.UpdateBlock)
	protected.Delete("/bio/blocks/:id", bioHandler.DeleteBlock)
	protected.Post("/bio/blocks/reorder", bioHandler.ReorderBlocks)
	protected.Post("/bio/links", bioHandler.AddLink)
	protected.Put("/bio/links/:id", bioHandler.UpdateLink)
	protected.Delete("/bio/links/:id", bioHandler.DeleteLink)
	protected.Put("/bio/profile", bioHandler.UpdateProfile)
	protected.Put("/bio/social", bioHandler.UpdateSocialLinks)

	// Pages
	protected.Get("/pages", pageHandler.List)
	protected.Post("/pages", pageHandler.Create)
	protected.Get("/pages/:id", pageHandler.Get)
	protected.Get("/pages/:id/draft", pageHandler.GetDraft)
	protected.Post("/pages/:id/save", pageHandler.Save)
	protected.Post("/pages/:id/publish", pageHandler.Publish)
	protected.Delete("/pages/:id", pageHandler.Delete)

	// Themes
	protected.Get("/themes/presets", themeHandler.ListPresets)
	protected.Get("/themes/custom", themeHandler.GetUserCustomTheme)
	protected.Post("/themes/custom", themeHandler.CreateCustom)
	protected.Delete("/themes/custom/:id", themeHandler.DeleteCustom)
	protected.Post("/themes/apply", themeHandler.Apply)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
