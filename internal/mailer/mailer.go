package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/Slightly-Techie/st-okr-api/config"
)



func sendEmail(to, subject, body string) error {
	smtpConfig := config.ENV

	auth := smtp.PlainAuth("", smtpConfig.SMTPUsername, smtpConfig.SMTPPassword, smtpConfig.SMTPHost)

	// Construct the email headers
	headers := make(map[string]string)
	headers["From"] = smtpConfig.SMTPUsername
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Construct the message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%s", smtpConfig.SMTPHost, smtpConfig.SMTPPort)

	// Create a custom TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(smtpConfig.SMTPUsername); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}



func SendWelcomeEmail(recipientEmail string, userName string) error {

	data := map[string]string{
		"userName": userName,
	}

	body, err := LoadTemplate("welcome", data)

	if err != nil {
		return err
	}

	return sendEmail(recipientEmail, "Welcome to OKR", body)
}