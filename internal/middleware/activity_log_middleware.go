package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

// ActivityLog records a "view" log entry for every successful GET request made
// by an authenticated user (e.g. opening /me, /roles, /accounts).
func ActivityLog(logService service.UserLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method != http.MethodGet {
			return
		}
		if c.Writer.Status() >= http.StatusBadRequest {
			return
		}

		userIDVal, ok := c.Get(ContextUserIDKey)
		if !ok {
			return
		}
		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			return
		}

		meta := utils.GetRequestMeta(c)
		module := moduleFromPath(meta.Path)

		if err := logService.LogView(meta, userID, module, "membuka halaman "+meta.Path); err != nil {
			log.Printf("gagal mencatat activity log: %v", err)
		}
	}
}

func moduleFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/api/v1"), "/")
	if trimmed == "" {
		return "unknown"
	}

	segments := strings.Split(trimmed, "/")
	module := "unknown"
	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			continue
		}
		module = segment
	}
	return module
}
