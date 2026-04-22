package controllers

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils" // Import utils lu bro
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSiswaDashboard(c *fiber.Ctx) error {
	// 1. Ambil ID dari Locals pake utils key
	// Pastiin tipe datanya sesuai sama yang di-set di middleware
	siswaID := c.Locals(utils.KeyUserID) 
	
	// Ambil InstansiID sebagai uuid.UUID (karena di middleware udah kita parse)
	instansiID, ok := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Instansi ID kaga kebaca, coba login ulang bro",
		})
	}

	// 2. Tarik data Siswa + Relasi Nilai
	var siswa models.Siswa

	if err := config.DB.Preload("Nilai.Mapel").Where("id = ? AND instansi_id = ?", siswaID, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Data lu kaga ketemu, mungkin beda sekolah atau emang gaada",
		})
	}

	// 3. Tarik data Template SKL
	var template models.TemplateSKL
	if err := config.DB.Where("instansi_id = ?", instansiID).First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Template SKL belom disiapin sama admin sekolah lu",
		})
	}

	// 4. Kirim datanya (Pake struktur yang konsisten buat frontend)
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"siswa":    siswa,
			"template": template,
			"catatan":  "Data aman, silakan cetak SKL!",
		},
	})
}