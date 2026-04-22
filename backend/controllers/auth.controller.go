package controllers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"
	"strings"
	"time"
	"log"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

const (
	SUPER_ADMIN_EMAIL = "superadmin@gmail.com"
	SUPER_ADMIN_PASS  = "super220508"
)


func SetAuthCookie(c *fiber.Ctx, token string) {

	c.ClearCookie("token")
	
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HTTPOnly = true
	cookie.Secure = false
	cookie.SameSite = "None"
	cookie.MaxAge = 86400 
	c.Cookie(cookie)
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
	}

	cleanEmail := strings.ToLower(strings.TrimSpace(input.Email))

	if cleanEmail == SUPER_ADMIN_EMAIL {
		if input.Password == SUPER_ADMIN_PASS {
			fakeUUID := uuid.New()
			token, err := utils.GenerateToken(
				fakeUUID.String(),
				models.RoleSuperAdmin,
				"Pusat",
				"Super Admin System",
				uuid.Nil,
				"pusat",
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin token"})
			}
			SetAuthCookie(c, token)
			return c.JSON(fiber.Map{
				"status":        "success",
				"success":       true,
				"message":       "Login Super Admin Berhasil!",
				"role":          models.RoleSuperAdmin,
				"nama_instansi": "Super Admin System",
				"tingkat" 		: "pusat",
				"slug":          "pusat",
				"instansi_id":   uuid.Nil.String(),
			})
		} else {
			return c.Status(401).JSON(fiber.Map{"message": "Password Super Admin salah!"})
		}
	}

	var hashedPassword string
	var userRole string
	var userID string

	var super models.SuperAdmin
	var admin models.Admin

	log.Println("[LOGIN] ========== START LOGIN ==========")
	log.Println("[LOGIN] Email:", cleanEmail)

	errSuper := config.DB.Where("LOWER(email) = ?", cleanEmail).First(&super).Error
	if errSuper == nil {
		log.Println("[LOGIN] ✅ Ketemu di super_admins")
		userRole = models.RoleSuperAdmin
		hashedPassword = super.Password
		userID = super.ID.String()
	} else {
		errAdmin := config.DB.Where("LOWER(email) = ?", cleanEmail).First(&admin).Error
		if errAdmin == nil {
			log.Println("[LOGIN] ✅ Ketemu di admins!")
			userRole = models.RoleAdmin
			hashedPassword = admin.Password
			userID = admin.ID.String()
		} else {
			return c.Status(404).JSON(fiber.Map{"message": "Email nggak kedaftar, bro"})
		}
	}

	if !utils.CheckPassword(input.Password, hashedPassword) {
		return c.Status(401).JSON(fiber.Map{"message": "Password salah!"})
	}

	rand.Seed(time.Now().UnixNano())
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	now := time.Now()
	updateData := map[string]interface{}{
		"otp":            otpCode,
		"otp_created_at": &now,
	}

	var errUpdate error
	if userRole == models.RoleSuperAdmin {
		errUpdate = config.DB.Model(&models.SuperAdmin{}).Where("id = ?", userID).Updates(updateData).Error
	} else {
		errUpdate = config.DB.Model(&models.Admin{}).Where("id = ?", userID).Updates(updateData).Error
	}

	if errUpdate != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal generate OTP"})
	}

	go sendOTPEmail(cleanEmail, otpCode)

	return c.JSON(fiber.Map{
		"status":  "otp_required",
		"message": "Kode OTP meluncur ke email lu!",
		"email":   cleanEmail,
	})
}

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

	SetAuthCookie(c, token)

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

func sendOTPEmail(toEmail string, code string) {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	jktTime := time.Now().In(time.FixedZone("WIB", 7*3600)).Format("15:04")
	msg := fmt.Sprintf("Subject: OTP SKL Digital [%s]\r\n\r\n"+
		"Kode OTP lu: %s\r\nIngat bro, kode ini cuma aktif 5 menit!",
		jktTime, code)

	auth := smtp.PlainAuth("", from, password, host)
	_ = smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logout berhasil!",
	})
}

// RestoreSuperAdmin - Kembalikan session ke Super Admin (untuk exit impersonate)
func RestoreSuperAdmin(c *fiber.Ctx) error {
    // Generate token baru dengan role super admin
    fakeUUID := uuid.New()
    token, err := utils.GenerateToken(
        fakeUUID.String(),
        models.RoleSuperAdmin,
        "Pusat",
        "Super Admin System",
        uuid.Nil,
        "pusat",
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin token"})
    }
    
    SetAuthCookie(c, token)
    
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Kembali ke mode Super Admin",
        "role":    models.RoleSuperAdmin,
    })
}	
