// Copyright 2020 NGR Softlab
//
package emailer

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/NGRsoftlab/ngr-logging"
	"mime/multipart"
	"net/smtp"
	"strings"
	"time"
)

// Struct for auth.
type loginAuth struct {
	username, password string
}

// Struct for attachments.
type AttachData struct {
	FileName string
	FileData []byte
}

// Auth with login and password.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Login auth start.
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next loginAuth.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}

/////////////////////////////////////////////

// Struct - sender for sending smtp packs.
type Sender struct {
	Login      string       // user login
	Email      string       // user email address
	Password   string       // user password
	ServerSMTP string       // smtp server string
	client     *smtp.Client // smtp client pointer
	message    []byte       // message text
	to         []string     // receivers of email
}

// Send smtp pack (mail).
func (s *Sender) Send() error {
	err := smtp.SendMail(s.ServerSMTP,
		LoginAuth(s.Login, s.Password),
		s.Login, s.to, s.message)

	if err != nil {
		logging.Logger.Error("SEND ERROR: ", err)
		return err
	}
	return nil
}

/////////////////////////////////////////////

// Creating new *Sender obj.
func NewSender(login, password, email, server string) *Sender {
	auth := Sender{
		Login:      login,
		Email:      email,
		Password:   password,
		ServerSMTP: server}
	return &auth
}

// Creating new email message.
func (s *Sender) NewMessage(subject string, to []string, body string, files []AttachData) error {
	logging.Logger.Info("files: ", len(files))

	attachments, err := attachFile(files)
	if err != nil {
		logging.Logger.Error(err)
		return err
	}

	withAttachments := len(attachments) > 0
	var headers = make(map[string]string)
	headers["From"] = s.Email
	headers["To"] = strings.Join(to, ";")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	var buf = bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	if withAttachments {
		buf.WriteString(fmt.Sprintf(`Content-Type: multipart/mixed; boundary="%s"`, boundary))
		buf.WriteString("\r\n\r\n")
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}
	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("\r\n" + body)
	if withAttachments {
		for k, v := range attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString("Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet\r\n")
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString("MIME-Version: 1.0\r\n")
			buf.WriteString(fmt.Sprintf(`Content-Disposition: attachment; filename="%s"`, k))
			buf.WriteString("\r\n\r\n")

			var b = make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
		}
		buf.WriteString("--")
	}
	s.to = to
	s.message = buf.Bytes()

	return nil
}

// Attaching file to mail. Returns attachments map: "filename": filedata.
func attachFile(files []AttachData) (map[string][]byte, error) {
	var attachments = make(map[string][]byte)
	for _, f := range files {
		attachments[f.FileName] = f.FileData
	}
	return attachments, nil
}
