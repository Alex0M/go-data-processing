package middleware

import (
	"out/internal/utils"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

//var secretKey = os.Getenv("SECRET_KEY")

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		err := utils.ExtractTokenID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
