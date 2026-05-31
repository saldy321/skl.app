package admin

import (
	"fmt"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AmbilTemplateSKL - Mengambil template SKL berdasarkan InstansiID (untuk Admin)
func AmbilTemplateSKL(c *fiber.Ctx) error {
	instansiIDStr := c.Query("instansi_id")

	if instansiIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Instansi ID wajib ada!",
		})
	}

	parsedUUID, err := uuid.Parse(instansiIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Format UUID Instansi tidak valid.",
		})
	}

	var template models.TemplateSKL

	result := config.DB.Where("instansi_id = ?", parsedUUID).First(&template)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(fiber.Map{"data": nil})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data template."})
	}

	return c.JSON(fiber.Map{"data": template})
}

// SimpanTemplateSKL - Menyimpan atau Update template SKL (untuk Admin)
func SimpanTemplateSKL(c *fiber.Ctx) error {
	instansiIDRaw := c.Locals(utils.KeyInstansiID)
	if instansiIDRaw == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi habis, login dulu!"})
	}
	instansiID := instansiIDRaw.(uuid.UUID)

	var input models.TemplateSKL

	if err := c.BodyParser(&input); err != nil {
		fmt.Println("ERROR BODY PARSER:", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"message": "Format data JSON tidak valid.",
			"detail":  err.Error(),
		})
	}

	input.InstansiID = instansiID

	var existing models.TemplateSKL
	err := config.DB.Where("instansi_id = ?", instansiID).First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if createErr := config.DB.Create(&input).Error; createErr != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat template baru: " + createErr.Error()})
			}
			return c.JSON(fiber.Map{"status": "success", "message": "Template baru berhasil dibuat!", "data": input})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Error saat mengecek data lama."})
	}

	if updateErr := config.DB.Model(&existing).Select("*").Omit("id", "created_at", "updated_at").Updates(input).Error; updateErr != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update template! " + updateErr.Error()})
	}

	var updatedTemplate models.TemplateSKL
	config.DB.Where("instansi_id = ?", instansiID).First(&updatedTemplate)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Template sukses di-update!",
		"data":    updatedTemplate,
	})
}