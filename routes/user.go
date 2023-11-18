package routes

import (
	"log/slog"
	"photo-app/controllers"
	"photo-app/handlers"
	"photo-app/middlewares"
	"photo-app/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserRoutes(r *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	userRepo := repositories.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo, logger)
	userHandler := handlers.NewUserHandler(userController)

	{
		r.POST("/register", userHandler.Register)
		r.POST("/login", userHandler.Login)
		r.Use(middlewares.AuthMiddleware(true))
		r.PUT("/me", userHandler.Update)
		r.DELETE("/me", userHandler.Delete)
	}
}
