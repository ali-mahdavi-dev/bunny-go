package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/service_layer/unit_of_work"
	"github.com/ali-mahdavi-dev/bunny-go/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

var errFailGetTokenFromDB = errors.New("fail to get token from DB")
var errTokenDoesNotExist = errors.New("token does nit exist")

func AuthMiddleware(cfg *config.Config, uow unit_of_work.PGUnitOfWork) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// validate alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return cfg.JWT.Secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store claims in context if needed
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user, err := uow.Token().FindByUserID(c, cast.ToUint64(claims["user_id"]))
			if err != nil {
				ginx.ResError(c, errFailGetTokenFromDB)
				return
			}
			if user.Token != tokenStr {
				ginx.ResError(c, errTokenDoesNotExist)
				return
			}
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
