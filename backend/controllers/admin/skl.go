package admin

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AmbilDataSKLAdmin - Mengambil data lengkap siswa dan template untuk dicetak oleh Admin
func AmbilDataSKLAdmin(c *fiber.Ctx) error {
	instansiIDRaw := c.Locals(utils.KeyInstansiID)
	if instansiIDRaw == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
	}
	instansiID := instansiIDRaw.(uuid.UUID)

	siswaIDParam := c.Params("id")
	siswaID, err := uuid.Parse(siswaIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID Siswa tidak valid"})
	}

	var siswa models.Siswa
	if err := config.DB.Preload("Nilai.Mapel").Where("id = ? AND instansi_id = ?", siswaID, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data siswa tidak ditemukan atau bukan milik sekolah ini"})
	}

	var template models.TemplateSKL
	if err := config.DB.Where("instansi_id = ?", instansiID).First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Template SKL belum dikonfigurasi untuk sekolah ini"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"siswa":    siswa,
			"template": template,
		},
	})
}