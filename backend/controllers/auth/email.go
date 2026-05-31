package auth

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
)

func SendOTPEmail(toEmail string, code string) {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	jktTime := time.Now().In(time.FixedZone("WIB", 7*3600)).Format("15:04")
	msg := fmt.Sprintf("Subject: OTP SKL Digital [%s]\r\n\r\n"+
		"Kode OTP lu: %s\r\nIngat bro, kode ini cuma aktif 5 menit!",
		jktTime, code)

	auth := smtp.PlainAuth("", from, password, host)
	_ = smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}