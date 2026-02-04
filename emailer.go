// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"strings"
	"time"
)

/////////////////////////////////////////////

// Sender struct - sender for sending smtp packs
type Sender struct {
	Auth   smtp.Auth
	tlsCfg *tls.Config

	Login    string // user login
	Password string // user password

	ServerSMTP string // smtp server string (full)
	ServerAddr string // smtp server addr string

	Email   string   // user email address (from)
	to      []string // receivers of email (to)
	message []byte   // message text
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

// NewTLSSender creating new *Sender obj
func NewTLSSender(auth smtp.Auth, tlsCfg *tls.Config, email, hostAddress, fullServerAddress string) (*Sender, error) {
	sender := Sender{
		Auth:       auth,
		tlsCfg:     tlsCfg,
		Email:      email,
		ServerSMTP: fullServerAddress,
		ServerAddr: hostAddress,
	}

	return &sender, nil
}

/////////////////////////////////////////////

// SendViaClient send smtp pack (mail)
func (s *Sender) SendViaClient() error {
	if err := validateLine(s.Email); err != nil {
		return err
	}
	for _, recp := range s.to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	client, err := smtp.Dial(s.ServerSMTP)
	if err != nil {
		return err
	}

	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok && s.tlsCfg != nil {
		config := s.tlsCfg
		if err = client.StartTLS(config); err != nil {
			return err
		}
	}

	if err = client.Auth(s.Auth); err != nil {
		return err
	}

	if err = client.Mail(s.Email); err != nil {
		return err
	}

	for _, addr := range s.to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(s.message)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// Send send smtp pack (mail) with login auth
func (s *Sender) Send() error {
	err := smtp.SendMail(s.ServerSMTP,
		LoginAuth(s.Login, s.Password),
		s.Email, s.to, s.message)

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
		s.Email, s.to, s.message)

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

/////////////////////////////////////////////

// validateLine checks to see if a line has CR or LF as per RFC 5321.
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("emailer: A line must not contain CR or LF")
	}
	return nil
}
