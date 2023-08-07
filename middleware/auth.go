package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(enabled bool, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "No authorization header",
			})
			return
		}

		if token != key {
			c.AbortWithStatusJSON(403, gin.H{
				"message": "Unauthorized, wrong token!",
			})
			return
		}
		c.Next()
	}
}
