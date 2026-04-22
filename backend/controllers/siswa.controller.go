package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllSiswa - Menampilkan semua siswa berdasarkan InstansiID dari middleware
// GetAllSiswa - Menampilkan semua siswa berdasarkan InstansiID & Filter Tahun dari Header
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

// GetSiswaByID - UNTUK CETAK & DETAIL: Mengambil 1 siswa berdasarkan ID & InstansiID
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

// CreateSiswa - Input manual siswa + Upload Foto
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

// UpdateSiswa - UPDATE SELEKTIF
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

// DeleteSiswa - Hapus permanen data & foto
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

// UploadFotoSiswa - ZIP Upload via Goroutines dengan tracking gagal
// UploadFotoSiswa - ZIP Upload via Goroutines dengan tracking gagal
func UploadFotoSiswa(c *fiber.Ctx) error {
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    file, err := c.FormFile("foto_zip")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Mana file ZIP-nya bro?"})
    }

    src, err := file.Open()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal buka ZIP"})
    }
    defer src.Close()

    zipReader, err := zip.NewReader(src, file.Size)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format file bukan ZIP"})
    }

    var wg sync.WaitGroup
    var successCount int32
    var failedList []string
    var mu sync.Mutex

    for _, f := range zipReader.File {
        if f.FileInfo().IsDir() || filepath.Base(f.Name)[0] == '.' {
            continue
        }

        wg.Add(1)

        go func(fileInZip *zip.File) {
            defer wg.Done()

            ext := filepath.Ext(fileInZip.Name)
            nisn := strings.TrimSuffix(filepath.Base(fileInZip.Name), ext)

            var siswa models.Siswa
            if err := config.DB.Where("nisn = ? AND instansi_id = ?", nisn, instansiID).First(&siswa).Error; err != nil {
                mu.Lock()
                failedList = append(failedList, nisn)
                mu.Unlock()
                return
            }

            imgFile, err := fileInZip.Open()
            if err != nil {
                mu.Lock()
                failedList = append(failedList, nisn+" (gagal baca)")
                mu.Unlock()
                return
            }
            defer imgFile.Close()

            newFilename := fmt.Sprintf("%s_%d%s", nisn, time.Now().UnixNano(), ext)
            targetPath := fmt.Sprintf("./public/uploads/siswa/%s", newFilename)

            dst, err := os.Create(targetPath)
            if err != nil {
                mu.Lock()
                failedList = append(failedList, nisn+" (gagal simpan)")
                mu.Unlock()
                return
            }
            defer dst.Close()

            if _, err := io.Copy(dst, imgFile); err != nil {
                mu.Lock()
                failedList = append(failedList, nisn+" (gagal copy)")
                mu.Unlock()
                return
            }

            if siswa.FotoSiswa != "" {
                os.Remove(fmt.Sprintf("./public/uploads/siswa/%s", siswa.FotoSiswa))
            }

            config.DB.Model(&siswa).Update("foto_siswa", newFilename)
            atomic.AddInt32(&successCount, 1)
        }(f)
    }

    wg.Wait()

    // ✅ URUTAN PENGECEKAN YANG BENAR
    // 1. Kalau gak ada yang sukses sama sekali → ERROR
    if successCount == 0 {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "Tidak ada foto yang berhasil diunggah. Seluruh NISN tidak terdaftar dalam sistem.",
        })
    }

    // 2. Kalau ada yang gagal (tapi ada yang sukses) → PARTIAL
    if len(failedList) > 0 {
        return c.JSON(fiber.Map{
            "status":        "partial",
            "message":       fmt.Sprintf("%d foto berhasil diunggah, %d foto gagal diproses.", successCount, len(failedList)),
            "success_count": successCount,
            "failed_count":  len(failedList),
            "failed_nisn":   failedList,
        })
    }

    // 3. Semua sukses → SUCCESS
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": fmt.Sprintf("Seluruh %d foto berhasil diunggah dan diproses.", successCount),
    })
}


func DeleteFotoSiswa(c *fiber.Ctx) error {
    id := c.Params("id")
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    var siswa models.Siswa
    if err := config.DB.Where("id = ? AND instansi_id = ?", id, instansiID).First(&siswa).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Siswa gak ketemu"})
    }

    if siswa.FotoSiswa != "" {
        os.Remove(fmt.Sprintf("./public/uploads/siswa/%s", siswa.FotoSiswa))
    }

    config.DB.Model(&siswa).Update("foto_siswa", "")

    return c.JSON(fiber.Map{"message": "Foto berhasil dihapus!"})
}


// GetSKLDataForAdmin mengambil data lengkap siswa dan template untuk dicetak oleh Admin
func GetSKLDataForAdmin(c *fiber.Ctx) error {
	instansiIDRaw := c.Locals(utils.KeyInstansiID)
	if instansiIDRaw == nil {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
	}
	instansiID := instansiIDRaw.(uuid.UUID)

	siswaIDParam := c.Params("id")
	siswaID, err := uuid.Parse(siswaIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID Siswa tidak valid"})
	}

	var siswa models.Siswa
	if err := config.DB.Preload("Nilai.Mapel").Where("id = ? AND instansi_id = ?", siswaID, instansiID).First(&siswa).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data siswa tidak ditemukan atau bukan milik sekolah ini"})
	}

	var template models.TemplateSKL
	if err := config.DB.Where("instansi_id = ?", instansiID).First(&template).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Template SKL belum dikonfigurasi untuk sekolah ini"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"siswa":    siswa, // Field tahun_lulus akan otomatis ada di sini
			"template": template,
		},
	})
}


func CheckGraduationEligibility(c *fiber.Ctx) error {
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    // 1. Ambil KKM dari Template SKL
    var template models.TemplateSKL
    if err := config.DB.Where("instansi_id = ?", instansiID).First(&template).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Template SKL/KKM belum disetting!"})
    }

    kkm := template.MinimalKelulusan // Asumsi field ini float64 di model TemplateSKL

    // 2. Ambil semua siswa aktif (belum lulus)
    var siswas []models.Siswa
    if err := config.DB.Preload("Nilai").Where("instansi_id = ? AND status_lulus = ?", instansiID, false).Find(&siswas).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal ambil data siswa"})
    }

    var layak []string       // List ID siswa yang layak
    var tidakLayak []fiber.Map // List detail siswa yang gagal

    for _, s := range siswas {
        if len(s.Nilai) == 0 {
            // Jika belum ada nilai sama sekali
            tidakLayak = append(tidakLayak, fiber.Map{
                "id":     s.ID.String(),
                "nisn":   s.NISN,
                "nama":   s.NamaSiswa,
                "alasan": "Belum ada input nilai",
            })
            continue
        }

        var total float64
        var hasFail bool
        var minNilai float64 = 100

        for _, n := range s.Nilai {
            total += n.NilaiAngka
            if n.NilaiAngka < minNilai {
                minNilai = n.NilaiAngka
            }
            // Syarat tambahan: Tidak boleh ada nilai di bawah 60 (Opsional, bisa dihapus jika hanya mau rata-rata)
            if n.NilaiAngka < 60 { 
                hasFail = true 
            }
        }

        rataRata := total / float64(len(s.Nilai))

        // LOGIKA KELULUSAN: Rata-rata >= KKM DAN Tidak ada nilai merah (<60)
        // Kamu bisa sesuaikan logika ini sesuai aturan sekolah
        if rataRata >= kkm && !hasFail {
            layak = append(layak, s.ID.String())
        } else {
            alasan := fmt.Sprintf("Rata-rata %.2f (Min: %.2f)", rataRata, kkm)
            if hasFail {
                alasan += fmt.Sprintf(", Ada nilai < 60 (Min: %.0f)", minNilai)
            }
            
            tidakLayak = append(tidakLayak, fiber.Map{
                "id":     s.ID.String(),
                "nisn":   s.NISN,
                "nama":   s.NamaSiswa,
                "rata":   rataRata,
                "alasan": alasan,
            })
        }
    }

    return c.JSON(fiber.Map{
        "status":        "success",
        "total_cek":     len(siswas),
        "layak_count":   len(layak),
        "invalid_count": len(tidakLayak),
        "layak_ids":     layak,          // Kirim ID ke frontend untuk mode "Hanya Layak"
        "invalid_list":  tidakLayak,     // Kirim detail kegagalan
        "kkm_used":      kkm,
    })
}

// ExecuteMassPromotion: Eksekusi update status lulus berdasarkan pilihan admin
func ExecuteMassPromotion(c *fiber.Ctx) error {
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    var input struct {
        TahunLulus string   `json:"tahun_lulus"`
        Mode       string   `json:"mode"` // "eligible_only" atau "force_all"
        StudentIDs []string `json:"student_ids"` // Digunakan jika mode eligible_only
        Reason     string   `json:"reason"` // Opsional: Alasan jika force_all
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
    }

    if input.TahunLulus == "" {
        return c.Status(400).JSON(fiber.Map{"message": "Tahun Lulus wajib diisi"})
    }

    query := config.DB.Model(&models.Siswa{}).Where("instansi_id = ?", instansiID)

    if input.Mode == "eligible_only" {
        // Update HANYA siswa yang ID-nya dikirim frontend (yang sudah divalidasi layak)
        if len(input.StudentIDs) == 0 {
            return c.JSON(fiber.Map{"status": "success", "message": "Tidak ada siswa layak untuk diproses."})
        }
        query = query.Where("id IN ?", input.StudentIDs)
    } else {
        // Force All: Update SEMUA siswa yang statusnya masih false (aktif)
        query = query.Where("status_lulus = ?", false)
    }

    result := query.Updates(map[string]interface{}{
        "status_lulus": true,
        "tahun_lulus":  input.TahunLulus,
    })

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal update database: " + result.Error.Error()})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": fmt.Sprintf("Berhasil! %d siswa telah dipromosikan menjadi lulusan tahun %s.", result.RowsAffected, input.TahunLulus),
        "count":   result.RowsAffected,
    })
}