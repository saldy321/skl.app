package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type SuperAdmin struct {
    ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
    Email     string         `gorm:"unique;not null" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    Role      string         `gorm:"default:super_admin" json:"role"`
    Otp          string    `gorm:"type:varchar(10)" json:"-"`
    OtpCreatedAt *time.Time `gorm:"type:timestamp;null" json:"-"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *SuperAdmin) BeforeCreate(tx *gorm.DB) (err error) {
    s.ID = uuid.New()
    return nil
}