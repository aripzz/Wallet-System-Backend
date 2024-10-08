package middleware

import (
	"Wallet-System-Backend/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type UserClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak ditemukan"})
		}

		// Remove "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(utils.GetJWTSecret()), nil
		})

		// Check for errors and validity
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
		}

		// Extract claims
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "could not parse claims"})
		}

		// Store user ID in the context for later use
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
