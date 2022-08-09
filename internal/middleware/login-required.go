package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func NewLoginRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		tokenJWT := c.Cookies("Bearer")

		tokenClaims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenJWT, tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			err = fiber.NewError(fiber.StatusUnauthorized, "wrong auth token")
			return
		}

		c.Locals("tokenClaims", tokenClaims)
		return c.Next()
	}
}
