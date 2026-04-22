package controllers

import (
	"skl-bakcend/models" 
	"skl-bakcend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
    "skl-bakcend/utils"
    "fmt"
)

// GetTemplateSKL - Mengambil template berdasarkan InstansiID
func GetTemplateSKL(c *fiber.Ctx) error {
	
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
			// Return null jika belum ada template, frontend akan pakai default
			return c.JSON(fiber.Map{"data": nil})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data template."})
	}

	return c.JSON(fiber.Map{"data": template})
}

// SaveTemplateSKL - Menyimpan atau Update template
// SaveTemplateSKL - Menyimpan atau Update template
// SaveTemplateSKL - Menyimpan atau Update template
func SaveTemplateSKL(c *fiber.Ctx) error {
   
    instansiIDRaw := c.Locals(utils.KeyInstansiID)
    if instansiIDRaw == nil {
        return c.Status(401).JSON(fiber.Map{"message": "Sesi habis, login dulu!"})
    }
    instansiID := instansiIDRaw.(uuid.UUID)

    var input models.TemplateSKL
    
    // Parse JSON body ke struct Model
    if err := c.BodyParser(&input); err != nil {
        // Tambahkan log error detail untuk debugging
        fmt.Println("ERROR BODY PARSER:", err.Error())
        return c.Status(400).JSON(fiber.Map{
            "message": "Format data JSON tidak valid.",
            "detail": err.Error(),
        })
    }

    // Pastikan InstansiID terisi dengan benar dari session (keamanan)
    input.InstansiID = instansiID

    var existing models.TemplateSKL
 
    // Cek apakah sudah ada template untuk instansi ini
    err := config.DB.Where("instansi_id = ?", instansiID).First(&existing).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            // CREATE BARU
            if createErr := config.DB.Create(&input).Error; createErr != nil {
                return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat template baru: " + createErr.Error()})
            }
            return c.JSON(fiber.Map{"status": "success", "message": "Template baru berhasil dibuat!", "data": input})
        }
        return c.Status(500).JSON(fiber.Map{"message": "Error saat mengecek data lama."})
    }

    // UPDATE YANG SUDAH ADA
    // --- PERBAIKAN KRUSIAL DI SINI ---
    // Gunakan .Select(...) untuk memaksa GORM mengupdate field tertentu bahkan jika nilainya zero value (false/0)
    // Atau, gunakan map[string]interface{} untuk update spesifik.
    // Tapi cara paling aman dan rapi adalah menggunakan Select dengan daftar field yang ingin diupdate.
    
    // Karena kita ingin update SEMUA field yang ada di struct input, kita bisa pakai Select("*") 
    // ATAU lebih baik, sebutkan field-field penting yang rentan di-skip oleh GORM.
    
    // Opsi 1: Update Semua Field (Termasuk Zero Value) dengan Select("*")
    // Ini akan mengupdate semua kolom di tabel sesuai isi struct 'input'
    if updateErr := config.DB.Model(&existing).Select("*").Omit("id", "created_at", "updated_at").Updates(input).Error; updateErr != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal update template! " + updateErr.Error()})
    }

    // KEMBALIKAN DATA TERBARU KE FRONTEND
    // Ini penting agar Frontend langsung dapat state terbaru yang pasti tersimpan di DB
    var updatedTemplate models.TemplateSKL
    config.DB.Where("instansi_id = ?", instansiID).First(&updatedTemplate)

    return c.JSON(fiber.Map{
        "status": "success", 
        "message": "Template sukses di-update!",
        "data": updatedTemplate, // Kirim balik data bersih dari DB
    })
}