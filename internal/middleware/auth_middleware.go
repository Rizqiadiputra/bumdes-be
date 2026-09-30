package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/liyansasongko/bumdes-be/internal/utils"
)

const ContextUserIDKey = "user_id"
const ContextRoleIDKey = "role_id"

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "token tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "format token tidak valid")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(jwtSecret, parts[1])
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "token tidak valid atau kadaluarsa")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextRoleIDKey, claims.RoleID)
		c.Next()
	}
}
