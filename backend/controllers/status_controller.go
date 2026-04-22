package controllers

import (
	"fmt"
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ProcessKelulusan: Logic utama buat nentuin siapa yang lulus
func ProcessKelulusan(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	
	// 1. Ambil KKM dari template (Biar Dinamis)
	var conf models.TemplateSKL
	if err := config.DB.Where("instansi_id = ?", instansiID).First(&conf).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Settingan KKM di Template SKL belum ada!"})
	}

	var siswas []models.Siswa
	// Preload Nilai buat dihitung rata-ratanya
	config.DB.Preload("Nilai").Where("instansi_id = ?", instansiID).Find(&siswas)

	for _, s := range siswas {
		if len(s.Nilai) == 0 { continue }
		
		var total float64
		for _, n := range s.Nilai {
			total += n.NilaiAngka
		}

		rataRata := total / float64(len(s.Nilai))
		// Bandingin sama KKM dari tabel TemplateSKL
		isLulus := rataRata >= conf.MinimalKelulusan

		config.DB.Model(&s).Updates(map[string]interface{}{
			"rata_rata_nilai": rataRata,
			"status_lulus":    isLulus,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success", 
		"message": fmt.Sprintf("Kelulusan diproses. KKM: %.2f", conf.MinimalKelulusan),
	})
}