// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSender(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		email    string
		server   string
	}{
		{
			name:     "valid sender params",
			login:    "user",
			password: "pass",
			email:    "test@test.com",
			server:   "smtp.test.com:587",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				sender := NewSender(tt.login, tt.password, tt.email, tt.server)

				require.Equal(t, tt.login, sender.Login)
				require.Equal(t, tt.password, sender.Password)
				require.Equal(t, tt.email, sender.Email)
				require.Equal(t, tt.server, sender.ServerSMTP)
			},
		)
	}
}

func TestNewMessage(t *testing.T) {
	s := NewSender(
		"",
		"",
		"",
		fmt.Sprintf("%v:%v", "", "587"),
	)

	body := `<h3>Test</h3>
Период: с ___ по ___</br>
</br>
1. Test1
2. Test2</br>
`

	err := s.NewMessage(
		&MessageParams{
			Topic:       "topic",
			ContentType: HtmlContentType,
			Charset:     Utf8Charset,
			Recipients:  []string{""},
			Body:        body,
			Files: []AttachData{
				{
					FileName: "test.txt",
					FileData: []byte("test"),
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewMessage_Table(t *testing.T) {
	sender := NewSender("user", "pass", "test@test.com", "smtp.test.com:587")

	tests := []struct {
		name   string
		params *MessageParams
	}{
		{
			name: "valid message without attachments",
			params: &MessageParams{
				Topic:       "test",
				ContentType: "text/plain",
				Charset:     "utf-8",
				Recipients:  []string{"recipient@test.com"},
				Body:        "Hello World",
			},
		},
		{
			name: "valid message with attachments",
			params: &MessageParams{
				Topic:       "test attach",
				ContentType: "text/plain",
				Charset:     "utf-8",
				Recipients:  []string{"test@test.com"},
				Files: []AttachData{
					{
						FileName: "test.txt",
						FileData: []byte("test"),
					},
					{
						FileName: "test1.txt",
						FileData: []byte("test1"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := sender.NewMessage(tt.params)
				require.NoError(t, err)
			},
		)
	}
}
