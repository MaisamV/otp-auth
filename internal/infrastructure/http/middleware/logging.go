package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	SkipPaths    []string // Paths to skip logging
	LogBody      bool     // Whether to log request/response body
	MaxBodySize  int      // Maximum body size to log
	SensitiveHeaders []string // Headers to mask in logs
}

// DefaultLoggingConfig returns a default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		SkipPaths: []string{"/health", "/live", "/ready"},
		LogBody:   false,
		MaxBodySize: 1024, // 1KB
		SensitiveHeaders: []string{"Authorization", "Cookie", "Set-Cookie"},
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging returns a logging middleware with the given configuration
func Logging(config LoggingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for specified paths
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body if logging is enabled
		var requestBody string
		if config.LogBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				if len(bodyBytes) <= config.MaxBodySize {
					requestBody = string(bodyBytes)
				} else {
					requestBody = string(bodyBytes[:config.MaxBodySize]) + "...[truncated]"
				}
				// Restore the request body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Wrap response writer to capture response body
		var responseBody string
		if config.LogBody {
			writer := &responseWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = writer
			defer func() {
				if writer.body.Len() <= config.MaxBodySize {
					responseBody = writer.body.String()
				} else {
					responseBody = writer.body.String()[:config.MaxBodySize] + "...[truncated]"
				}
			}()
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		// Build query string
		if raw != "" {
			path = path + "?" + raw
		}

		// Mask sensitive headers
		headers := make(map[string]string)
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				value := values[0]
				for _, sensitiveHeader := range config.SensitiveHeaders {
					if key == sensitiveHeader {
						value = "[MASKED]"
						break
					}
				}
				headers[key] = value
			}
		}

		// Log the request
		logEntry := map[string]interface{}{
			"timestamp":    start.Format(time.RFC3339),
			"client_ip":    clientIP,
			"method":       method,
			"path":         path,
			"status_code":  statusCode,
			"latency":      latency.String(),
			"body_size":    bodySize,
			"user_agent":   userAgent,
			"headers":      headers,
		}

		if config.LogBody {
			logEntry["request_body"] = requestBody
			logEntry["response_body"] = responseBody
		}

		// Log based on status code
		if statusCode >= 500 {
			log.Printf("[ERROR] %+v", logEntry)
		} else if statusCode >= 400 {
			log.Printf("[WARN] %+v", logEntry)
		} else {
			log.Printf("[INFO] %+v", logEntry)
		}
	}
}

// SimpleLogging returns a simple logging middleware
func SimpleLogging() gin.HandlerFunc {
	return Logging(DefaultLoggingConfig())
}

// DetailedLogging returns a detailed logging middleware that logs request/response bodies
func DetailedLogging() gin.HandlerFunc {
	config := DefaultLoggingConfig()
	config.LogBody = true
	return Logging(config)
}