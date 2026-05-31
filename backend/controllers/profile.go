package controllers

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMyProfile(c *fiber.Ctx) error {
	userID, okID := c.Locals(utils.KeyUserID).(string)
	role, okRole := c.Locals(utils.KeyRole).(string)

	if !okID || !okRole {
		return c.Status(401).JSON(fiber.Map{"message": "Session abis, login lagi bro"})
	}

	// SUPER ADMIN (Hardcoded)
	if role == models.RoleSuperAdmin {
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":             userID,
				"email":          "superadmin@gmail.com",
				"role":           role,
				"nama_instansi":  "Super Admin System",
                "tingkat":         "pusat",
				"slug":           "pusat",
				"instansi_id":    uuid.Nil.String(),
                
			},
		})
	}

	// ADMIN SEKOLAH
	if role == models.RoleAdmin {
		var admin models.Admin
		if err := config.DB.Preload("Instansi").First(&admin, "id = ?", userID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "Admin gak ada"})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":             admin.ID.String(),
				"email":          admin.Email,
				"role":           models.RoleAdmin,
				"nama_instansi":  admin.Instansi.NamaInstansi,
                "tingkat":        admin.Instansi.TingkatSekolah,
				"slug":           admin.Instansi.Slug,
				"instansi_id":    admin.InstansiID.String(),
                "foto_profile":   admin.FotoProfile,
                 "logo_instansi":  admin.Instansi.LogoInstansi,
			},
		})
	}

	return c.Status(403).JSON(fiber.Map{"message": "Role gak dikenali"})
}