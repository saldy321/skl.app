package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Instansi struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	NamaInstansi  string         `gorm:"type:varchar(255);not null" json:"nama_instansi"`
	KodeInstansi  string         `gorm:"type:varchar(50);unique;not null" json:"kode_instansi"`
	Alamat        string         `gorm:"type:text" json:"alamat"`
	TingkatSekolah string        `gorm:"type:varchar(50);not null" json:"tingkat_sekolah"` // SMK/SMA/SMP
	Slug 			string		 `gorm:"type:varchar(255);unique;not null" json:"slug"`
	 LogoInstansi  string         `gorm:"type:varchar(255)" json:"logo_instansi"` // ← TAMBAHKAN
	  TampilkanLogo bool           `gorm:"default:true" json:"tampilkan_logo"`
	 WaktuBukaPengumuman *time.Time `gorm:"type:datetime" json:"waktu_buka_pengumuman"` 
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *Instansi) BeforeCreate(tx *gorm.DB) (err error) {
	
	i.ID = uuid.New()

	
	if i.Slug == "" {
       
        cleanSlug := strings.ToLower(i.NamaInstansi)
        cleanSlug = strings.ReplaceAll(cleanSlug, " ", "-")
        
        i.Slug = cleanSlug
    }
	
	return nil
}