package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func ValidateToken(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusForbidden, gin.H{"message": "No token provided!"})
			c.Abort()
			return
		}

		args := strings.Split(tokenString, " ")
		if args[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
			c.Abort()
			return
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(args[1], claims, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
			c.Abort()
			return
		}
		c.Set("email", claims["email"])
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
