package handler

import (
	"net/http"
	"strings"

	"ozinse/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code": "UNAUTHORIZED",
				"message":    "Требуется авторизация",
				"details":    nil,
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtService.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code": "TOKEN_INVALID",
				"message":    "Недействительный токен",
				"details":    nil,
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists || roleID.(int) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error_code": "FORBIDDEN",
				"message":    "Доступ запрещён. Требуются права администратора.",
				"details":    nil,
			})
			return
		}
		c.Next()
	}
}