package super_admin

import (
	"encoding/json"
	"fmt"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"log"
	"os"
	"path/filepath"
	
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ==================== BACKUP & RESTORE DATA INSTANSI ====================

// BackupManual - Super Admin klik tombol "Backup Sekarang"
func BackupManual(c *fiber.Ctx) error {
	// 1. Cek role: hanya Super Admin
	role := c.Locals(utils.KeyRole).(string)
	if role != "super_admin" {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
	}

	// 2. Ambil ID instansi dari URL
	instansiID := c.Params("id")
	parsedID, err := uuid.Parse(instansiID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format ID Instansi tidak valid"})
	}

	// 3. Jalankan mesin backup
	backup, err := JalankanBackup(parsedID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat backup: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Backup berhasil dibuat!",
		"data":    backup,
	})
}

// RiwayatBackup - Lihat riwayat backup 60 hari terakhir
func RiwayatBackup(c *fiber.Ctx) error {
	// 1. Cek role
	role := c.Locals(utils.KeyRole).(string)
	if role != "super_admin" {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
	}

	// 2. Ambil ID instansi
	instansiID := c.Params("id")

	// 3. Cari backup 60 hari terakhir
	batasWaktu := time.Now().Add(-60 * 24 * time.Hour)

	var backups []models.BackupInstansi
	config.DB.Where("instansi_id = ? AND created_at >= ?", instansiID, batasWaktu).
		Order("created_at DESC").
		Find(&backups)

	// 4. Kelompokkan per tanggal (biar rapi)
	hasil := make(map[string][]models.BackupInstansi)
	for _, b := range backups {
		tanggal := b.CreatedAt.Format("02 January 2006")
		hasil[tanggal] = append(hasil[tanggal], b)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   hasil,
	})
}

// PulihkanBackup - Kembalikan data instansi ke snapshot tertentu
// PulihkanBackup - Kembalikan data instansi ke snapshot tertentu
func PulihkanBackup(c *fiber.Ctx) error {
	// 1. Cek role
	role := c.Locals(utils.KeyRole).(string)
	if role != "super_admin" {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
	}

	// 2. Ambil ID backup
	backupID := c.Params("backupId")

	// 3. Cari backup di database
	var backup models.BackupInstansi
	if err := config.DB.First(&backup, "id = ?", backupID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Backup tidak ditemukan"})
	}

	// 4. Parse JSON backup
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(backup.DataJSON), &data); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Format backup rusak, tidak bisa dipulihkan"})
	}

	instansiID := backup.InstansiID.String()

	// 5. Mulai transaksi
	tx := config.DB.Begin()

	// 6. HAPUS DATA SAAT INI (URUTAN PENTING: nilai & mapel dulu, baru siswa)
	if err := tx.Unscoped().Where("instansi_id = ?", instansiID).Delete(&models.Nilai{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus nilai: " + err.Error()})
	}
	if err := tx.Unscoped().Where("instansi_id = ?", instansiID).Delete(&models.Mapel{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus mapel: " + err.Error()})
	}
	if err := tx.Unscoped().Where("instansi_id = ?", instansiID).Delete(&models.Siswa{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus siswa: " + err.Error()})
	}

	if siswaData, ok := data["siswa"]; ok {
		jsonBytes, _ := json.Marshal(siswaData)
		var siswaList []models.Siswa
		if err := json.Unmarshal(jsonBytes, &siswaList); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal parse data siswa"})
		}
		for _, s := range siswaList {
			s.InstansiID = backup.InstansiID
			if err := tx.Create(&s).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"message": "Gagal restore siswa: " + err.Error()})
			}
		}
	}

	if mapelData, ok := data["mapel"]; ok {
		jsonBytes, _ := json.Marshal(mapelData)
		var mapelList []models.Mapel
		if err := json.Unmarshal(jsonBytes, &mapelList); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal parse data mapel"})
		}
		for _, m := range mapelList {
			m.InstansiID = backup.InstansiID
			if err := tx.Create(&m).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"message": "Gagal restore mapel: " + err.Error()})
			}
		}
	}

	
	if nilaiData, ok := data["nilai"]; ok {
		jsonBytes, _ := json.Marshal(nilaiData)
		var nilaiList []models.Nilai
		if err := json.Unmarshal(jsonBytes, &nilaiList); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"message": "Gagal parse data nilai"})
		}
		for _, n := range nilaiList {
			n.InstansiID = backup.InstansiID
			if err := tx.Create(&n).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"message": "Gagal restore nilai: " + err.Error()})
			}
		}
	}

	// 10. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal commit transaksi: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Data berhasil dikembalikan ke backup: %s", backup.NamaBackup),
	})
}
// 🔧 MESIN BACKUP - Fungsi internal (tidak diekspos ke API)
func JalankanBackup(instansiID uuid.UUID) (*models.BackupInstansi, error) {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

    log.Printf("[BACKUP] Memulai backup untuk instansi ID: %s", instansiID)
    
    var siswa []models.Siswa
    var nilai []models.Nilai
    var mapel []models.Mapel

    // Ambil data
    config.DB.Unscoped().Where("instansi_id = ?", instansiID).Find(&siswa)
    config.DB.Unscoped().Where("instansi_id = ?", instansiID).Find(&nilai)
    config.DB.Unscoped().Where("instansi_id = ?", instansiID).Find(&mapel)

    // Bungkus jadi paket
    paket := map[string]interface{}{
        "siswa": siswa,
        "nilai": nilai,
        "mapel": mapel,
    }

    bungkusan, err := json.Marshal(paket)
    if err != nil {
        return nil, err
    }

    // ========== SIMPAN KE FILE ==========
    backupDir := "./backups"
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return nil, err
    }

    fileName := fmt.Sprintf("backup_%s_%s.json", instansiID.String(), time.Now().Format("2006-01-02_15-04-05"))
    filePath := filepath.Join(backupDir, fileName)

    if err := os.WriteFile(filePath, bungkusan, 0644); err != nil {
        return nil, err
    }

    // ========== SIMPAN KE DATABASE ==========
    backup := &models.BackupInstansi{
        ID:          uuid.New(),
        InstansiID:  instansiID,
        NamaBackup:  fmt.Sprintf("Snapshot - %s", now.Format("02 Jan 2006 15:04 WIB")),
        DataJSON:    string(bungkusan),
        FilePath:    filePath,
        JumlahSiswa: len(siswa),
        JumlahNilai: len(nilai),
        JumlahMapel: len(mapel),
        CreatedAt:   now,
    }

    if err := config.DB.Create(&backup).Error; err != nil {
        return nil, err
    }

    log.Printf("[BACKUP] SUCCESS - Backup ID: %s, File: %s", backup.ID, filePath)
    return backup, nil
}

func DownloadBackup(c *fiber.Ctx) error {
    // 1. Cek role
    role := c.Locals(utils.KeyRole).(string)
    if role != "super_admin" {
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
    }

    // 2. Ambil ID backup
    backupID := c.Params("backupId")

    // 3. Cari backup di database
    var backup models.BackupInstansi
    if err := config.DB.First(&backup, "id = ?", backupID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Backup tidak ditemukan"})
    }

    // 4. Cek apakah file exists
    if _, err := os.Stat(backup.FilePath); os.IsNotExist(err) {
        return c.Status(404).JSON(fiber.Map{"message": "File backup tidak ditemukan"})
    }

    // 5. Kirim file sebagai download
    return c.Download(backup.FilePath, backup.NamaBackup+".json")
}

func ImportBackup(c *fiber.Ctx) error {
    // 1. Cek role
    role := c.Locals(utils.KeyRole).(string)
    if role != "super_admin" {
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
    }

    // 2. Ambil ID instansi
    instansiID := c.Params("id")
    parsedID, err := uuid.Parse(instansiID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format ID Instansi tidak valid"})
    }

    // 3. Parse JSON dari body
    var data map[string]interface{}
    if err := c.BodyParser(&data); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format JSON tidak valid"})
    }

    // 4. Validasi struktur
    if _, ok := data["siswa"]; !ok {
        return c.Status(400).JSON(fiber.Map{"message": "Key 'siswa' tidak ditemukan"})
    }
    if _, ok := data["nilai"]; !ok {
        return c.Status(400).JSON(fiber.Map{"message": "Key 'nilai' tidak ditemukan"})
    }
    if _, ok := data["mapel"]; !ok {
        return c.Status(400).JSON(fiber.Map{"message": "Key 'mapel' tidak ditemukan"})
    }

    // 5. Mulai transaksi
    tx := config.DB.Begin()

    // 6. Hapus data lama
    tx.Unscoped().Where("instansi_id = ?", parsedID).Delete(&models.Siswa{})
    tx.Unscoped().Where("instansi_id = ?", parsedID).Delete(&models.Nilai{})
    tx.Unscoped().Where("instansi_id = ?", parsedID).Delete(&models.Mapel{})

    // 7. Import siswa
    if siswaData, ok := data["siswa"].([]interface{}); ok {
        for _, item := range siswaData {
            s := item.(map[string]interface{})
            siswa := models.Siswa{
                ID:          uuid.New(),
                InstansiID:  parsedID,
                NISN:        s["nisn"].(string),
                NamaSiswa:   s["nama_siswa"].(string),
                TempatLahir: s["tempat_lahir"].(string),
                TanggalLahir: s["tanggal_lahir"].(string),
                Kelas:       s["kelas"].(string),
                Jurusan:     s["jurusan"].(string),
                TahunLulus:  s["tahun_lulus"].(string),
                StatusLulus: s["status_lulus"].(bool),
            }
            tx.Create(&siswa)
        }
    }

    // 8. Import mapel
  // 8. Import mapel
if mapelData, ok := data["mapel"].([]interface{}); ok {
    for _, item := range mapelData {
        m := item.(map[string]interface{})
        mapel := models.Mapel{
            InstansiID: parsedID,
            NamaMapel:  m["nama_mapel"].(string),
        }
        tx.Create(&mapel)
    }
}

    // 9. Import nilai
  // 9. Import nilai
if nilaiData, ok := data["nilai"].([]interface{}); ok {
    for _, item := range nilaiData {
        n := item.(map[string]interface{})
        nilai := models.Nilai{
            InstansiID: parsedID,
            SiswaID:    uuid.MustParse(n["siswa_id"].(string)),
            MapelID:    uint(n["mapel_id"].(float64)),  // ← konversi ke uint
            NilaiAngka: n["nilai_angka"].(float64),
            TahunAjaran: n["tahun_ajaran"].(string),
            Semester:   int(n["semester"].(float64)),
        }
        tx.Create(&nilai)
    }
}

    // 10. Commit transaksi
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return c.Status(500).JSON(fiber.Map{"message": "Gagal commit: " + err.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Data berhasil di-import",
    })
}