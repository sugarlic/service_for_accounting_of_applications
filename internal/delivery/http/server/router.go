package server

import (
	"application-service/internal/config"
	handlers "application-service/internal/delivery/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter wires up the Gin engine with middlewares, CORS and routes.
func NewRouter(
	cfg *config.Config,
	logger *zap.Logger,
	applicationHandler *handlers.ApplicationHandler,
	integrationHandler *handlers.IntegrationHandler) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// Limit request body to 1 MB to prevent abuse.
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	r.Static("/static", "./web")

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	admin := r.Group("/admin")
	{
		admin.POST("/applications", applicationHandler.CreateApplication)
		admin.GET("/applications", applicationHandler.ListApplications)
		admin.GET("/applications/:id", applicationHandler.GetApplication)
		admin.PATCH("/applications/:id/status", applicationHandler.UpdateApplicationStatus)
	}

	integration := r.Group("/api/v1")
	integration.Use(APIKeyMiddleware(cfg.CRMAPIKey))
	{
		integration.POST("/applications", integrationHandler.CreateApplication)
		integration.GET("/applications", integrationHandler.ListApplications)
	}

	logRoutes(r, logger)
	return r
}
