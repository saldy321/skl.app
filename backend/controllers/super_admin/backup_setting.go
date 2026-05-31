package super_admin

import (
    "time"
    "skl-bakcend/config"
    "skl-bakcend/models"
    "skl-bakcend/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type SettingInput struct {
    JamBackup   int  `json:"jam_backup"`
    MenitBackup int  `json:"menit_backup"`
    RetensiHari int  `json:"retensi_hari"`
    AutoBackup  bool `json:"auto_backup"`
}

func AmbilSettingBackup(c *fiber.Ctx) error {
    role := c.Locals(utils.KeyRole).(string)
    if role != "super_admin" {
        return c.Status(403).JSON(fiber.Map{"message": "Hanya Super Admin!"})
    }

    var setting models.BackupSetting
    err := config.DB.First(&setting).Error
    if err != nil {
        setting = models.BackupSetting{
            ID:          uuid.New(),
            JamBackup:   2,
            MenitBackup: 0,
            RetensiHari: 60,
            AutoBackup:  true,
        }
        config.DB.Create(&setting)
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   setting,
    })
}

func SimpanSettingBackup(c *fiber.Ctx) error {
    role := c.Locals(utils.KeyRole).(string)
    if role != "super_admin" {
        return c.Status(403).JSON(fiber.Map{"message": "Hanya Super Admin!"})
    }

    var input SettingInput
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
    }

    // Validasi
    if input.JamBackup < 0 || input.JamBackup > 23 {
        return c.Status(400).JSON(fiber.Map{"message": "Jam harus antara 0-23"})
    }
    if input.MenitBackup < 0 || input.MenitBackup > 59 {
        return c.Status(400).JSON(fiber.Map{"message": "Menit harus antara 0-59"})
    }

    var setting models.BackupSetting
    config.DB.First(&setting)

    config.DB.Model(&setting).Updates(map[string]interface{}{
        "jam_backup":    input.JamBackup,
        "menit_backup":  input.MenitBackup,
        "retensi_hari":  input.RetensiHari,
        "auto_backup":   input.AutoBackup,
        "updated_at":    time.Now(),
    })

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Pengaturan backup berhasil disimpan!",
    })
}