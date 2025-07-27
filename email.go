package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
)

type EmailAttachment struct {
	Path     string
	Filename string
}

type Document interface {
	GetNumber() string
	GetFilename() string // optional: helps decide the file name when saving/downloading
}

func SendEmailWithAttachment(to, subject, body string, doc Document, attachmentPath string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("EMAIL_SENDER_ADDRESS")
	senderName := os.Getenv("EMAIL_SENDER_NAME")

	// Ler o ficheiro PDF
	fileBytes, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("erro ao ler o ficheiro: %w", err)
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	boundary := writer.Boundary()
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", senderName, sender)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

	// Cabeçalhos principais
	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	// Parte do corpo (texto)
	partWriter, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	bodyEncoded := quotedprintable.NewWriter(partWriter)
	_, err = bodyEncoded.Write([]byte(body))
	if err != nil {
		return err
	}
	err = bodyEncoded.Close()
	if err != nil {
		return err
	}

	// Parte do anexo
	attachmentHeader := make(textproto.MIMEHeader)
	attachmentHeader.Set("Content-Type", "application/pdf")
	attachmentHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.GetFilename()))
	attachmentHeader.Set("Content-Transfer-Encoding", "base64")

	part, _ := writer.CreatePart(attachmentHeader)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(fileBytes)))
	base64.StdEncoding.Encode(encoded, fileBytes)
	_, err = part.Write(encoded)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Enviar o email
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{to}, append([]byte(msg.String()), b.Bytes()...))
}
