package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, otp string) error {
	from := os.Getenv("MAIL_USERNAME")
	fromName := os.Getenv("MAIL_FROM_NAME")
	pass := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")

	// logoURL := os.Getenv("APP_URL") + "/logo.png"

	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8" />
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #ffffff;
				color: #333;
				padding: 20px;
			}
			.container {
				background: #fff;
				border-radius: 10px;
				padding: 20px;
				max-width: 500px;
				margin: auto;
				box-shadow: 0 4px 6px rgba(0,0,0,0.1);
				text-align: center;
				border-top: 6px solid #57A32E;
			}
			.logo {
				width: 100px;
				margin-bottom: 20px;
			}
			h2 {
				color: #57A32E;
			}
			.otp {
				font-size: 28px;
				font-weight: bold;
				color: #57A32E; /* OTP warna base */
				margin: 20px 0;
			}
			.btn {
				display: inline-block;
				padding: 12px 24px;
				background: #57A32E;
				color: #fff !important;
				font-weight: bold;
				text-decoration: none;
				border-radius: 6px;
				margin-top: 20px;
			}
			.footer {
				margin-top: 30px;
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<img src="./public/logo.png" class="logo" alt="Logo">
			<h2>Kode OTP Reset Password</h2>
			<p>Gunakan kode OTP berikut untuk reset password akunmu:</p>
			<div class="otp">%s</div>
			<a href="#" class="btn">Reset Password</a>
			<p>Kode berlaku selama <b>10 menit</b>.</p>
			<div class="footer">
				&copy; %d %s. All rights reserved.
			</div>
		</div>
	</body>
	</html>
`, otp, 2025, fromName)

	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		fmt.Sprintf("From: %s <%s>\r\n", fromName, from) +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", from, pass, host)
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
}
