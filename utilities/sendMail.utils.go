package utilities

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"

	"github.com/khanhvtn/netevent-go/models"
)

/* Send email from project email to others */
func SendMail(eventName string, listReceiver []*models.Participant) error {
	//send data
	from := os.Getenv("PROJECT_EMAIL")
	password := os.Getenv("PROJECT_EMAIL_PASSWORD")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	t, _ := template.ParseFiles(dir + "/templates/mail.template.html")

	//send email to all receivers
	for _, v := range listReceiver {
		var body bytes.Buffer
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: Invitation for %s event \n%s\n\n", eventName, mimeHeaders)))
		t.Execute(&body, struct {
			Name string
		}{
			Name: v.Name,
		})
		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{v.Email}, body.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}
