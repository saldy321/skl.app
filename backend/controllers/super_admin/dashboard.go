package super_admin

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
)


func GetSuperDashboard(c *fiber.Ctx) error {
	role := c.Locals(utils.KeyRole).(string)
	if role != models.RoleSuperAdmin {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak, khusus Super Admin!"})
	}

	var totalSekolah, totalAdmin, totalSiswaNasional int64

	config.DB.Model(&models.Instansi{}).Count(&totalSekolah)
	config.DB.Model(&models.Admin{}).Count(&totalAdmin)
	config.DB.Model(&models.Siswa{}).Count(&totalSiswaNasional)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"totalSekolah":       totalSekolah,
			"totalAdmin":         totalAdmin,
			"totalSiswaNasional": totalSiswaNasional,
		},
	})
}