    package controllers

    import (
        "skl-bakcend/config"
        "skl-bakcend/models"
        "skl-bakcend/utils"
        "github.com/gofiber/fiber/v2"
        "github.com/google/uuid"
        "log"
        "time"
        "strings"
    )

    // RegisterFullSchool - CREATE instansi + admin sekaligus
    func RegisterFullSchool(c *fiber.Ctx) error {
        var input struct {
            NamaInstansi   string `json:"nama_instansi"`
            KodeInstansi   string `json:"kode_instansi"`
            TingkatSekolah string `json:"tingkat_sekolah"`
            Alamat         string `json:"alamat"`
            Email          string `json:"email"`
            Password       string `json:"password"`
        }

        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format input ngaco, bro"})
        }

        if input.NamaInstansi == "" || input.Email == "" || input.Password == "" {
            return c.Status(400).JSON(fiber.Map{"message": "Data sekolah dan admin wajib diisi semua!"})
        }

        tx := config.DB.Begin()

        var cek models.Instansi
        if err := tx.Where("kode_instansi = ?", input.KodeInstansi).First(&cek).Error; err == nil {
            tx.Rollback()
            return c.Status(400).JSON(fiber.Map{"message": "Kode sekolah ini udah kedaftar!"})
        }

        // Generate slug otomatis dari nama
        slug := strings.ToLower(strings.ReplaceAll(input.NamaInstansi, " ", "-"))

        newInstansi := models.Instansi{
            NamaInstansi:   input.NamaInstansi,
            KodeInstansi:   input.KodeInstansi,
            TingkatSekolah: input.TingkatSekolah,
            Alamat:         input.Alamat,
            Slug:           slug,
        }

        if err := tx.Create(&newInstansi).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin sekolah: " + err.Error()})
        }

        hashedPassword, _ := utils.HashPassword(input.Password)
        newAdmin := models.Admin{
            Email:      input.Email,
            Password:   hashedPassword,
            InstansiID: newInstansi.ID,
        }

        if err := tx.Create(&newAdmin).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin akun admin (Email mungkin duplikat)"})
        }

        tx.Commit()

        // Kirim email di background
        go func(email, nama, slug, pass string) {
            err := utils.SendCredentialEmail(email, nama, slug, pass)
            if err != nil {
                log.Printf("[EMAIL ERROR] Gagal kirim ke %s: %v", email, err)
            } else {
                log.Printf("[EMAIL SUCCESS] Kredensial terkirim ke %s", email)
            }
        }(input.Email, input.NamaInstansi, input.KodeInstansi, input.Password)

        return c.Status(201).JSON(fiber.Map{
            "status":  "success",
            "message": "Mantap! Sekolah & Admin berhasil didaftarin sekaligus!",
            "data": fiber.Map{
                "instansi": newInstansi,
                "admin":    newAdmin.Email,
            },
        })
    }

    // GetAllInstansi - READ semua instansi
    func GetAllInstansi(c *fiber.Ctx) error {
        var instansis []models.Instansi
        config.DB.Find(&instansis)
        return c.JSON(fiber.Map{"status": "success", "data": instansis})
    }

    // GetInstansiByID - READ instansi by ID
    func GetInstansiByID(c *fiber.Ctx) error {
        id := c.Params("id")
        var instansi models.Instansi
        if err := config.DB.First(&instansi, "id = ?", id).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"message": "Sekolah gak ketemu!"})
        }
        return c.JSON(fiber.Map{"status": "success", "data": instansi})
    }

    // UpdateInstansi - UPDATE instansi (lengkap dengan slug)
    func UpdateInstansi(c *fiber.Ctx) error {
        id := c.Params("id")
        var instansi models.Instansi

        if err := config.DB.First(&instansi, "id = ?", id).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"message": "Sekolahnya gak ketemu, mau edit apaan?"})
        }

        var input struct {
            NamaInstansi   string `json:"nama_instansi"`
            KodeInstansi   string `json:"kode_instansi"`
            Slug           string `json:"slug"`
            TingkatSekolah string `json:"tingkat_sekolah"`
            Alamat         string `json:"alamat"`
        }

        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Input edit lo ada yang salah nih, bro"})
        }

        updates := map[string]interface{}{
            "nama_instansi":   input.NamaInstansi,
            "kode_instansi":   input.KodeInstansi,
            "tingkat_sekolah": input.TingkatSekolah,
            "alamat":          input.Alamat,
            "slug":            input.Slug,
        }

        if err := config.DB.Model(&instansi).Updates(updates).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal update data sekolah: " + err.Error()})
        }

        config.DB.First(&instansi, "id = ?", id)

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Data sekolah berhasil di-update!",
            "data":    instansi,
        })
    }

    // DeleteInstansi - HARD DELETE permanen
    func DeleteInstansi(c *fiber.Ctx) error {
        id := c.Params("id")

        tx := config.DB.Begin()

        // Hapus data terkait
        if err := tx.Unscoped().Where("instansi_id = ?", id).Delete(&models.Siswa{}).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus data siswa: " + err.Error()})
        }

        if err := tx.Unscoped().Where("instansi_id = ?", id).Delete(&models.Nilai{}).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus data nilai: " + err.Error()})
        }

        if err := tx.Unscoped().Where("instansi_id = ?", id).Delete(&models.Admin{}).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus data admin: " + err.Error()})
        }

        if err := tx.Unscoped().Where("instansi_id = ?", id).Delete(&models.Mapel{}).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus data mapel: " + err.Error()})
        }

        if err := tx.Unscoped().Delete(&models.Instansi{}, "id = ?", id).Error; err != nil {
            tx.Rollback()
            return c.Status(500).JSON(fiber.Map{"message": "Gagal hapus instansi: " + err.Error()})
        }

        tx.Commit()

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Instansi dan semua data terkait berhasil dihapus permanen!",
        })
    }

    // ResetPasswordAdminByInstansi - Reset password admin via InstansiID
    func ResetPasswordAdminByInstansi(c *fiber.Ctx) error {
        role := c.Locals(utils.KeyRole).(string)
        if role != models.RoleSuperAdmin {
            return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak!"})
        }

        instansiID := c.Params("instansi_id")
        if instansiID == "" {
            return c.Status(400).JSON(fiber.Map{"message": "ID Instansi wajib"})
        }

        var input struct {
            NewPassword string `json:"new_password"`
        }

        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
        }

        if len(input.NewPassword) < 6 {
            return c.Status(400).JSON(fiber.Map{"message": "Password minimal 6 karakter"})
        }

        var admin models.Admin
        if err := config.DB.Where("instansi_id = ?", instansiID).First(&admin).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"message": "Admin untuk sekolah ini tidak ditemukan"})
        }

        hashedPassword, _ := utils.HashPassword(input.NewPassword)

        if err := config.DB.Model(&admin).Update("password", hashedPassword).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal update password"})
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Password admin berhasil direset!",
        })
    }

    // SetWaktuBukaPengumuman - Atur waktu buka pengumuman
    func SetWaktuBukaPengumuman(c *fiber.Ctx) error {
        instansiIDRaw := c.Locals(utils.KeyInstansiID)
        if instansiIDRaw == nil {
            return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
        }
        instansiID := instansiIDRaw.(uuid.UUID)

        var input struct {
            WaktuBuka string `json:"waktu_buka"`
        }

        if err := c.BodyParser(&input); err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
        }

        loc, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal load timezone server"})
        }

        waktu, err := time.ParseInLocation("2006-01-02 15:04:05", input.WaktuBuka, loc)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"message": "Format waktu harus YYYY-MM-DD HH:MM:SS"})
        }

        if err := config.DB.Model(&models.Instansi{}).Where("id = ?", instansiID).Update("waktu_buka_pengumuman", waktu).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"message": "Gagal update jadwal"})
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Jadwal pengumuman berhasil diatur!",
        })
    }

    // CheckInstansi - Cek instansi berdasarkan slug (publik)
  // CheckInstansi - Cek instansi berdasarkan slug (publik)
func CheckInstansi(c *fiber.Ctx) error {
    slug := c.Params("slug")
    var instansi models.Instansi
    
    // ✅ TAMBAHKAN logo_instansi & tampilkan_logo di SELECT
    if err := config.DB.Select("id, nama_instansi, tingkat_sekolah, slug, waktu_buka_pengumuman, logo_instansi, tampilkan_logo").Where("slug = ?", slug).First(&instansi).Error; err != nil {
        return c.Status(200).JSON(fiber.Map{
            "success": false,
            "message": "Sekolah tidak terdaftar",
        })
    }

    loc, _ := time.LoadLocation("Asia/Jakarta")
    sekarang := time.Now().In(loc)

    isClosed := true
    var sisaWaktu int64 = 0

    if instansi.WaktuBukaPengumuman != nil {
        waktuBukaWIB := instansi.WaktuBukaPengumuman.In(loc)

        if sekarang.After(waktuBukaWIB) {
            isClosed = false
        } else {
            sisaWaktu = int64(waktuBukaWIB.Sub(sekarang).Seconds())
        }
    } else {
        isClosed = false
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "nama_instansi":   instansi.NamaInstansi,
            "tingkat":         instansi.TingkatSekolah,
            "is_closed":       isClosed,
            "countdown":       sisaWaktu,
            "logo_instansi":   instansi.LogoInstansi,   // ✅ TAMBAHKAN
            "tampilkan_logo":  instansi.TampilkanLogo,  // ✅ TAMBAHKAN
        },
    })
}   
    // ImpersonateAsAdmin - Super admin login sebagai admin sekolah tertentu
func ImpersonateAsAdmin(c *fiber.Ctx) error {
    // 1. Cek role user yang login harus super admin
    role := c.Locals(utils.KeyRole).(string)
    if role != models.RoleSuperAdmin {
        return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak! Khusus Super Admin"})
    }

    // 2. Ambil instansi ID dari parameter URL
    instansiID := c.Params("instansi_id")
    if instansiID == "" {
        return c.Status(400).JSON(fiber.Map{"message": "ID Instansi wajib diisi"})
    }

    parsedInstansiID, err := uuid.Parse(instansiID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format ID Instansi tidak valid"})
    }

    // 3. Cari instansi
    var instansi models.Instansi
    if err := config.DB.First(&instansi, "id = ?", parsedInstansiID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Instansi tidak ditemukan"})
    }

    // 4. Cari admin sekolah terkait
    var admin models.Admin
    if err := config.DB.Where("instansi_id = ?", parsedInstansiID).First(&admin).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"message": "Admin untuk sekolah ini tidak ditemukan"})
    }

    // 5. Generate token baru dengan role "admin"
    token, err := utils.GenerateToken(
        admin.ID.String(),
        models.RoleAdmin,
        instansi.TingkatSekolah,
        instansi.NamaInstansi,
        parsedInstansiID,
        instansi.Slug,
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal bikin token impersonate"})
    }

    // 6. Set cookie baru
    SetAuthCookie(c, token)

    // 7. Return response
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Berhasil masuk sebagai admin sekolah",
        "slug":    instansi.Slug,
        "role":    models.RoleAdmin,
    })
}

// UpdateTampilkanLogo - Toggle tampilan logo di halaman login siswa
func UpdateTampilkanLogo(c *fiber.Ctx) error {
    instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)

    var input struct {
        TampilkanLogo bool `json:"tampilkan_logo"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
    }

    result := config.DB.Model(&models.Instansi{}).
        Where("id = ?", instansiID).
        Update("tampilkan_logo", input.TampilkanLogo)

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{"message": "Gagal update setting"})
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Setting logo berhasil diupdate",
        "data": fiber.Map{
            "tampilkan_logo": input.TampilkanLogo,
        },
    })
}