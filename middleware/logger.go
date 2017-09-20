package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LogrusLogger instance a Logger middleware that uses logrus for logging.
// The output is in structured logging with the following format:
// time="2017-09-19T15:27:37+02:00" level=info client="::1" duration=66ns method=GET path=/ status=404
func LogrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		end := time.Now()
		duration := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.WithFields(log.Fields{
			"status":   status,
			"duration": duration,
			"client":   clientIP,
			"method":   method,
			"path":     path,
		}).Infoln(comment)
	}
}
