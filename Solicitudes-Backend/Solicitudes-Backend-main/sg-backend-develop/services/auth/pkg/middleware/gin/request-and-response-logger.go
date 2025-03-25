package sdkmwr

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
)

type HttpLoggingOptions struct {
	LogLevel       string
	IncludeHeaders bool
	IncludeBody    bool
	ExcludedPaths  []string
}

// INFO: registra y loggea las solicitudes HTTP entrantes y las respuestas salientes
func RequestAndResponseLogger(options HttpLoggingOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path is excluded
		for _, path := range options.ExcludedPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// Registrar la solicitud entrante
		startTime := time.Now()
		logger.Infof("Incoming request: %s %s", c.Request.Method, c.Request.URL)

		if options.IncludeHeaders {
			logger.Infof("Request headers: %v", c.Request.Header)
		}

		if options.IncludeBody {
			// Log the request body if needed (make sure to handle large bodies appropriately)
			// NOTE: You might need to handle the body as a buffer and reset it afterward
		}

		// Procesar la solicitud
		c.Next()

		// Registrar la respuesta saliente
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()
		logger.Infof("Response: %d, Latency: %v", statusCode, latency)
	}
}
