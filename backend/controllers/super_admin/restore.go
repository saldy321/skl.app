package super_admin

import (
	"skl-bakcend/controllers" // Import untuk SetAuthCookie
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RestoreSuperAdmin - Kembalikan session ke Super Admin (untuk exit impersonate)
func RestoreSuperAdmin(c *fiber.Ctx) error {
	// Generate token baru dengan role super admin
	fakeUUID := uuid.New()
	token, err := utils.GenerateToken(
		fakeUUID.String(),
		models.RoleSuperAdmin,
		"Pusat",
		"Super Admin System",
		uuid.Nil,
		"pusat",
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin token"})
	}

	controllers.SetAuthCookie(c, token)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Kembali ke mode Super Admin",
		"role":    models.RoleSuperAdmin,
	})
}