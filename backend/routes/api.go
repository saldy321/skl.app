package routes

import (
    "skl-bakcend/config"
    "skl-bakcend/controllers/auth"
    "skl-bakcend/controllers/super_admin"
     "skl-bakcend/controllers"
    "skl-bakcend/middleware"
      "skl-bakcend/controllers/siswa"
    "skl-bakcend/controllers/admin"

    "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
  
    app.Post("/api/login", auth.Login)
    app.Post("/api/verify-otp", auth.VerifyOTP)
    app.Post("/api/restore-super-admin", super_admin.RestoreSuperAdmin)
    app.Get("/api/check-instansi/:slug", controllers.CheckInstansi)

  
    superAuth := app.Group("/api/super", middleware.Auth, middleware.IsSuperAdmin)
    superAuth.Get("/dashboard", super_admin.GetSuperDashboard)
    superAuth.Post("/instansi", super_admin.RegisterFullSchool)
    superAuth.Get("/instansi", super_admin.GetAllInstansi)
    superAuth.Get("/instansi/:id", super_admin.GetInstansiByID)
    superAuth.Put("/instansi/:id", super_admin.UpdateInstansi)
    superAuth.Delete("/instansi/:id", super_admin.DeleteInstansi)
    superAuth.Post("/admin-register", super_admin.RegisterAdminSekolah)
    superAuth.Get("/admin-list", super_admin.GetAllAdminSekolah)
    superAuth.Post("/instansi/:instansi_id/reset-password", super_admin.ResetPasswordAdminByInstansi)
    superAuth.Put("/admin/:id/email", super_admin.UpdateAdminEmail)
    superAuth.Post("/impersonate/:instansi_id", super_admin.ImpersonateAsAdmin)
    superAuth.Get("/nisn/search", super_admin.SearchNISN)
    superAuth.Post("/nisn/force-delete", super_admin.DeleteNISN)
    superAuth.Post("/instansi/:id/backup", super_admin.BackupManual)
    superAuth.Get("/instansi/:id/backups", super_admin.RiwayatBackup)
    superAuth.Post("/instansi/restore/:backupId", super_admin.PulihkanBackup)
    superAuth.Get("/backup-setting", super_admin.AmbilSettingBackup)
superAuth.Put("/backup-setting", super_admin.SimpanSettingBackup)
superAuth.Get("/instansi/:id/backups/download/:backupId", super_admin.DownloadBackup)
superAuth.Post("/instansi/:id/import", super_admin.ImportBackup)


  
 app.Get("/api/me", middleware.Auth, controllers.GetMyProfile)
    app.Get("/api/logout", middleware.Auth, controllers.Logout)
    app.Post("/api/upload-foto", middleware.Auth, admin.UploadFotoProfile)
    app.Delete("/api/delete-foto", middleware.Auth, admin.DeleteFotoProfile)

    
    tenant := app.Group("/api/:slug", middleware.CheckTenant(config.DB))
    tenant.Post("/login-siswa", auth.LoginSiswa)

    portal := tenant.Group("/portal", middleware.AuthSiswa)
    portal.Get("/dashboard", siswa.GetSiswaDashboard)
    portal.Get("/setting/pesan", admin.GetSettingsPublic)
portal.Get("/setting/background", admin.GetBackgroundSettings)

    authTenant := tenant.Group("", middleware.Auth)
    authTenant.Get("/dashboard", admin.GetAdminDashboard)

    adminGroup := authTenant.Group("/admin", middleware.Admin)
      adminGroup.Post("/check-graduation", admin.CekKelayakanLulus)   
    adminGroup.Post("/execute-promotion", admin.ProsesLulusMasal)   
    adminGroup.Post("/setting/waktu-buka", admin.SetWaktuBukaPengumuman)
    adminGroup.Post("/siswa/import", admin.ImportSiswaExcel)
    adminGroup.Post("/upload-siswa", admin.UploadFotoSiswa)
    adminGroup.Delete("/siswa/:id/foto", admin.DeleteFotoSiswa)
    adminGroup.Post("/setting/tampilkan-logo", admin.UpdateTampilkanLogo)
      adminGroup.Get("/setting/pesan", admin.GetSettings)
    adminGroup.Put("/setting/pesan", admin.UpdateSettings)
    adminGroup.Get("/setting/background", admin.GetBackgroundSettings)
  adminGroup.Post("/setting/background/upload", admin.UploadBackground)
  adminGroup.Delete("/setting/background", admin.DeleteBackground)


     siswaGroup := adminGroup.Group("/siswa") 
    siswaGroup.Get("/", admin.GetAllSiswa)
    siswaGroup.Get("/search", admin.SearchSiswa)
      siswaGroup.Delete("/all", admin.DeleteAllSiswa)
    siswaGroup.Post("/", admin.CreateSiswa)
    siswaGroup.Get("/:id", admin.GetSiswaByID)
    siswaGroup.Put("/:id", admin.UpdateSiswa)
    siswaGroup.Delete("/:id", admin.DeleteSiswa) 
  

    adminGroup.Get("/mapel", admin.GetAllMapel)
    adminGroup.Post("/mapel", admin.CreateMapel)
    adminGroup.Put("/mapel/:id", admin.UpdateMapel)
    adminGroup.Delete("/mapel/:id", admin.DeleteMapel)

    adminGroup.Get("/template-skl", admin.AmbilTemplateSKL)
    adminGroup.Post("/template-skl", admin.SimpanTemplateSKL)
    adminGroup.Get("/cetak/:id", admin.AmbilDataSKLAdmin) 

    adminGroup.Get("/nilai/filters", admin.GetFilterSiswa)
    adminGroup.Get("/nilai/template", admin.DownloadTemplateNilai)
    adminGroup.Post("/nilai/import", admin.ImportNilaiExcel)
    adminGroup.Get("/nilai/leger", admin.GetLegerNilai)
    adminGroup.Get("/nilai/filter-options", admin.GetFilterSiswa)

}