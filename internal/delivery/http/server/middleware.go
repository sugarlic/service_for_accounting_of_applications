package server

import "github.com/gin-gonic/gin"

func APIKeyMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		if apiKey != expectedKey {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
