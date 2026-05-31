package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type SettingBackground struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
    InstansiID  uuid.UUID `gorm:"type:char(36);uniqueIndex;not null"`
    Background  string    `gorm:"type:varchar(255)"` // path file background
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (s *SettingBackground) BeforeCreate(tx *gorm.DB) error {
    s.ID = uuid.New()
    return nil
}