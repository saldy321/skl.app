package controllers

import (
	"github.com/gofiber/fiber/v2"
)


func SetAuthCookie(c *fiber.Ctx, token string) {

	c.ClearCookie("token")
	
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = false
	cookie.SameSite = "None"
	cookie.MaxAge = 86400 
	c.Cookie(cookie)
}