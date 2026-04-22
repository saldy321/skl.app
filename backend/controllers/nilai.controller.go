package controllers

import (
	"bytes"
	"strconv"
	"strings"
	"skl-bakcend/config"
	"skl-bakcend/models"
	"skl-bakcend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetFilterSiswa(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	var kelas, jurusan []string

	config.DB.Model(&models.Siswa{}).Where("instansi_id = ?", instansiID).Distinct().Pluck("kelas", &kelas)
	config.DB.Model(&models.Siswa{}).Where("instansi_id = ?", instansiID).Distinct().Pluck("jurusan", &jurusan)

	return c.JSON(fiber.Map{"kelas": kelas, "jurusan": jurusan})
}

func DownloadTemplateNilai(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	filterKelas := c.Query("kelas")
	filterJurusan := c.Query("jurusan")

	var mapels []models.Mapel
	config.DB.Where("instansi_id = ?", instansiID).Order("id asc").Find(&mapels)

	var siswas []models.Siswa
	db := config.DB.Where("instansi_id = ?", instansiID)
	if filterKelas != "" { db = db.Where("kelas = ?", filterKelas) }
	if filterJurusan != "" { db = db.Where("jurusan = ?", filterJurusan) }
	db.Order("nama_siswa asc").Find(&siswas)

	f := excelize.NewFile()
	sheet := "Sheet1"
	
	f.SetCellValue(sheet, "A1", "ID_SISWA")
	f.SetCellValue(sheet, "B1", "NO")
	f.SetCellValue(sheet, "C1", "NISN")
	f.SetCellValue(sheet, "D1", "NAMA SISWA")

	for i, m := range mapels {
		col, _ := excelize.ColumnNumberToName(i + 5)
		f.SetCellValue(sheet, col+"1", m.NamaMapel)
	}

	for i, s := range siswas {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.ID.String())
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), i+1)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), s.NISN)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), s.NamaSiswa)
	}

	f.SetColVisible(sheet, "A", false)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=Template_Nilai.xlsx")

	var buf bytes.Buffer
	f.Write(&buf)
	return c.Send(buf.Bytes())
}

func ImportNilaiExcel(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	file, err := c.FormFile("file_excel")
	if err != nil { return c.Status(400).JSON(fiber.Map{"message": "File kosong"}) }

	src, _ := file.Open()
	defer src.Close()
	f, _ := excelize.OpenReader(src)
	rows, _ := f.GetRows("Sheet1")
	
	if len(rows) < 2 { return c.Status(400).JSON(fiber.Map{"message": "Data kosong"}) }
	
	header := rows[0]
	var allMapels []models.Mapel
	config.DB.Where("instansi_id = ?", instansiID).Find(&allMapels)

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		for i, row := range rows {
			if i == 0 || len(row) < 4 || row[0] == "" { continue }
			
			siswaID, errParse := uuid.Parse(row[0])
			if errParse != nil { continue }

			for colIdx := 4; colIdx < len(row); colIdx++ {
				if colIdx >= len(header) { break }
				
				valStr := strings.TrimSpace(row[colIdx])
				if valStr == "" { continue }
				
				nilaiAngka, _ := strconv.ParseFloat(valStr, 64)
				namaMapelHeader := strings.TrimSpace(strings.ToLower(header[colIdx]))

				var mID uint
				for _, m := range allMapels {
					if strings.TrimSpace(strings.ToLower(m.NamaMapel)) == namaMapelHeader {
						mID = m.ID
						break
					}
				}

				if mID != 0 {
					tx.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "siswa_id"}, {Name: "mapel_id"}},
						DoUpdates: clause.AssignmentColumns([]string{"nilai_angka", "updated_at"}),
					}).Create(&models.Nilai{
						InstansiID:  instansiID,
						SiswaID:     siswaID,
						MapelID:     mID,
						NilaiAngka:  nilaiAngka,
						TahunAjaran: "2025/2026",
						Semester:    2,
					})
				}
			}
		}
		return nil
	})

	if err != nil { return c.Status(500).JSON(fiber.Map{"message": err.Error()}) }
	return c.JSON(fiber.Map{"message": "Selesai! Data berhasil disinkronkan."})
}

func GetLegerNilai(c *fiber.Ctx) error {
	instansiID := c.Locals(utils.KeyInstansiID).(uuid.UUID)
	filterKelas := c.Query("kelas")
	filterJurusan := c.Query("jurusan")
	
	// <--- TAMBAHKAN INI: Ambil parameter tahun_lulus dari Frontend
	filterTahun := c.Query("tahun_lulus")

	var mapels []models.Mapel
	config.DB.Where("instansi_id = ?", instansiID).Order("id asc").Find(&mapels)

	var siswas []models.Siswa
	db := config.DB.Where("instansi_id = ?", instansiID)
	
	if filterKelas != "" { db = db.Where("kelas = ?", filterKelas) }
	if filterJurusan != "" { db = db.Where("jurusan = ?", filterJurusan) }
	
	// <--- TAMBAHKAN INI: Filter siswa berdasarkan tahun lulus jika ada paramnya
	if filterTahun != "" {
		db = db.Where("tahun_lulus = ?", filterTahun)
	}
	
	db.Order("nama_siswa asc").Find(&siswas)

	var legerData []map[string]interface{}
	
	for _, s := range siswas {
		row := map[string]interface{}{
			"nama": s.NamaSiswa,
			"nisn": s.NISN,
		}
		
		var nilais []models.Nilai
		config.DB.Where("siswa_id = ?", s.ID).Find(&nilais)
		
		var totalNilai float64
		var jmlMapel int
		
		for _, m := range mapels {
			row[m.NamaMapel] = "-"
			for _, n := range nilais {
				if n.MapelID == m.ID {
					row[m.NamaMapel] = n.NilaiAngka
					totalNilai += n.NilaiAngka
					jmlMapel++
					break
				}
			}
		}

		if jmlMapel > 0 {
			row["rata_rata"] = totalNilai / float64(jmlMapel)
		} else {
			row["rata_rata"] = 0
		}
		
		legerData = append(legerData, row)
	}

	return c.JSON(fiber.Map{
		"mapels": mapels,
		"data":   legerData,
	})
}