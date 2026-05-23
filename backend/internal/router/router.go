package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"project-alpha/backend/internal/config"
	"project-alpha/backend/internal/handlers"
	"project-alpha/backend/internal/repositories"
	"project-alpha/backend/internal/response"
	"project-alpha/backend/internal/services"
)

func New(cfg config.Config, db *gorm.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.Use(corsMiddleware(cfg.CORSAllowedOrigins))

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group("/api/v1")
	registerAPIRoutes(api, db)

	return engine
}

func registerAPIRoutes(api *gin.RouterGroup, db *gorm.DB) {
	tagRepository := repositories.NewTagRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	tagService := services.NewTagService(tagRepository)
	ticketService := services.NewTicketService(ticketRepository)
	tagHandler := handlers.NewTagHandler(tagService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	api.GET("", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	tagHandler.RegisterRoutes(api)
	ticketHandler.RegisterRoutes(api)
}
