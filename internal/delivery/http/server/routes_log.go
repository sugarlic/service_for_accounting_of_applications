package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// logRoutes prints all registered routes at startup for easier debugging.
func logRoutes(r *gin.Engine, logger *zap.Logger) {
	routes := r.Routes()
	logger.Info("routes registered", zap.Int("count", len(routes)))
	for _, rt := range routes {
		logger.Info("route",
			zap.String("method", rt.Method),
			zap.String("path", rt.Path),
			zap.String("handler", rt.Handler),
		)
	}
}
