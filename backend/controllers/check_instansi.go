package controllers

import (
	"time"

	"skl-bakcend/config"
	"skl-bakcend/models"

	"github.com/gofiber/fiber/v2"
)

func CheckInstansi(c *fiber.Ctx) error {
    slug := c.Params("slug")
    var instansi models.Instansi
    
   
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
            "logo_instansi":   instansi.LogoInstansi,   
            "tampilkan_logo":  instansi.TampilkanLogo,  
        },
    })
} 