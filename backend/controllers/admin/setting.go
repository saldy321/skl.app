package admin

import (
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

  func SetWaktuBukaPengumuman(c *fiber.Ctx) error {
        instansiIDRaw := c.Locals(utils.KeyInstansiID)
        if instansiIDRaw == nil {
            return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
        }
        instansiID := instansiIDRaw.(uuid.UUID)

        var input struct {
            WaktuBuka string `json:"waktu_buka"`
        }

        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
        }

        loc, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal load timezone server"})
        }

        waktu, err := time.ParseInLocation("2006-01-02 15:04:05", input.WaktuBuka, loc)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format waktu harus YYYY-MM-DD HH:MM:SS"})
        }

        if err := config.DB.Model(&models.Instansi{}).Where("id = ?", instansiID).Update("waktu_buka_pengumuman", waktu).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal update jadwal"})
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Jadwal pengumuman berhasil diatur!",
        })
    }

	func UpdateTampilkanLogo(c *fiber.Ctx) error {
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    var input struct {
        TampilkanLogo bool `json:"tampilkan_logo"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
    }

    result := config.DB.Model(&models.Instansi{}).
        Where("id = ?", instansiID).
        Update("tampilkan_logo", input.TampilkanLogo)

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal update setting"})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Setting logo berhasil diupdate",
        "data": fiber.Map{
            "tampilkan_logo": input.TampilkanLogo,
        },
    })
}