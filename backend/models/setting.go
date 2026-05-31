package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type BackupSetting struct {
    ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
    JamBackup   int       `gorm:"default:2" json:"jam_backup"`
    MenitBackup int       `gorm:"default:0" json:"menit_backup"`
    RetensiHari int       `gorm:"default:60" json:"retensi_hari"`
    AutoBackup  bool      `gorm:"default:true" json:"auto_backup"`
    UpdatedAt   time.Time `json:"updated_at"`
}


func (BackupSetting) AfterAutoMigrate(tx *gorm.DB) error {
    var count int64
    tx.Model(&BackupSetting{}).Count(&count)
    if count == 0 {
        return tx.Create(&BackupSetting{
            ID:          uuid.New(),
            JamBackup:   2,
            MenitBackup: 0,  
            RetensiHari: 60,
            AutoBackup:  true,
        }).Error
    }
    return nil
}