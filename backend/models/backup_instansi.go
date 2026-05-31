package models

import (
	"time"
	"github.com/google/uuid"
)

type BackupInstansi struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	InstansiID  uuid.UUID `gorm:"type:char(36);index;not null" json:"instansi_id"`
	NamaBackup  string    `gorm:"type:varchar(255)" json:"nama_backup"`
	FilePath 	string 	`gorm:"type:varchar(500)" json:"file_path"`
	DataJSON    string    `gorm:"type:longtext" json:"data_json"`
	JumlahSiswa int       `json:"jumlah_siswa"`
	JumlahNilai int       `json:"jumlah_nilai"`
	JumlahMapel int       `json:"jumlah_mapel"`
	CreatedAt   time.Time `json:"created_at"`
}