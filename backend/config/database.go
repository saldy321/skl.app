package config

import (
	"fmt"
	"log"
	"os"
	"skl-bakcend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB	

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	dsn := os.Getenv("DB_URL")

if !strings.Contains(dsn, "parseTime=True") {
        if strings.Contains(dsn, "?") {
            dsn += "&parseTime=True&loc=Local"
        } else {
            dsn += "?parseTime=True&loc=Local"
        }
    }


	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal koneksi databse",err )
	}
	fmt.Println("database terhubung")

	err = DB.AutoMigrate(
		&models.Admin{}, 
		&models.SuperAdmin{},
		&models.Instansi{},
		&models.TemplateSKL{},
		&models.Siswa{},
		&models.Mapel{},
		&models.Nilai{},
		&models.BackupInstansi{},
		&models.BackupSetting{},
		&models.Setting_pesan{},
		 &models.SettingBackground{},)
    if err != nil {
        fmt.Println(" Gagal Migrate:", err)
    } else {
        fmt.Println(" Tabel berhasil di-migrate!")
    }
}