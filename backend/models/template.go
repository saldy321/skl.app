package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateSKL struct {
	gorm.Model

	InstansiID uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"instansi_id"`
	NamaSurat    string `json:"nama_surat"`
	NoSurat      string `json:"no_surat"`
	TanggalSurat string `json:"tanggal_surat"`
	NamaKepsek   string `json:"nama_kepsek"`
	NipKepsek    string `json:"nip_kepsek"`
	FileHeader string `gorm:"type:longtext" json:"file_header"`
	Stempel    string `gorm:"type:longtext" json:"stempel"`
	TtdKepsek  string `gorm:"type:longtext" json:"ttd_kepsek"`
	DasarSurat   string `gorm:"type:longtext" json:"dasar_surat"`
	IsiSurat     string `gorm:"type:longtext" json:"isi_surat"`
	PenutupSurat string `gorm:"type:longtext" json:"penutup_surat"`
	WidthStempel int `json:"width_stempel" gorm:"default:100"`
	WidthTtd     int `json:"width_ttd" gorm:"default:150"`
	MarginTop    int `json:"margin_top" gorm:"default:20"`
	PakaiStempel        bool    `json:"pakai_stempel" gorm:"default:true"`
	PakaiTtd            bool    `json:"pakai_ttd" gorm:"default:true"`
	PakaiTtdQrCode      bool    `json:"pakai_ttd_qrcode" gorm:"default:false"`
	PakaiNamaWali       bool    `json:"pakai_nama_wali" gorm:"default:false"`
	PakaiFoto           bool    `json:"pakai_foto" gorm:"default:true"`
	PakaiKelompokMapel  bool    `json:"pakai_kelompok_mapel" gorm:"default:true"`
	TampilkanNilaiAdmin bool    `json:"tampilkan_nilai_admin" gorm:"default:true"`
	TampilkanNilaiSiswa bool    `json:"tampilkan_nilai_siswa" gorm:"default:false"`
	MinimalKelulusan float64 `json:"minimal_kelulusan" gorm:"default:75"`
	MinimalNilai     float64 `json:"minimal_nilai"`
}