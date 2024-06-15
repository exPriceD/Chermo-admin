package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var JwtKey = []byte(os.Getenv("SECRET_KEY"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("museum_id", claims.MuseumID)
		c.Next()
	}
}
