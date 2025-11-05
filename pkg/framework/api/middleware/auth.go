package middleware

import (
	"net/http"
	"strings"

	httpapi "shikposh-backend/pkg/framework/api/http"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		tokenStr := parts[1]

		// Parse JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.Cfg.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Extract claims and validate token from DB
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := cast.ToUint64(claims["user_id"])
			ctx := c.Context()
			user, err := m.Uow.Token(ctx).FindByUserID(ctx, userID)
			if err != nil {
				return httpapi.ResError(c, errFailGetTokenFromDB)
			}
			if user.Token != tokenStr {
				return httpapi.ResError(c, errTokenDoesNotExist)
			}

			// Store user_id in Fiber context
			c.Locals("user_id", userID)
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		return c.Next()
	}
}
