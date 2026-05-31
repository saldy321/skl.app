package super_admin

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func GetAllAdminSekolah(c *fiber.Ctx) error {
	var admins []models.Admin
	if err := config.DB.Preload("Instansi").Find(&admins).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal tarik data admin"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   admins,
	})
}

// RegisterAdminSekolah - Daftar admin baru (untuk Super Admin)
func RegisterAdminSekolah(c *fiber.Ctx) error {
	var input struct {
		Email      string    `json:"email"`
		Password   string    `json:"password"`
		InstansiID uuid.UUID `json:"instansi_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input data ngaco, cek lagi bro"})
	}

	if input.Email == "" || input.Password == "" || input.InstansiID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email, password, dan instansi wajib diisi!"})
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	newAdmin := models.Admin{
		Email:      input.Email,
		Password:   hashedPassword,
		InstansiID: input.InstansiID,
	}

	if err := config.DB.Create(&newAdmin).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal daftar admin: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Admin sekolah berhasil ditambahkan!",
		"data":    newAdmin,
	})
}

// UpdateAdminEmail - Edit email admin (untuk Super Admin)
func UpdateAdminEmail(c *fiber.Ctx) error {
	id := c.Params("id")

	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	if input.Email == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Email tidak boleh kosong"})
	}

	var admin models.Admin
	if err := config.DB.First(&admin, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Admin tidak ditemukan"})
	}

	if err := config.DB.Model(&admin).Update("email", input.Email).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update email: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Email admin berhasil diupdate",
	})
}