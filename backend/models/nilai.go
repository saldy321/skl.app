package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nilai struct {
	gorm.Model
	InstansiID uuid.UUID `gorm:"type:char(36);index;not null" json:"instansi_id"`
	

	SiswaID    uuid.UUID `gorm:"type:char(36);uniqueIndex:idx_nilai_kolektif;not null" json:"siswa_id"`
	MapelID    uint      `gorm:"uniqueIndex:idx_nilai_kolektif;not null" json:"mapel_id"`
	
	NilaiAngka float64   `gorm:"type:decimal(5,2);default:0" json:"nilai_angka"`
	TahunAjaran string    `gorm:"type:varchar(10)" json:"tahun_ajaran"`
	Semester    int       `gorm:"type:int" json:"semester"`


	Mapel      Mapel     `gorm:"foreignKey:MapelID" json:"mapel"`
	Siswa      Siswa     `gorm:"foreignKey:SiswaID" json:"-"`
}