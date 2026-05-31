package admin

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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