package main

import (
	"skl-bakcend/config" 
	"skl-bakcend/routes" 
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Cookie",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,	
		ExposeHeaders: "Set-Cookie",
	}))

	// --- INI KUNCINYA ---
	// Pastikan folder "./public/uploads" ADA secara fisik di project kamu
	// Pastikan baris ini DIPANGGIL sebelum routes.Setup jika ada route wildcard /*
	app.Static("/uploads", "./public/uploads")
	// -------------------

	routes.Setup(app)

	app.Listen(":3000")
}