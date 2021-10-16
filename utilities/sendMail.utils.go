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

/* Send email from project email to others */
func SendActivateAccountMail(newUser *models.User) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	sender := models.NewSender()
	//send email to all receivers
	m := models.NewMail()
	m.To = append(m.To, newUser.Email)
	m.Subject = "Activate Account"
	t, err := template.ParseFiles(dir + "/templates/activateAccountMail.template.html")
	if err != nil {
		return err
	}
	m.MailTemplate = &models.MailTemplate{Template: t, Data: struct {
		Email         string
		ActivatedLink string
	}{
		ActivatedLink: fmt.Sprintf("https://net-event.herokuapp.com/activate/%s", newUser.ID.Hex()),
		Email:         newUser.Email,
	}}
	if err := sender.Send(m); err != nil {
		return err
	}
	return nil
}
