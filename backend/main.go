package main

import (
	"log"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/controllers/super_admin"
	"skl-bakcend/models"
	"skl-bakcend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cookie",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Set-Cookie",
	}))

	app.Static("/uploads", "./public/uploads")

	routes.Setup(app)

	
	go jalankanAutoBackup()

	app.Listen(":3000")
}
func jalankanAutoBackup() {
    log.Println("[AUTO BACKUP] Scheduler started, checking every 20 seconds")
    
    // Track backup terakhir (untuk cek duplikasi)
    var lastBackupTimestamp time.Time
    
    for {
        // STEP 1: Load setting
        var setting models.BackupSetting
        if err := config.DB.First(&setting).Error; err != nil {
            log.Printf("[AUTO BACKUP] Gagal load setting: %v, akan coba lagi 20 detik", err)
            time.Sleep(20 * time.Second)
            continue
        }

        // STEP 2: Kalau auto backup mati, skip
        if !setting.AutoBackup {
            time.Sleep(20 * time.Second)
            continue
        }

        now := time.Now().In(time.FixedZone("WIB", 7*3600)) // WIB
        jamTarget := setting.JamBackup
        menitTarget := setting.MenitBackup

        // STEP 4: Cek apakah sekarang waktu backup (TEPAT, tanpa toleransi)
        isBackupTime := (now.Hour() == jamTarget && now.Minute() == menitTarget)

        // STEP 5: Kalau sudah waktunya dan belum dijalankan hari ini
        if isBackupTime && now.Sub(lastBackupTimestamp) > 1*time.Minute {
            log.Printf("[AUTO BACKUP] ========== MULAI BACKUP OTOMATIS ==========")
            log.Printf("[AUTO BACKUP] Waktu: %s | Target: %02d:%02d", now.Format("15:04:05"), jamTarget, menitTarget)

            // Ambil semua instansi
            var instansiList []models.Instansi
            if err := config.DB.Find(&instansiList).Error; err != nil {
                log.Printf("[AUTO BACKUP] Gagal ambil daftar instansi: %v", err)
                time.Sleep(20 * time.Second)
                continue
            }

            successCount := 0
            failCount := 0

            // Backup satu per satu
            for _, instansi := range instansiList {
                log.Printf("[AUTO BACKUP] Memproses backup untuk: %s", instansi.NamaInstansi)
                _, err := super_admin.JalankanBackup(instansi.ID)
                if err != nil {
                    log.Printf("[AUTO BACKUP]  GAGAL backup %s: %v", instansi.NamaInstansi, err)
                    failCount++
                } else {
                    log.Printf("[AUTO BACKUP]  BERHASIL backup %s", instansi.NamaInstansi)
                    successCount++
                }
            }

            // Hapus backup lama (hard delete)
            batas := time.Now().Add(-time.Duration(setting.RetensiHari) * 24 * time.Hour)
            result := config.DB.Unscoped().Where("created_at < ?", batas).Delete(&models.BackupInstansi{})
            if result.Error != nil {
                log.Printf("[AUTO BACKUP] Gagal hapus backup lama: %v", result.Error)
            } else if result.RowsAffected > 0 {
                log.Printf("[AUTO BACKUP] Menghapus %d backup yang lebih dari %d hari", result.RowsAffected, setting.RetensiHari)
            }

            log.Printf("[AUTO BACKUP] ========== SELESAI ==========")
            log.Printf("[AUTO BACKUP] Backup: %d berhasil, %d gagal", successCount, failCount)
            
            // Update timestamp backup terakhir
            lastBackupTimestamp = now
        }

        // STEP 6: Tunggu 20 detik sebelum cek lagi
        time.Sleep(20 * time.Second)
    }
}