package super_admin

import (
	"fmt"
	"log"
	"os"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
func DeleteNISN(c *fiber.Ctx) error {
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