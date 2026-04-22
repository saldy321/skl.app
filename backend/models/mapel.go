package models

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Mapel struct {
	
	ID         uint           `gorm:"primarykey" json:"id"` 
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	InstansiID uuid.UUID `gorm:"type:char(36);index;not null" json:"instansi_id"`
	NamaMapel  string    `gorm:"type:varchar(100);not null" json:"nama_mapel"` 
	Kelompok   string    `gorm:"type:varchar(10)" json:"kelompok"`
	Jurusan    string    `gorm:"type:varchar(100)" json:"jurusan"`
	Slug       string    `gorm:"type:varchar(150);uniqueIndex" json:"slug"`
}

func (m *Mapel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Slug == "" {
		baseSlug := strings.ToLower(strings.ReplaceAll(m.NamaMapel, " ", "-"))
		m.Slug = fmt.Sprintf("%s-%s", baseSlug, uuid.New().String()[:5])
	}
	return nil
}