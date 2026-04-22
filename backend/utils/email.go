package utils

import (
	"os"
	"gopkg.in/gomail.v2"
	"fmt"
	"strconv"
)

func SendCredentialEmail(toEmail, namaInstansi, slug, password string) error {
	// Ambil config dari .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	// Konversi port ke integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", fmt.Sprintf("Akun Admin Baru - %s", namaInstansi))

	// Desain Email Simple & Clean
	body := fmt.Sprintf(`
	<div style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; max-width: 600px; margin: 0 auto; padding: 30px; border: 1px solid #e2e8f0; border-radius: 12px; background-color: #ffffff;">
		<div style="text-align: center; margin-bottom: 20px;">
			<h2 style="color: #1e293b; margin: 0;">Selamat Datang di Sistem SKL!</h2>
			<p style="color: #64748b; font-size: 14px;">Admin Sekolah <strong>%s</strong></p>
		</div>
		
		<p style="color: #334155; line-height: 1.6;">Halo Admin,</p>
		<p style="color: #334155; line-height: 1.6;">Akun Anda telah berhasil dibuat oleh Super Admin. Berikut adalah kredensial login Anda:</p>
		
		<div style="background-color: #f1f5f9; padding: 20px; border-radius: 8px; margin: 25px 0; border-left: 4px solid #4f46e5;">
			<p style="margin: 8px 0; font-size: 14px;"><strong>🌐 Link Login:</strong> <a href="http://localhost:5173/%s/login" style="color: #4f46e5; text-decoration: none;">Klik Disini untuk Login</a></p>
			<p style="margin: 8px 0; font-size: 14px;"><strong>📧 Email:</strong> %s</p>
			<p style="margin: 8px 0; font-size: 14px;"><strong>🔑 Password:</strong> <span style="background: #e2e8f0; padding: 4px 8px; border-radius: 4px; font-family: monospace; color: #0f172a;">%s</span></p>
		</div>

		<p style="font-size: 12px; color: #94a3b8; margin-top: 30px; border-top: 1px solid #e2e8f0; padding-top: 10px;">
			⚠️ Harap segera ganti password Anda setelah login pertama kali demi keamanan data sekolah.
		</p>
	</div>
	`, namaInstansi, slug, toEmail, password)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}