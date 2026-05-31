package auth

import (
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"skl-bakcend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"skl-bakcend/controllers"

)

func LoginSiswa(c *fiber.Ctx) error {
log.Println("[LoginSiswa] ========== HANDLER DIPANGGIL ==========")
	instansiID, ok := c.Locals("tenant_id").(uuid.UUID)
	if !ok {
		
		slug := c.Params("slug")
		
		var instansi models.Instansi
		if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Sekolah tidak ditemukan",
			})
		}
		instansiID = instansi.ID
	}

	type LoginRequest struct {
		NISN         string `json:"nisn"`
		TanggalLahir string `json:"tanggal_lahir"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Format data login salah",
		})
	}

	if req.NISN == "" || req.TanggalLahir == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "NISN dan Tanggal Lahir wajib diisi",
		})
	}

	var siswa models.Siswa
	err := config.DB.Where("nisn = ? AND tanggal_lahir = ? AND instansi_id = ?",
		req.NISN, req.TanggalLahir, instansiID).First(&siswa).Error

	if err != nil {
		log.Println("[LOGIN SISWA] ERROR: Siswa tidak ditemukan")
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "NISN atau Tanggal Lahir tidak cocok",
		})
	}

	log.Println("[LOGIN SISWA] Siswa ditemukan:", siswa.NamaSiswa)

	var instansi models.Instansi
	if err := config.DB.Where("id = ?", instansiID).First(&instansi).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Data sekolah tidak ditemukan",
		})
	}

	token, err := utils.GenerateToken(
		siswa.ID.String(),
		"siswa",
		instansi.TingkatSekolah,
		instansi.NamaInstansi,
		instansi.ID,
		instansi.Slug,
	)
	if err != nil {
		log.Println("[LOGIN SISWA]) ERROR generate token:",err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat sesi login",
		})
	}

	controllers.SetAuthCookie(c, token)
	log.Println("[LOGIN SISWA] SUCCESS! Cookie sudah di-set")

	return c.JSON(fiber.Map{
		"status":        "success",
		"success":       true,
		"message":       "Login Berhasil! Halo " + siswa.NamaSiswa,
		"role":          "siswa",
		"slug":          instansi.Slug,
		"nama_instansi": instansi.NamaInstansi,
		"tingkat":		instansi.TingkatSekolah,
		"instansi_id":   instansi.ID.String(),
		"user": fiber.Map{
			"id":    siswa.ID,
			"nisn":  siswa.NISN,
			"nama":  siswa.NamaSiswa,
			"kelas": siswa.Kelas,
		},
	})
}