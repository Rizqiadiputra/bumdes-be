package utils

import "github.com/gin-gonic/gin"

type RequestMeta struct {
	Method    string
	Path      string
	IPAddress string
	UserAgent string
}

func GetRequestMeta(c *gin.Context) RequestMeta {
	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	return RequestMeta{
		Method:    c.Request.Method,
		Path:      path,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
}
