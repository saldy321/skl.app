package admin

import (
    "skl-bakcend/config"
    "skl-bakcend/models"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "path/filepath"
    "os"
	  
)

type SettingBackgroundRequest struct {
    Background string `json:"background"`
}

// GET background setting (untuk siswa & admin)
func GetBackgroundSettings(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    var instansi models.Instansi
    if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    var setting models.SettingBackground
    err := config.DB.Where("instansi_id = ?", instansi.ID).First(&setting).Error
    
    if err == gorm.ErrRecordNotFound {
        return c.Status(200).JSON(fiber.Map{
            "success": true,
            "data": SettingBackgroundRequest{
                Background: "",
            },
        })
    }
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data",
        })
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "data": SettingBackgroundRequest{
            Background: setting.Background,
        },
    })
}

// UPLOAD background (admin only)
func UploadBackground(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    // Ambil file dari form
    file, err := c.FormFile("background")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "File tidak ditemukan",
        })
    }
    
    // Validasi tipe file
    ext := filepath.Ext(file.Filename)
    allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
    if !allowedExt[ext] {
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Format file harus JPG, JPEG, PNG, atau WEBP",
        })
    }
    
    // Cari instansi
    var instansi models.Instansi
    if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    // Buat folder jika belum ada
    uploadDir := "./public/uploads/backgrounds"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gagal membuat folder upload",
        })
    }
    
    // Generate nama file unik
    filename := slug + "_bg" + ext
    filepath := uploadDir + "/" + filename
    
    // Simpan file
    if err := c.SaveFile(file, filepath); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gagal menyimpan file",
        })
    }

      /*img, err := imaging.Open(filepath)
    if err == nil {
        // Resize ke lebar 1920px (tinggi otomatis proporsional)
        resized := imaging.Resize(img, 1920, 0, imaging.Lanczos)
        
        // Simpan ulang dengan kualitas 85%
        out, err := os.Create(filepath)
        if err == nil {
            defer out.Close()
            // Konversi ke JPEG biar konsisten
            jpeg.Encode(out, resized, &jpeg.Options{Quality: 85})
        }
    }
		*/
    
    // Simpan ke database
    var setting models.SettingBackground
    err = config.DB.Where("instansi_id = ?", instansi.ID).First(&setting).Error
    
    if err == gorm.ErrRecordNotFound {
        setting = models.SettingBackground{
            InstansiID: instansi.ID,
            Background: filename,
        }
        if err := config.DB.Create(&setting).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menyimpan ke database",
            })
        }
    } else {
        setting.Background = filename
        if err := config.DB.Save(&setting).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{
                "success": false,
                "message": "Gagal update database",
            })
        }
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "message": "Background berhasil diupload",
        "data": SettingBackgroundRequest{
            Background: filename,
        },
    })
}

// DELETE background (admin only)
func DeleteBackground(c *fiber.Ctx) error {
    slug := c.Params("slug")
    
    var instansi models.Instansi
    if err := config.DB.Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Instansi tidak ditemukan",
        })
    }
    
    var setting models.SettingBackground
    if err := config.DB.Where("instansi_id = ?", instansi.ID).First(&setting).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "success": false,
            "message": "Setting background tidak ditemukan",
        })
    }
    
    // Hapus file
    if setting.Background != "" {
        filePath := "./public/uploads/backgrounds/" + setting.Background
        os.Remove(filePath)
    }
    
    // Hapus dari database
    if err := config.DB.Delete(&setting).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gagal menghapus data",
        })
    }
    
    return c.Status(200).JSON(fiber.Map{
        "success": true,
        "message": "Background berhasil dihapus",
    })
}