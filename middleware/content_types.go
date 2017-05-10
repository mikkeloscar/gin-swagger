package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikkeloscar/gin-swagger/api"
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
			problem := api.Problem{
				Title:  "Unsupported media type.",
				Status: http.StatusUnsupportedMediaType,
				Detail: fmt.Sprintf("unsupported media type '%s', only %s are allowed",
					reqContentType,
					contentTypes,
				),
			}
			c.Writer.Header.Set("Content-Type", "application/problem+json")
			c.JSON(problem.Status, problem)
			c.Abort()
			return
		}
		c.Next()
	}
}
