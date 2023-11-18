package routes

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"photo-app/controllers"
	"photo-app/handlers"
	"photo-app/middlewares"
	"photo-app/repositories"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPhotoRoutes(api *gin.RouterGroup, static *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	repo := repositories.NewPhotoRepository(db)
	userRepo := repositories.NewUserRepository(db)
	controller := controllers.NewPhotoController(repo, userRepo, logger)
	handler := handlers.NewPhotoHandler(controller)

	{
		re := regexp.MustCompile(`^/photos/[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}/([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
		static.Use(func(ctx *gin.Context) {
			if !re.MatchString(ctx.Request.URL.String()) {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
			ctx.Next()
		}, middlewares.AuthMiddleware(false), func(ctx *gin.Context) {
			photoID := re.FindStringSubmatch(ctx.Request.URL.String())[1]
			isAllowed, err := controller.IsAllowedToView(ctx, photoID)
			if !isAllowed || err != nil {
				ctx.AbortWithStatus(404)
				return
			}
			ctx.Next()
		})

		static.Static("", filepath.Join(os.Getenv("PHOTO_DIR")))
	}

	{
		api.GET("", handler.GetAll)
		api.GET("/by/:username", handler.GetByOwner)
		api.Use(middlewares.AuthMiddleware(true))
		api.GET("/my", handler.GetMine)
		api.POST("", handler.Create)
		api.PUT("/:id", handler.Update)
		api.DELETE("/:id", handler.Delete)
	}
}
