package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestContentTypesMiddleware(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	for _, ti := range []struct {
		msg                  string
		reqContentType       string
		expectedContentTypes []string
		statusCode           int
	}{
		{
			msg:                  "valid content-type",
			reqContentType:       "application/json",
			expectedContentTypes: []string{"application/json"},
			statusCode:           http.StatusOK,
		},
		{
			msg:                  "invalid content-type",
			reqContentType:       "application/not+json",
			expectedContentTypes: []string{"application/json"},
			statusCode:           http.StatusUnsupportedMediaType,
		},
		{
			msg:                  "valid multiple content-types",
			reqContentType:       "application/json",
			expectedContentTypes: []string{"application/json", "application/xml"},
			statusCode:           http.StatusOK,
		},
	} {
		t.Run(ti.msg, func(t *testing.T) {
			router := gin.New()
			router.Use(ContentTypes(ti.expectedContentTypes...))
			router.GET("/content-type", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, err := http.NewRequest("GET", "/content-type", nil)
			if err != nil {
				t.Errorf("should not fail: %s", err)
			}
			req.Header.Set("Content-Type", ti.reqContentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != ti.statusCode {
				t.Errorf("expected response code %d, got %d", ti.statusCode, w.Code)
			}
		})
	}
}
