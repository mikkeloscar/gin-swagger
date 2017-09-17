package restapi

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LogrusLogger instance a Logger
// Example: os.Stdout, a file opened in write mode, a socket...
func LogrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		//	defaultLogFormat = "timestamp=%v status_code=%d response_time=%v client=%s http_method=%s path=%s msg=%s\n"
		log.WithFields(log.Fields{
			"timestamp":     end.Format("2006-01-02-15:04:05"),
			"status_code":   statusCode,
			"response_time": latency,
			"client:":       clientIP,
			"method":        method,
			"path":          path,
		}).Infoln(comment)
	}
}
