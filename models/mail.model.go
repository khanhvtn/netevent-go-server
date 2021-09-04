package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type Sender struct {
	auth smtp.Auth
}

type Mail struct {
	To           []string
	CC           []string
	BCC          []string
	Subject      string
	Body         string
	Attachments  map[string][]byte
	MailTemplate *MailTemplate
}

type MailTemplate struct {
	Template *template.Template
	Data     interface{}
}

func NewSender() *Sender {
	auth := smtp.PlainAuth("", os.Getenv("PROJECT_EMAIL"), os.Getenv("PROJECT_EMAIL_PASSWORD"), os.Getenv("PROJECT_EMAIL_HOST"))
	return &Sender{auth}
}

func (s *Sender) Send(m *Mail) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("PROJECT_EMAIL_HOST"), os.Getenv("PROJECT_EMAIL_PORT")), s.auth, os.Getenv("PROJECT_EMAIL"), m.To, m.ToBytes())
}

func NewMail() *Mail {
	return &Mail{
		Attachments:  make(map[string][]byte),
		MailTemplate: nil,
	}
}

func (m *Mail) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}
func (m *Mail) AttachByteFile(fileName string, fileByte []byte) error {
	m.Attachments[fileName] = fileByte
	return nil
}

func (m *Mail) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
	buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	if m.MailTemplate != nil {
		buf.WriteString("Content-Type: text/html; charset=utf-8\n")
		m.MailTemplate.Template.Execute(buf, m.MailTemplate.Data)
	}
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
		buf.WriteString("--")
	}
	return buf.Bytes()
}
