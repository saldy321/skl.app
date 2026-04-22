package controllers

import (
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"fmt"
	"gorm.io/gorm/clause"
    "time"
    "os"
    "log"
    "strings"
        "path/filepath"
        "io"
)

// ==================== PROFILE & DASHBOARD ====================

// GetMyProfile - Ambil profil user yang sedang login
func GetMyProfile(c *fiber.Ctx) error {
	userID, okID := c.Locals(utils.KeyUserID).(string)
	role, okRole := c.Locals(utils.KeyRole).(string)

	if !okID || !okRole {
		return c.Status(401).JSON(fiber.Map{"message": "Session abis, login lagi bro"})
	}

	// SUPER ADMIN (Hardcoded)
	if role == models.RoleSuperAdmin {
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":             userID,
				"email":          "superadmin@gmail.com",
				"role":           role,
				"nama_instansi":  "Super Admin System",
                "tingkat":         "pusat",
				"slug":           "pusat",
				"instansi_id":    uuid.Nil.String(),
                
			},
		})
	}

	// ADMIN SEKOLAH
	if role == models.RoleAdmin {
		var admin models.Admin
		if err := config.DB.Preload("Instansi").First(&admin, "id = ?", userID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "Admin gak ada"})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"id":             admin.ID.String(),
				"email":          admin.Email,
				"role":           models.RoleAdmin,
				"nama_instansi":  admin.Instansi.NamaInstansi,
                "tingkat":        admin.Instansi.TingkatSekolah,
				"slug":           admin.Instansi.Slug,
				"instansi_id":    admin.InstansiID.String(),
                "foto_profile":   admin.FotoProfile,
                 "logo_instansi":  admin.Instansi.LogoInstansi,
			},
		})
	}

	return c.Status(403).JSON(fiber.Map{"message": "Role gak dikenali"})
}

// GetSuperDashboard - Dashboard Super Admin
func GetSuperDashboard(c *fiber.Ctx) error {
	role := c.Locals(utils.KeyRole).(string)
	if role != models.RoleSuperAdmin {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak, khusus Super Admin!"})
	}

	var totalSekolah, totalAdmin, totalSiswaNasional int64

	config.DB.Model(&models.Instansi{}).Count(&totalSekolah)
	config.DB.Model(&models.Admin{}).Count(&totalAdmin)
	config.DB.Model(&models.Siswa{}).Count(&totalSiswaNasional)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"totalSekolah":       totalSekolah,
			"totalAdmin":         totalAdmin,
			"totalSiswaNasional": totalSiswaNasional,
		},
	})
}

// GetAdminDashboard - Dashboard Admin Sekolah
func GetAdminDashboard(c *fiber.Ctx) error {
	val := c.Locals(utils.KeyInstansiID)
	instansiID, ok := val.(uuid.UUID)

	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi instansi tidak valid, login ulang bro"})
	}

	var totalSiswa, siswaLulus int64

	config.DB.Model(&models.Siswa{}).Where("instansi_id = ?", instansiID).Count(&totalSiswa)
	config.DB.Model(&models.Siswa{}).Where("instansi_id = ? AND status_lulus = ?", instansiID, true).Count(&siswaLulus)

	return c.JSON(fiber.Map{
		"totalSiswa": totalSiswa,
		"siswaLulus": siswaLulus,
		"belumNilai": 0,
	})
}

// ==================== MANAJEMEN ADMIN ====================

// GetAllAdminSekolah - Lihat semua admin (untuk Super Admin)
func GetAllAdminSekolah(c *fiber.Ctx) error {
	var admins []models.Admin
	if err := config.DB.Preload("Instansi").Find(&admins).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal tarik data admin"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   admins,
	})
}

// RegisterAdminSekolah - Daftar admin baru (untuk Super Admin)
func RegisterAdminSekolah(c *fiber.Ctx) error {
	var input struct {
		Email      string    `json:"email"`
		Password   string    `json:"password"`
		InstansiID uuid.UUID `json:"instansi_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input data ngaco, cek lagi bro"})
	}

	if input.Email == "" || input.Password == "" || input.InstansiID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email, password, dan instansi wajib diisi!"})
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	newAdmin := models.Admin{
		Email:      input.Email,
		Password:   hashedPassword,
		InstansiID: input.InstansiID,
	}

	if err := config.DB.Create(&newAdmin).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal daftar admin: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Admin sekolah berhasil ditambahkan!",
		"data":    newAdmin,
	})
}

// UpdateAdminEmail - Edit email admin (untuk Super Admin)
func UpdateAdminEmail(c *fiber.Ctx) error {
	id := c.Params("id")

	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	if input.Email == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Email tidak boleh kosong"})
	}

	var admin models.Admin
	if err := config.DB.First(&admin, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Admin tidak ditemukan"})
	}

	if err := config.DB.Model(&admin).Update("email", input.Email).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update email: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Email admin berhasil diupdate",
	})
}

// ==================== MANAJEMEN SISWA ====================

// ImportSiswaExcel - Import data siswa dari file Excel
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

// ==================== MANAJEMEN NISN (SUPER ADMIN) ====================

// SearchNISN - Cari NISN di seluruh sistem (Hanya Super Admin)
func SearchNISN(c *fiber.Ctx) error {
    // Cek role Super Admin
    role := c.Locals(utils.KeyRole).(string)
    if role != models.RoleSuperAdmin {
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
    }

    nisn := c.Query("nisn")
    if nisn == "" {
        return c.Status(400).JSON(fiber.Map{"message": "NISN wajib diisi"})
    }

    var result struct {
        ID           uuid.UUID `json:"id"`
        NISN         string    `json:"nisn"`
        NamaSiswa    string    `json:"nama_siswa"`
        NamaInstansi string    `json:"nama_instansi"`
        InstansiID   uuid.UUID `json:"instansi_id"`
        Kelas        string    `json:"kelas"`
        Jurusan      string    `json:"jurusan"`
        StatusLulus  bool      `json:"status_lulus"`
        CreatedAt    time.Time `json:"created_at"`
        FotoSiswa    string    `json:"foto_siswa"`
    }

    err := config.DB.Table("siswas").
        Select("siswas.id, siswas.nisn, siswas.nama_siswa, instansis.nama_instansi, instansis.id as instansi_id, siswas.kelas, siswas.jurusan, siswas.status_lulus, siswas.created_at, siswas.foto_siswa").
        Joins("JOIN instansis ON siswas.instansi_id = instansis.id").
        Where("siswas.nisn = ? AND siswas.deleted_at IS NULL", nisn).
        First(&result).Error

    if err != nil {
        return c.JSON(fiber.Map{
            "found":   false,
            "message": "NISN tidak ditemukan di sistem",
        })
    }

    return c.JSON(fiber.Map{
        "found": true,
        "data":  result,
    })
}

// ForceDeleteNISN - Hapus permanen NISN (Hanya Super Admin)
func ForceDeleteNISN(c *fiber.Ctx) error {
    // Cek role Super Admin
    role := c.Locals(utils.KeyRole).(string)
    if role != models.RoleSuperAdmin {
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin."})
    }

    superAdminID := c.Locals(utils.KeyUserID).(string)

    var input struct {
        SiswaID string `json:"siswa_id"`
        Alasan  string `json:"alasan"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
    }

    if input.SiswaID == "" || input.Alasan == "" {
        return c.Status(400).JSON(fiber.Map{"message": "ID Siswa dan Alasan wajib diisi"})
    }

    // Cari siswa
    var siswa models.Siswa
    if err := config.DB.First(&siswa, "id = ?", input.SiswaID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Siswa tidak ditemukan"})
    }

    // Hapus foto jika ada
    if siswa.FotoSiswa != "" {
        os.Remove(fmt.Sprintf("./public/uploads/siswa/%s", siswa.FotoSiswa))
    }

    // Hapus nilai terkait
    config.DB.Unscoped().Where("siswa_id = ?", siswa.ID).Delete(&models.Nilai{})

    // Hapus siswa permanen (Unscoped untuk hard delete karena lo pake gorm.Model)
    if err := config.DB.Unscoped().Delete(&siswa).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus data: " + err.Error()})
    }

    // Catat log audit
    log.Printf("[AUDIT] Super Admin %s menghapus permanen NISN %s (Siswa: %s, Sekolah: %s). Alasan: %s",
        superAdminID, siswa.NISN, siswa.NamaSiswa, siswa.InstansiID, input.Alasan)

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": fmt.Sprintf("NISN %s berhasil dihapus permanen dari sistem.", siswa.NISN),
    })
}

func UploadFotoProfile(c *fiber.Ctx) error {
	userID := c.Locals(utils.KeyUserID).(string)
	role := c.Locals(utils.KeyRole).(string)

	log.Println("[UPLOAD] ========== START UPLOAD ==========")
	log.Println("[UPLOAD] UserID:", userID)
	log.Println("[UPLOAD] Role:", role)

	if role != models.RoleAdmin && role != models.RoleSuperAdmin {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak"})
	}

	file, err := c.FormFile("foto")
	if err != nil {
		log.Println("[UPLOAD] Error get file:", err)
		return c.Status(400).JSON(fiber.Map{"message": "File tidak ditemukan"})
	}

	log.Println("[UPLOAD] File:", file.Filename, "Size:", file.Size)

	if file.Size > 2*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"message": "Ukuran file maksimal 2MB"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExt[ext] {
		return c.Status(400).JSON(fiber.Map{"message": "Format file harus JPG atau PNG"})
	}

	filename := fmt.Sprintf("admin_%s_%d%s", userID[:8], time.Now().Unix(), ext)
	targetPath := fmt.Sprintf("./public/uploads/admin/%s", filename)
	logoPath := fmt.Sprintf("./public/uploads/instansi/%s", filename)

	log.Println("[UPLOAD] Target path:", targetPath)
	log.Println("[UPLOAD] Logo path:", logoPath)

	os.MkdirAll("./public/uploads/admin", 0755)
	os.MkdirAll("./public/uploads/instansi", 0755)

	if err := c.SaveFile(file, targetPath); err != nil {
		log.Println("[UPLOAD] Gagal simpan file:", err)
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan foto"})
	}
	log.Println("[UPLOAD] File tersimpan di:", targetPath)

	if err := copyFile(targetPath, logoPath); err != nil {
		log.Println("[UPLOAD] Gagal copy logo:", err)
	} else {
		log.Println("[UPLOAD] Logo tersalin ke:", logoPath)
	}

	if role == models.RoleAdmin {
		var admin models.Admin
		if err := config.DB.Preload("Instansi").First(&admin, "id = ?", userID).Error; err != nil {
			log.Println("[UPLOAD] Admin tidak ditemukan:", err)
			return c.Status(404).JSON(fiber.Map{"message": "Admin tidak ditemukan"})
		}

		log.Println("[UPLOAD] Admin ditemukan:", admin.Email)
		log.Println("[UPLOAD] Instansi ID:", admin.InstansiID)

		if admin.FotoProfile != "" {
			oldPath := fmt.Sprintf("./public/uploads/admin/%s", admin.FotoProfile)
			os.Remove(oldPath)
			log.Println("[UPLOAD] Hapus foto lama:", oldPath)
		}

		if admin.Instansi.LogoInstansi != "" {
			oldLogoPath := fmt.Sprintf("./public/uploads/instansi/%s", admin.Instansi.LogoInstansi)
			os.Remove(oldLogoPath)
			log.Println("[UPLOAD] Hapus logo lama:", oldLogoPath)
		}

		resultAdmin := config.DB.Model(&models.Admin{}).
			Where("id = ?", userID).
			Update("foto_profile", filename)

		log.Println("[UPLOAD] Update admin - RowsAffected:", resultAdmin.RowsAffected, "Error:", resultAdmin.Error)

		resultInstansi := config.DB.Model(&models.Instansi{}).
			Where("id = ?", admin.InstansiID).
			Update("logo_instansi", filename)

		log.Println("[UPLOAD] Update instansi - RowsAffected:", resultInstansi.RowsAffected, "Error:", resultInstansi.Error)

		var verifyAdmin models.Admin
		config.DB.First(&verifyAdmin, "id = ?", userID)
		log.Println("[UPLOAD] VERIFIKASI - Admin foto_profile setelah update:", verifyAdmin.FotoProfile)

		var verifyInstansi models.Instansi
		config.DB.First(&verifyInstansi, "id = ?", admin.InstansiID)
		log.Println("[UPLOAD] VERIFIKASI - Instansi logo_instansi setelah update:", verifyInstansi.LogoInstansi)
	}

	log.Println("[UPLOAD] ========== UPLOAD SELESAI ==========")

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Foto profil & logo sekolah berhasil diupload",
		"data": fiber.Map{
			"foto_profile": filename,
			"logo_instansi": filename,
		},
	})
}

// DeleteFotoProfile - Hapus foto profil & logo sekolah
func DeleteFotoProfile(c *fiber.Ctx) error {
	userID := c.Locals(utils.KeyUserID).(string)
	role := c.Locals(utils.KeyRole).(string)

	if role == models.RoleAdmin {
		var admin models.Admin
		config.DB.Preload("Instansi").First(&admin, "id = ?", userID)

		if admin.FotoProfile != "" {
			os.Remove(fmt.Sprintf("./public/uploads/admin/%s", admin.FotoProfile))
			config.DB.Model(&admin).Update("foto_profile", "")
		}

		if admin.Instansi.LogoInstansi != "" {
			os.Remove(fmt.Sprintf("./public/uploads/instansi/%s", admin.Instansi.LogoInstansi))
			config.DB.Model(&models.Instansi{}).Where("id = ?", admin.InstansiID).Update("logo_instansi", "")
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Foto profil & logo sekolah berhasil dihapus",
	})
}

// copyFile - Helper function untuk copy file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}