package routes

import (
    "skl-bakcend/config"
    "skl-bakcend/controllers"
    "skl-bakcend/middleware"

    "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
    // ========== PUBLIC ROUTES (NO AUTH) ==========
    app.Post("/api/login", controllers.Login)
    app.Post("/api/verify-otp", controllers.VerifyOTP)
    app.Post("/api/restore-super-admin", controllers.RestoreSuperAdmin)
    app.Get("/api/check-instansi/:slug", controllers.CheckInstansi)

    // ========== SUPER ADMIN GROUP (HARUS DI ATAS TENANT!) ==========
    superAuth := app.Group("/api/super", middleware.Auth, middleware.IsSuperAdmin)
    superAuth.Get("/dashboard", controllers.GetSuperDashboard)
    superAuth.Post("/instansi", controllers.RegisterFullSchool)
    superAuth.Get("/instansi", controllers.GetAllInstansi)
    superAuth.Get("/instansi/:id", controllers.GetInstansiByID)
    superAuth.Put("/instansi/:id", controllers.UpdateInstansi)
    superAuth.Delete("/instansi/:id", controllers.DeleteInstansi)
    superAuth.Post("/admin-register", controllers.RegisterAdminSekolah)
    superAuth.Get("/admin-list", controllers.GetAllAdminSekolah)
    superAuth.Post("/instansi/:instansi_id/reset-password", controllers.ResetPasswordAdminByInstansi)
    superAuth.Put("/admin/:id/email", controllers.UpdateAdminEmail)
    superAuth.Post("/impersonate/:instansi_id", controllers.ImpersonateAsAdmin)
    superAuth.Get("/nisn/search", controllers.SearchNISN)
    superAuth.Post("/nisn/force-delete", controllers.ForceDeleteNISN)

    // ========== PROTECTED ROUTES (AUTH REQUIRED, NO TENANT) ==========
 app.Get("/api/me", middleware.Auth, controllers.GetMyProfile)
    app.Get("/api/logout", middleware.Auth, controllers.Logout)
    app.Post("/api/upload-foto", middleware.Auth, controllers.UploadFotoProfile)
    app.Delete("/api/delete-foto", middleware.Auth, controllers.DeleteFotoProfile)

    // ========== TENANT GROUP (HARUS DI BAWAH SUPER ADMIN!) ==========
    tenant := app.Group("/api/:slug", middleware.CheckTenant(config.DB))
    tenant.Post("/login-siswa", controllers.LoginSiswa)

    portal := tenant.Group("/portal", middleware.AuthSiswa)
    portal.Get("/dashboard", controllers.GetSiswaDashboard)

    authTenant := tenant.Group("", middleware.Auth)
    authTenant.Get("/dashboard", controllers.GetAdminDashboard)

    admin := authTenant.Group("/admin", middleware.Admin)
    admin.Post("/check-graduation", controllers.CheckGraduationEligibility)
    admin.Post("/execute-promotion", controllers.ExecuteMassPromotion)
    admin.Post("/setting/waktu-buka", controllers.SetWaktuBukaPengumuman)
    admin.Post("/siswa/import", controllers.ImportSiswaExcel)
    admin.Post("/upload-siswa", controllers.UploadFotoSiswa)
    admin.Delete("/siswa/:id/foto", controllers.DeleteFotoSiswa)
    admin.Post("/setting/tampilkan-logo", controllers.UpdateTampilkanLogo)

    siswa := admin.Group("/siswa")
    siswa.Get("/", controllers.GetAllSiswa)
    siswa.Post("/", controllers.CreateSiswa)
    siswa.Get("/:id", controllers.GetSiswaByID)
    siswa.Put("/:id", controllers.UpdateSiswa)
    siswa.Delete("/:id", controllers.DeleteSiswa)

    admin.Get("/mapel", controllers.GetAllMapel)
    admin.Post("/mapel", controllers.CreateMapel)
    admin.Put("/mapel/:id", controllers.UpdateMapel)
    admin.Delete("/mapel/:id", controllers.DeleteMapel)

    admin.Get("/template-skl", controllers.GetTemplateSKL)
    admin.Post("/template-skl", controllers.SaveTemplateSKL)
    admin.Get("/cetak/:id", controllers.GetSKLDataForAdmin)

    admin.Get("/nilai/filters", controllers.GetFilterSiswa)
    admin.Get("/nilai/template", controllers.DownloadTemplateNilai)
    admin.Post("/nilai/import", controllers.ImportNilaiExcel)
    admin.Get("/nilai/leger", controllers.GetLegerNilai)
    admin.Get("/nilai/filter-options", controllers.GetFilterSiswa)

}