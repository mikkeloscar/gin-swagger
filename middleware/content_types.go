package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentTypes is a middleware that checks the request Content-Type header and
// responds with 415 if it doesn't match any of the expected Content-Types.
func ContentTypes(contentTypes ...string) gin.HandlerFunc {
	validContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, typ := range contentTypes {
		validContentTypes[typ] = struct{}{}
	}

	return func(c *gin.Context) {
		reqContentType := c.ContentType()
		if _, ok := validContentTypes[reqContentType]; !ok {
			msg := gin.H{
				"code": http.StatusUnsupportedMediaType,
				"message": fmt.Sprintf("unsupported media type '%s', only %s are allowed",
					reqContentType,
					contentTypes,
				),
			}
			c.JSON(http.StatusUnsupportedMediaType, msg)
			c.Abort()
			return
		}
		c.Next()
	}
}
