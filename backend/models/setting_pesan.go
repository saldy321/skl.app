package models
import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Setting_pesan struct {
    ID               uuid.UUID `gorm:"type:char(36);primaryKey"`
    InstansiID       uuid.UUID `gorm:"type:char(36);uniqueIndex;not null"`
    PesanLulus       string    `gorm:"type:text"`
    PesanTidakLulus  string    `gorm:"type:text"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}

func (s *Setting_pesan) BeforeCreate(tx *gorm.DB) error {
    s.ID = uuid.New()
    return nil
}