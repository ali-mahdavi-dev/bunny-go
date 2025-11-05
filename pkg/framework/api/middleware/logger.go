package middleware

import (
	"time"

	"shikposh-backend/pkg/framework/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) DefaultStructuredLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

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

		entry := logging.Info("HTTP Request").
			WithAny("path", path).
			WithAny("client_ip", clientIP).
			WithAny("method", method).
			WithAny("latency", latency).
			WithAny("status_code", status).
			WithAny("body_size", len(respBody))

		if err != nil {
			entry = entry.WithError(err)
		}

		entry.Log()

		return err
	}
}
