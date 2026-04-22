package middleware

import (
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"skl-bakcend/utils"
)

func AuthSiswa(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")

	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "akses dilarang, token tidak ditemukan",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(401, "Method signing tidak valid")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"message": "Sesi lu udah abis atau token palsu bro",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["role"] != "siswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Ini area khusus siswa, dilarang masuk!",
		})
	}

	c.Locals(utils.KeyUserID, claims["id"])
	c.Locals(utils.KeyRole, claims["role"])
	c.Locals(utils.KeyTingkatSekolah, claims["tingkat"])
	c.Locals(utils.KeyInstansi, claims["instansi"])

	if instIDStr, ok := claims["instansi_id"].(string); ok {
		c.Locals(utils.KeyInstansiID, utils.ParseUUID(instIDStr))
	}

	return c.Next()
}