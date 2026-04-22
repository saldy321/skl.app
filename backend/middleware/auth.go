package middleware

import (
	"fmt"
	"os"
	"skl-bakcend/models"
	"skl-bakcend/utils" 
	"gorm.io/gorm"
"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	
)

// Auth adalah satpam utama buat ngecek Token JWT
func Auth(c *fiber.Ctx) error {
	// 1. Ambil Token dari Cookie
	tokenString := c.Cookies("token") 

	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"message": "Akses ditolak, token anda tidak ada!"})
	}

	// 2. Parse & Validasi JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Token palsu atau expired!"})
	}

	// 3. Extract Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Gagal baca isi token"})
	}

	// Simpan data dasar ke Locals
	c.Locals(utils.KeyUserID, claims["id"])
	c.Locals(utils.KeyRole, claims["role"])
	c.Locals(utils.KeyInstansi, claims["nama_instansi"])
	c.Locals(utils.KeyTingkatSekolah, claims["tingkat_sekolah"])
	// Di dalam func Auth, setelah set KeyInstansiID
	c.Locals(utils.KeySlug, claims["slug"]) 



	// 4. Handle Instansi ID (Penting untuk membedakan Super Admin vs Admin Sekolah)
	instIDStr, _ := claims["instansi_id"].(string)
	
	// Jika instansi_id kosong atau "null" (khas Super Admin Hardcoded)
	if instIDStr == "" || instIDStr == uuid.Nil.String() {
		// Set ke Nil UUID agar konsisten
		c.Locals(utils.KeyInstansiID, uuid.Nil)
	} else {
		parsedID, err := uuid.Parse(instIDStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": "Format ID Instansi di token rusak!"})
		}
		c.Locals(utils.KeyInstansiID, parsedID)
	}

	return c.Next()
}

// IsSuperAdmin Middleware khusus untuk memastikan user adalah Super Admin
func IsSuperAdmin(c *fiber.Ctx) error {
	role := c.Locals(utils.KeyRole)

	if role != "super_admin" && role != models.RoleSuperAdmin {
		return c.Status(403).JSON(fiber.Map{"message": "Akses Ditolak: Khusus Super Admin!"})
	}
	return c.Next()
}

// Admin Middleware khusus untuk memastikan user adalah Admin Sekolah
func Admin(c *fiber.Ctx) error {
	role := c.Locals(utils.KeyRole)
	if role != "admin" && role != models.RoleAdmin {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses Ditolak: Khusus Admin Sekolah!",
		})
	}
	return c.Next()
}


func CheckTenant(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        log.Println("[CheckTenant] ========== MASUK ==========")
        log.Println("[CheckTenant] Path:", c.Path())
        log.Println("[CheckTenant] Method:", c.Method())
        
        slug := c.Params("slug")
        log.Println("[CheckTenant] Slug:", slug)
        
        // Jika tidak ada slug di URL, lewati saja
        if slug == "" {
            log.Println("[CheckTenant] Slug kosong, NEXT")
            return c.Next()
        }

        var instansi models.Instansi
        // Cari instansi berdasarkan slug dari URL
        err := db.Where("slug = ?", slug).Select("id, slug").First(&instansi).Error
        if err != nil {
            log.Println("[CheckTenant] ERROR: Instansi gak ketemu:", err)
            return c.Status(404).JSON(fiber.Map{
                "message": "URL Sekolah salah, instansi gak ketemu!",
            })
        }

        log.Println("[CheckTenant] Instansi ditemukan:", instansi.ID)
        // Simpan ID Instansi dari DATABASE ke Locals dengan key 'tenant_id'
        c.Locals("tenant_id", instansi.ID)

        // ✅ PERBAIKAN: Cek apakah ini route login siswa (publik)
        path := c.Path()
        if strings.Contains(path, "/login-siswa") {
            log.Println("[CheckTenant] Route login-siswa, SKIP AUTH, NEXT")
            // Login siswa tidak perlu cek token
            return c.Next()
        }

        log.Println("[CheckTenant] Bukan login-siswa, cek role...")
        
        // ========== BAGIAN YANG KURANG ==========
        valRole := c.Locals(utils.KeyRole)
        
        if valRole != nil {
            role := valRole.(string)

            // SUPER ADMIN BOLOS SEMUA SEKOLAH
            if role == "super_admin" || role == models.RoleSuperAdmin {
                log.Println("[CheckTenant] Super Admin, BOLOS")
                return c.Next()
            }

            valToken := c.Locals(utils.KeyInstansiID) 
            if valToken == nil {
                log.Println("[CheckTenant] ERROR: KeyInstansiID nil")
                return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid, login lagi!"})
            }
            
            tokenInstansiID := valToken.(uuid.UUID)

            // VERIFIKASI: Apakah ID di URL sama dengan ID di Token?
            if instansi.ID != tokenInstansiID {
                log.Println("[CheckTenant] ERROR: ID mismatch")
                return c.Status(403).JSON(fiber.Map{
                    "message": "Anda tidak memiliki akses ke sekolah ini!",
                })
            }
        }
        // ========== END BAGIAN YANG KURANG ==========

        log.Println("[CheckTenant] SELESAI, NEXT")
        return c.Next()
    }
}