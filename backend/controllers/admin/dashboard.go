package admin

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAdminDashboard(c *fiber.Ctx) error {
	val := c.Locals(utils.KeyInstansiID)
	instansiID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi instansi tidak valid, login ulang bro"})
	}

	var totalSiswa int64
	var siswaLulus int64
	var siswaTidakLulus int64
	var totalMapel int64
	var totalNilai int64
	var belumNilai int64

	// 1. Total Siswa
	config.DB.Model(&models.Siswa{}).Where("instansi_id = ?", instansiID).Count(&totalSiswa)

	// 2. Siswa Lulus
	config.DB.Model(&models.Siswa{}).Where("instansi_id = ? AND status_lulus = ?", instansiID, true).Count(&siswaLulus)

	// 3. Siswa Tidak Lulus
	config.DB.Model(&models.Siswa{}).Where("instansi_id = ? AND status_lulus = ?", instansiID, false).Count(&siswaTidakLulus)

	// 4. Total Mapel
	config.DB.Model(&models.Mapel{}).Where("instansi_id = ?", instansiID).Count(&totalMapel)

	// 5. Total Nilai
	config.DB.Model(&models.Nilai{}).Where("instansi_id = ?", instansiID).Count(&totalNilai)

	// 6. Belum Nilai (opsional, bisa dihitung lebih lanjut)
	belumNilai = 0 // default, bisa dihitung nanti

	return c.JSON(fiber.Map{
		"totalSiswa":      totalSiswa,
		"siswaLulus":      siswaLulus,
		"siswaTidakLulus": siswaTidakLulus,
		"totalMapel":      totalMapel,
		"totalNilai":      totalNilai,
		"belumNilai":      belumNilai,
	})
}