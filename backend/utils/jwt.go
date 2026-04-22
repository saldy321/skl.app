package utils

import (
	"os"
	"time"
	 
	"github.com/golang-jwt/jwt/v5"

)
func GenerateToken(id string, role string, tingkat string, instansi string, instansiID interface{}, slug string) (string, error) {
     claims := jwt.MapClaims{
		"id":              id,
		"role":            role,
		"tingkat_sekolah": tingkat,
		"nama_instansi":   instansi,
		"instansi_id":     instansiID,
		"slug":				slug,
		"exp":             time.Now().Add(24 * time.Hour).Unix(),
	 }

	 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	 return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
 }