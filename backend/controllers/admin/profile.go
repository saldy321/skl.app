package admin

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"skl-bakcend/controllers"

	"github.com/gofiber/fiber/v2"
)

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

	if err := controllers.CopyFile(targetPath, logoPath); err != nil {
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