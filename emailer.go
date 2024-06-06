// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"strings"
	"time"
)

/////////////////////////////////////////////

// Sender struct - sender for sending smtp packs
type Sender struct {
	Login      string       // user login
	Email      string       // user email address
	Password   string       // user password
	ServerSMTP string       // smtp server string
	client     *smtp.Client // smtp client pointer
	message    []byte       // message text
	to         []string     // receivers of email
}

// NewSender creating new *Sender obj
func NewSender(login, password, email, server string) *Sender {
	auth := Sender{
		Login:      login,
		Email:      email,
		Password:   password,
		ServerSMTP: server}
	return &auth
}

/////////////////////////////////////////////

// Send send smtp pack (mail) with login auth
func (s *Sender) Send() error {
	err := smtp.SendMail(s.ServerSMTP,
		LoginAuth(s.Login, s.Password),
		s.Login, s.to, s.message)

	if err != nil {
		logger.Errorf("send error: %s", err.Error())
		return err
	}
	return nil
}

// SendWithAuth send smtp pack (mail) with custom auth
func (s *Sender) SendWithAuth(auth smtp.Auth) error {
	err := smtp.SendMail(s.ServerSMTP,
		auth,
		s.Login, s.to, s.message)

	if err != nil {
		logger.Errorf("send error: %s", err.Error())
		return err
	}
	return nil
}

/////////////////////////////////////////////

// NewMessage creating new email message
func (s *Sender) NewMessage(params *MessageParams) error {
	logger.Infof("files: %d", len(params.Files))

	attachments, err := attachFile(params.Files)
	if err != nil {
		logger.Error(err)
		return err
	}

	withAttachments := len(attachments) > 0
	var headers = make(map[string]string)
	headers["From"] = s.Email
	headers["To"] = strings.Join(params.Recipients, ";")
	headers["Subject"] = params.Topic
	headers["MIME-Version"] = "1.0"
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	var buf = bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	if withAttachments {
		buf.WriteString(fmt.Sprintf(`Content-Type: %s; boundary="%s"`, MultipartMixedContentType, boundary))
		buf.WriteString("\r\n\r\n")
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}
	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=%s\r\n", params.ContentType, params.Charset))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("\r\n" + params.Body)
	if withAttachments {
		for k, v := range attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", VndSheetContentType))
			buf.WriteString(fmt.Sprintf("Content-Transfer-Encoding: %s\r\n", Base64Charset))
			buf.WriteString("MIME-Version: 1.0\r\n")
			buf.WriteString(fmt.Sprintf(`Content-Disposition: attachment; filename="%s"`, k))
			buf.WriteString("\r\n\r\n")

			var b = make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
		}
		buf.WriteString("--")
	}
	s.to = params.Recipients
	s.message = buf.Bytes()

	return nil
}
