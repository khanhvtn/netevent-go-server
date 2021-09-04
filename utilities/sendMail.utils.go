package utilities

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"

	"github.com/khanhvtn/netevent-go/models"
	"github.com/skip2/go-qrcode"
)

/* Send email from project email to others */
func SendMail(eventName, participantId, eventId string, listReceiver []*models.Participant) error {
	// //send data
	// from := os.Getenv("PROJECT_EMAIL")
	// password := os.Getenv("PROJECT_EMAIL_PASSWORD")

	// // smtp server configuration.
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// // Authentication.
	// auth := smtp.PlainAuth("", from, password, smtpHost)
	// dir, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }
	// t, _ := template.ParseFiles(dir + "/templates/mail.template.html")

	// //send email to all receivers
	// for _, v := range listReceiver {
	// 	var body bytes.Buffer
	// 	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// 	body.WriteString(fmt.Sprintf("Subject: Invitation for %s event \n%s\n\n", eventName, mimeHeaders))
	// 	t.Execute(&body, struct {
	// 		Name string
	// 	}{
	// 		Name: v.Name,
	// 	})
	// 	// Sending email.
	// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{v.Email}, body.Bytes())
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	sender := models.NewSender()
	//encrypt event id and participant id
	valueQrCode := struct {
		EventID       string `json:"eventId"`
		ParticipantID string `json:"participantId"`
	}{
		EventID:       eventId,
		ParticipantID: participantId,
	}
	valueQrCodeJson, err := json.Marshal(valueQrCode)
	if err != nil {
		return err
	}
	encryptedValueQrCodeJson, err := Encrypt(valueQrCodeJson)
	if err != nil {
		return err
	}
	//generate qrcode
	qrCode, err := qrcode.Encode(string(encryptedValueQrCodeJson), qrcode.Medium, 256)
	if err != nil {
		return err
	}
	//send email to all receivers
	for _, v := range listReceiver {
		m := models.NewMail()
		m.To = append(m.To, v.Email)
		m.Subject = fmt.Sprintf("Invitation for %s event", eventName)
		t, err := template.ParseFiles(dir + "/templates/mail.template.html")
		if err != nil {
			return err
		}
		m.MailTemplate = &models.MailTemplate{Template: t, Data: struct{ Name string }{
			Name: v.Name,
		}}

		m.AttachByteFile("qrCode.png", qrCode)
		if err := sender.Send(m); err != nil {
			return err
		}
	}
	return nil
}
