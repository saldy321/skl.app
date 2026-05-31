package admin

import (
	"fmt"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)




func CekKelayakanLulus(c *fiber.Ctx) error {
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

func ProsesLulusMasal(c *fiber.Ctx) error {
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