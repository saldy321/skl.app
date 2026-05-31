package admin

import (
    "skl-bakcend/config"
    "skl-bakcend/models"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

type SettingPesanRequest struct {
    PesanLulus      string `json:"pesan_lulus"`
    PesanTidakLulus string `json:"pesan_tidak_lulus"`
}

// Fungsi internal buat ambil data setting
func getSettingBySlug(slug string) (*SettingPesanRequest, error) {
    var instansi models.Instansi
    if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return nil, err
    }
    
    var setting models.Setting_pesan
    err := config.DB.Where("instansi_id = ?", instansi.ID).First(&setting).Error
    
    if err == gorm.ErrRecordNotFound {
        return &SettingPesanRequest{
            PesanLulus:      "",
            PesanTidakLulus: "",
        }, nil
    }
    
    if err != nil {
        return nil, err
    }
    
    return &SettingPesanRequest{
        PesanLulus:      setting.PesanLulus,
        PesanTidakLulus: setting.PesanTidakLulus,
    }, nil
}

// Untuk siswa (public)
func GetSettingsPublic(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    data, err := getSettingBySlug(slug)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}

// Untuk admin (GET)
func GetSettings(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    data, err := getSettingBySlug(slug)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}

// PUT /:slug/admin/setting/pesan
func UpdateSettings(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    var req SettingPesanRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Invalid request body: " + err.Error(),
        })
    }
    
    var instansi models.Instansi
    if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    var setting models.Setting_pesan
    err := config.DB.Where("instansi_id = ?", instansi.ID).First(&setting).Error
    
    if err == gorm.ErrRecordNotFound {
        setting = models.Setting_pesan{
            InstansiID:      instansi.ID,
            PesanLulus:      req.PesanLulus,
            PesanTidakLulus: req.PesanTidakLulus,
        }
        if err := config.DB.Create(&setting).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menyimpan data: " + err.Error(),
            })
        }
    } else if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gagal memproses data: " + err.Error(),
        })
    } else {
        setting.PesanLulus = req.PesanLulus
        setting.PesanTidakLulus = req.PesanTidakLulus
        if err := config.DB.Save(&setting).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "message": "Gagal update data: " + err.Error(),
            })
        }
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "message": "Pengaturan berhasil disimpan",
        "data":    req,
    })
}