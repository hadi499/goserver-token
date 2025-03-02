package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    //ambil token dari header Authorization
    tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
      c.Abort()
      return 
    }

    // periksa apa ada token di blacklist
    if IsTokenBlacklisted(tokenString) {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalidated"})
      c.Abort()
      return 
    }

    c.Next()

  }
}
