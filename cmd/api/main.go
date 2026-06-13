package main

import (
	"log"

	"ozinse/internal/config"
	"ozinse/internal/database"
	"ozinse/internal/handler"
	"ozinse/internal/repository"
	"ozinse/internal/service"
	"ozinse/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwtService := jwt.NewService(cfg.JWTSecret, cfg.RefreshSecret)

	// Repos
	userRepo := repository.NewUserRepo(db)
	tokenRepo := repository.NewTokenRepo(db)
	projectRepo := repository.NewProjectRepo(db)
	refRepo := repository.NewReferenceRepo(db)
	favRepo := repository.NewFavoriteRepo(db)
	adminRepo := repository.NewAdminRepo(db)

	// Services
	authService := service.NewAuthService(userRepo, tokenRepo, jwtService, cfg.BaseURL)
	profileService := service.NewProfileService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	refService := service.NewReferenceService(refRepo)
	favService := service.NewFavoriteService(favRepo)
	adminService := service.NewAdminService(adminRepo, cfg.BaseURL)
	
	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	contentHandler := handler.NewContentHandler(projectService)
	refHandler := handler.NewReferenceHandler(refService)
	favHandler := handler.NewFavoriteHandler(favService)
	adminHandler := handler.NewAdminHandler(adminService)

	r := gin.Default()
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Public
		api.GET("/projects", contentHandler.GetProjects)
		api.GET("/projects/featured", contentHandler.GetFeatured)
		api.GET("/projects/:id", contentHandler.GetProjectByID)
		api.GET("/projects/:id/seasons", contentHandler.GetSeasons)
		api.GET("/categories", refHandler.GetCategories)
		api.GET("/genres", refHandler.GetGenres)
		api.GET("/age-ratings", refHandler.GetAgeRatings)

		// Protected
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware(jwtService))
		{
			protected.GET("/profile/me", profileHandler.GetProfile)
			protected.PUT("/profile/me", profileHandler.UpdateProfile)
			protected.PUT("/profile/password", profileHandler.ChangePassword)

			protected.GET("/favorites", favHandler.GetFavorites)
			protected.POST("/favorites/:project_id", favHandler.AddFavorite)
			protected.DELETE("/favorites/:project_id", favHandler.RemoveFavorite)
		}

		// Admin
		admin := api.Group("/admin")
		admin.Use(handler.AuthMiddleware(jwtService), handler.AdminMiddleware())
		{
			// Projects
			admin.GET("/projects", contentHandler.GetProjects)
			admin.POST("/projects", adminHandler.CreateProject)
			admin.PUT("/projects/:id", adminHandler.UpdateProject)
			admin.DELETE("/projects/:id", adminHandler.DeleteProject)
			admin.POST("/projects/:id/seasons", adminHandler.CreateSeason)
			admin.POST("/seasons/:season_id/episodes", adminHandler.CreateEpisode)
			admin.PUT("/projects/featured-order", adminHandler.UpdateFeaturedOrder)

			// Categories
			admin.POST("/categories", adminHandler.CreateCategory)
			admin.PUT("/categories/:id", adminHandler.UpdateCategory)
			admin.DELETE("/categories/:id", adminHandler.DeleteCategory)

			// Genres
			admin.POST("/genres", adminHandler.CreateGenre)
			admin.PUT("/genres/:id", adminHandler.UpdateGenre)
			admin.DELETE("/genres/:id", adminHandler.DeleteGenre)

			// Age Ratings
			admin.POST("/age-ratings", adminHandler.CreateAgeRating)
			admin.PUT("/age-ratings/:id", adminHandler.UpdateAgeRating)
			admin.DELETE("/age-ratings/:id", adminHandler.DeleteAgeRating)

			// Users
			admin.GET("/users", adminHandler.GetUsers)
			admin.POST("/users/:user_id/assign-role", adminHandler.AssignRole)

			// Roles
			admin.GET("/roles", adminHandler.GetRoles)
			admin.POST("/roles", adminHandler.CreateRole)
			admin.PUT("/roles/:id", adminHandler.UpdateRole)
			admin.DELETE("/roles/:id", adminHandler.DeleteRole)
			
			// Upload
			admin.POST("/upload", adminHandler.UploadFile)
		}
	}

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	r.Run(":" + cfg.Port)
}