package auth

import (
	"strings"
	"time"

	"skl-bakcend/config"
	"skl-bakcend/controllers"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func VerifyOTP(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input nggak valid"})
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(input.Email))
	var userID, userRole, tingkatSekolah, namaInstansi, slug string
	var instansiID interface{}
	var dbOtp string
	var dbOtpTime *time.Time

	var super models.SuperAdmin
	var admin models.Admin

	if err := config.DB.Where("LOWER(email) = ?", cleanEmail).First(&super).Error; err == nil {
		dbOtp, dbOtpTime = super.Otp, super.OtpCreatedAt
		userID, userRole = super.ID.String(), models.RoleSuperAdmin
		tingkatSekolah, namaInstansi, slug = models.TingkatPusat, "Pusat Sistem", "pusat"
		instansiID = uuid.Nil.String()
	} else if err := config.DB.Preload("Instansi").Where("LOWER(email) = ?", cleanEmail).First(&admin).Error; err == nil {
		dbOtp, dbOtpTime = admin.Otp, admin.OtpCreatedAt
		userID, userRole = admin.ID.String(), models.RoleAdmin
		tingkatSekolah, namaInstansi = admin.Instansi.TingkatSekolah, admin.Instansi.NamaInstansi
		instansiID = admin.InstansiID
		slug = admin.Instansi.Slug
	} else {
		return c.Status(404).JSON(fiber.Map{"message": "User kaga ketemu"})
	}

	if dbOtp == "" || dbOtp != input.OTP {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Kode OTP salah!"})
	}

	if dbOtpTime == nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Data waktu OTP kosong"})
	}

	duration := time.Since(*dbOtpTime)
	if duration.Minutes() > 5 {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "OTP udah expired!"})
	}

	updateMap := map[string]interface{}{"otp": "", "otp_created_at": nil}
	if userRole == models.RoleSuperAdmin {
		config.DB.Model(&models.SuperAdmin{}).Where("id = ?", userID).Updates(updateMap)
	} else {
		config.DB.Model(&models.Admin{}).Where("id = ?", userID).Updates(updateMap)
	}

	token, err := utils.GenerateToken(userID, userRole, tingkatSekolah, namaInstansi, instansiID, slug)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin token"})
	}

	controllers.SetAuthCookie(c, token)

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"success":       true,
		"message":       "Login Berhasil!",
		"role":          userRole,
		"nama_instansi": namaInstansi,
		"tingkat":		tingkatSekolah,
		"instansi_id":   instansiID,
		"slug":          slug,
	})
}