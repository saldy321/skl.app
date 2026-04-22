package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Email        string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password     string         `gorm:"type:varchar(255);not null" json:"-"` 
	Role         string         `gorm:"type:varchar(50);default:admin" json:"role"` 
	InstansiID   uuid.UUID      `gorm:"type:char(36);not null" json:"instansi_id"`
	Instansi     Instansi       `gorm:"foreignKey:InstansiID" json:"instansi"` 
	    FotoProfile string     `gorm:"type:varchar(255)" json:"foto_profile"`
	Otp           string        `gorm:"type:varchar(10)" json:"-"`
	OtpCreatedAt  *time.Time     `gorm:"type:timestamp;null" json:"-"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return nil
}