package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Siswa struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	InstansiID   uuid.UUID `gorm:"type:char(36);index;not null" json:"instansi_id"`
	FotoSiswa    string    `gorm:"type:varchar(255)" json:"foto_siswa"`
	NISN         string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nisn"`
	NamaSiswa    string    `gorm:"type:varchar(100);not null" json:"nama_siswa"`
	TempatLahir  string    `gorm:"type:varchar(100)" json:"tempat_lahir"`
	TanggalLahir string    `gorm:"type:varchar(50)" json:"tanggal_lahir"` 
	JenisKelamin string    `gorm:"type:enum('L','P')" json:"jenis_kelamin"`
	Kelas        string    `gorm:"type:varchar(50)" json:"kelas"`
	Jurusan      string    `gorm:"type:varchar(100)" json:"jurusan"`
	 TahunLulus  string `gorm:"type:varchar(4);index" json:"tahun_lulus"`
	NamaWali     string    `gorm:"type:varchar(100)" json:"nama_wali"`
	StatusLulus  bool      `gorm:"default:false" json:"status_lulus"`
	Nilai        []Nilai   `gorm:"foreignKey:SiswaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"nilai"`
}


func (s *Siswa) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}