package auth

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"log"

	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"skl-bakcend/controllers"
)

const (
	SUPER_ADMIN_EMAIL =	"superadmin@gmail.com"
	SUPER_ADMIN_PASS = "super220508"
)

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
			controllers.SetAuthCookie(c, token)
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

	go SendOTPEmail(cleanEmail, otpCode)

	return c.JSON(fiber.Map{
		"status":  "otp_required",
		"message": "Kode OTP meluncur ke email lu!",
		"email":   cleanEmail,
	})
}
