package v1

import (
	_ "EffectiveMobile/docs"
	"EffectiveMobile/pkg/postgres"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// NewRouter -.
// @title       AUTO API
// @description API of the AUTO service
// @version     1.0
// @host        localhost:8080
// @BasePath    /
func NewRouter(handler *gin.Engine, l *zap.Logger, db postgres.Service) {
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	newAutoRoutes(handler, l, db)
}
