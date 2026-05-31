package admin


import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"     
	"gorm.io/gorm/clause"              
)




func GetAllSiswa(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

	// 1. Ambil parameter tahun_lulus dari Query URL (dikirim oleh Frontend via Store)
	tahunLulus := c.Query("tahun_lulus")

	var siswas []models.Siswa
	
	// 2. Mulai Query Dasar GORM
	query := config.DB.Preload("Nilai").Where("instansi_id = ?", instansiID)

	// 3. Jika ada filter tahun, tambahkan kondisi WHERE
	if tahunLulus != "" {
		query = query.Where("tahun_lulus = ?", tahunLulus)
	}
	
	// 4. Eksekusi Query
	result := query.Find(&siswas)
	
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal narik data siswa"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   siswas,
	})
}

func GetSiswaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format ID Siswa salah bro"})
	}

	var siswa models.Siswa
	if err := config.DB.Preload("Nilai").Where("id = ? AND instansi_id = ?", parsedID, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Siswa gak ketemu!"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   siswa, // Field tahun_lulus otomatis terambil
	})
}

func CreateSiswa(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

	siswa := models.Siswa{
		InstansiID:   instansiID,
		NISN:         c.FormValue("nisn"),
		NamaSiswa:    c.FormValue("nama_siswa"),
		TempatLahir:  c.FormValue("tempat_lahir"),
		TanggalLahir: c.FormValue("tanggal_lahir"),
		JenisKelamin: c.FormValue("jenis_kelamin"),
		Kelas:        c.FormValue("kelas"),
		Jurusan:      c.FormValue("jurusan"),
		NamaWali:     c.FormValue("nama_wali"),
		StatusLulus:  c.FormValue("status_lulus") == "true",
		
		// <--- TAMBAHAN: Ambil Tahun Lulus dari Form
		TahunLulus:   c.FormValue("tahun_lulus"), 
	}

	// Handle Upload Foto
	file, err := c.FormFile("foto_siswa")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", siswa.NISN, time.Now().Unix(), ext)
		targetPath := fmt.Sprintf("./public/uploads/siswa/%s", filename)
		
		if err := c.SaveFile(file, targetPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan foto"})
		}
		siswa.FotoSiswa = filename
	}

	if err := config.DB.Create(&siswa).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Gagal tambah siswa, NISN mungkin udah ada"})
	}

	return c.JSON(fiber.Map{
		"message": "Siswa berhasil ditambah!",
		"data":    siswa,
	})
}

func UpdateSiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

	var siswa models.Siswa
	if err := config.DB.Where("id = ? AND instansi_id = ?", id, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Siswa gak ketemu!"})
	}

	updates := make(map[string]interface{})

	// Manual binding biar data lama gak ketiban kosong
	if val := c.FormValue("nama_siswa"); val != "" { updates["nama_siswa"] = val }
	if val := c.FormValue("nisn"); val != "" { updates["nisn"] = val }
	if val := c.FormValue("tempat_lahir"); val != "" { updates["tempat_lahir"] = val }
	if val := c.FormValue("tanggal_lahir"); val != "" { updates["tanggal_lahir"] = val }
	if val := c.FormValue("kelas"); val != "" { updates["kelas"] = val }
	if val := c.FormValue("jurusan"); val != "" { updates["jurusan"] = val }
	if val := c.FormValue("jenis_kelamin"); val != "" { updates["jenis_kelamin"] = val }
	
	// <--- TAMBAHAN: Update Tahun Lulus jika ada inputannya
	if val := c.FormValue("tahun_lulus"); val != "" { 
		updates["tahun_lulus"] = val 
	}
	
	statusStr := c.FormValue("status_lulus")
	if statusStr != "" {
		updates["status_lulus"] = (statusStr == "true")
	}

	// Handle foto baru
	file, err := c.FormFile("foto_siswa")
	if err == nil {
		if siswa.FotoSiswa != "" {
			os.Remove(fmt.Sprintf("./public/uploads/siswa/%s", siswa.FotoSiswa))
		}
		
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", siswa.NISN, time.Now().Unix(), ext)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/siswa/%s", filename))
		updates["foto_siswa"] = filename
	}

	if err := config.DB.Model(&siswa).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update database"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Data siswa berhasil diupdate!",
		"data": siswa,
	})
}

func DeleteSiswa(c *fiber.Ctx) error {
	id := c.Params("id")
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

	var siswa models.Siswa
	if err := config.DB.Where("id = ? AND instansi_id = ?", id, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Siswa gak ada"})
	}

	if siswa.FotoSiswa != "" {
		os.Remove(fmt.Sprintf("./public/uploads/siswa/%s", siswa.FotoSiswa))
	}

	if err := config.DB.Unscoped().Delete(&siswa).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus permanen"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"message": "Siswa berhasil dimusnahkan!",
	})
}


func ImportSiswaExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "File gak ketemu!"})
	}

	val := c.Locals(utils.KeyInstansiID)
	parsedInstansiID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "sesi instansi tidak valid silahkan login ulang"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "gagal buka stream file"})
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal buka file excel"})
	}

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil || len(rows) <= 1 {
		return c.Status(400).JSON(fiber.Map{"message": "Sheet kosong bro"})
	}

	var listSiswa []models.Siswa

	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		tahunLulus := "2026"
		if len(row) > 7 && row[7] != "" {
			tahunLulus = row[7]
		}

		siswa := models.Siswa{
			InstansiID:   parsedInstansiID,
			NISN:         row[0],
			NamaSiswa:    row[1],
			TempatLahir:  row[2],
			TanggalLahir: row[3],
			JenisKelamin: row[4],
			Kelas:        row[5],
			Jurusan:      row[6],
			TahunLulus:   tahunLulus,
			StatusLulus:  false,
		}
		listSiswa = append(listSiswa, siswa)
	}

	if err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "nisn"}},
		DoNothing: true,
	}).Create(&listSiswa).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal simpan ke DB: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("%d data siswa berhasil diimport!", len(listSiswa)),
	})
}

// DeleteAllSiswa - Hapus semua data siswa untuk instansi tertentu (berdasarkan tahun lulus)
func DeleteAllSiswa(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	tahunLulus := c.Query("tahun_lulus")

	if tahunLulus == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Tahun lulus wajib dipilih"})
	}

	// Cek jumlah siswa yang akan dihapus
	var count int64
	if err := config.DB.Model(&models.Siswa{}).
		Where("instansi_id = ? AND tahun_lulus = ?", instansiID, tahunLulus).
		Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal cek data siswa"})
	}

	if count == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Tidak ada siswa dengan tahun lulus tersebut"})
	}

	// Hapus semua siswa (hard delete karena pake Unscoped)
	if err := config.DB.Unscoped().
		Where("instansi_id = ? AND tahun_lulus = ?", instansiID, tahunLulus).
		Delete(&models.Siswa{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus siswa: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("%d data siswa berhasil dihapus", count),
	})
}

// SearchSiswa - Pencarian siswa berdasarkan keyword
func SearchSiswa(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	keyword := c.Query("keyword")
	tahunLulus := c.Query("tahun_lulus")

	if keyword == "" {
		return GetAllSiswa(c) // balik ke semua data
	}

	var siswas []models.Siswa
	query := config.DB.Preload("Nilai").Where("instansi_id = ?", instansiID)

	if tahunLulus != "" {
		query = query.Where("tahun_lulus = ?", tahunLulus)
	}

	// Pencarian berdasarkan NISN atau Nama Siswa
	query = query.Where("nisn LIKE ? OR nama_siswa LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Find(&siswas).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mencari data"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   siswas,
	})
}

