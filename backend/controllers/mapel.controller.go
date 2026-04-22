package controllers

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllMapel(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	var mapels []models.Mapel
	config.DB.Where("instansi_id = ?", instansiID).Order("id DESC").Find(&mapels)
	return c.JSON(fiber.Map{"status": "success", "data": mapels})
}

func CreateMapel(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	var mapel models.Mapel
	if err := c.BodyParser(&mapel); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input error"})
	}
	mapel.InstansiID = instansiID
	if err := config.DB.Create(&mapel).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan"})
	}
	return c.JSON(fiber.Map{"status": "success", "data": mapel})
}

func UpdateMapel(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	id := c.Params("id")
	
	var mapel models.Mapel
	// Cari dulu datanya
	if err := config.DB.Where("id = ? AND instansi_id = ?", id, instansiID).First(&mapel).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data tidak ditemukan"})
	}

	// Masukin data baru dari body ke variabel mapel
	if err := c.BodyParser(&mapel); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Payload salah"})
	}

	config.DB.Save(&mapel) 
	return c.JSON(fiber.Map{"status": "success", "message": "Berhasil diupdate"})
}

func DeleteMapel(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	id := c.Params("id")
	
	// Pake Unscoped() biar ilang permanen (biar gak menuhin DB kalo salah input)
	result := config.DB.Unscoped().Where("id = ? AND instansi_id = ?", id, instansiID).Delete(&models.Mapel{})
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Gagal hapus, data ga ada"})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Data dihapus"})
}