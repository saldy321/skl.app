package controllers

import (
	"github.com/gofiber/fiber/v2"
)



func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logout berhasil!",
	})
}
