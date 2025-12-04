package middleware

import (
	"net/http"
	"strings"

	"loading_time/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware проверяет JWT и допустимые роли
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		// ЛОГИРУЕМ успешный разбор токена
		logrus.Infof("AuthMiddleware: ParseJWT OK -> user_id=%v, role=%v", claims.UserID, claims.Role)

		// сохраняем данные в контекст Gin
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// Проверяем роль, если переданы allowedRoles
		if len(allowedRoles) > 0 {
			allowed := false
			for _, r := range allowedRoles {
				if r == claims.Role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied for role: " + claims.Role})
				return
			}
		}

		c.Next()
	}
}

// ModeratorMiddleware — требует роль "port_operator"
func ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, exists := c.Get("role")

		if !exists {
			logrus.Warn("ModeratorMiddleware: role not found in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Moderator access required"})
			return
		}
		role, ok := roleAny.(string)
		logrus.Infof("ModeratorMiddleware: role from context=%v ok=%v", roleAny, ok)

		if !ok || role != "port_operator" {
			logrus.Warnf("ModeratorMiddleware: access denied for role=%v", role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Moderator access required"})
			return
		}

		logrus.Infof("ModeratorMiddleware: access granted for role=%v", role)
		c.Next()
	}
}
