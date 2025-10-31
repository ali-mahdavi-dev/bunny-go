package middleware

import (
	"strings"
	"time"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) DefaultStructuredLogger() fiber.Handler {
	logger := logging.NewLogger(m.Cfg)
	return func(c *fiber.Ctx) error {
		if strings.Contains(c.Path(), "swagger") {
			return c.Next()
		}

		start := time.Now()

		// Capture request body
		reqBody := c.Body()                                     // []byte
		c.Request().SetBodyRaw(append([]byte(nil), reqBody...)) // restore body

		// Run next handler
		err := c.Next()

		// Capture response body
		respBody := c.Response().Body()

		// Collect metadata
		latency := time.Since(start)
		clientIP := c.IP()
		method := c.Method()
		status := c.Response().StatusCode()
		path := c.OriginalURL()

		keys := map[logging.ExtraKey]interface{}{
			logging.Path:         path,
			logging.ClientIp:     clientIP,
			logging.Method:       method,
			logging.Latency:      latency,
			logging.StatusCode:   status,
			logging.RequestBody:  string(reqBody),
			logging.ResponseBody: string(respBody),
			logging.BodySize:     len(respBody),
			logging.ErrorMessage: "",
		}

		if err != nil {
			keys[logging.ErrorMessage] = err.Error()
		}

		logger.Info(logging.RequestResponse, logging.Api, "", keys)

		return err
	}
}
