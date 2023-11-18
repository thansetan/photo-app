package app

import (
	"fmt"
	"log/slog"
	"photo-app/helpers"
	"photo-app/routes"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type app struct {
	port   uint
	db     *gorm.DB
	r      *gin.Engine
	logger *slog.Logger
}

func New(conf helpers.App, db *gorm.DB, logger *slog.Logger) *app {
	return &app{
		port:   conf.Port,
		db:     db,
		r:      gin.Default(),
		logger: logger,
	}
}

func (app *app) Start() error {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	app.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := app.r.Group("/api/v1")

	users := v1.Group("/users")
	{
		routes.NewUserRoutes(users, app.db, app.logger)
	}

	photosApi := v1.Group("/photos")
	photosStatic := app.r.Group("/photos")
	{
		routes.NewPhotoRoutes(photosApi, photosStatic, app.db, app.logger)
	}

	app.logger.Info("Server starting", "port", app.port)
	if err := app.r.Run(fmt.Sprintf(":%d", app.port)); err != nil {
		return err
	}

	return nil
}
