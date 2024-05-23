// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
			Files: []AttachData{{
				FileName: "test.txt",
				FileData: []byte("test"),
			}},
		})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTestLoginAuth(t *testing.T) {
	type testParams struct {
		host        string
		port        int
		tlsOn       bool
		tlsConfig   *tls.Config
		user        string
		password    string
		connTimeout time.Duration
	}

	tests := []struct {
		name          string
		params        testParams
		wantAuthError bool
		wantConnError bool
	}{
		{
			name: "connection error (bad host)",
			params: testParams{
				host:        "test",
				port:        25,
				tlsOn:       false,
				user:        "test",
				password:    "test",
				connTimeout: time.Second * 2,
			},
			wantConnError: true,
			wantAuthError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connErr, authErr := TestLoginAuth(
				tt.params.host, tt.params.port,
				tt.params.tlsOn, tt.params.tlsConfig,
				tt.params.user, tt.params.password,
				tt.params.connTimeout)
			if tt.wantConnError {
				require.Error(t, connErr)
			}
			if tt.wantAuthError {
				require.Error(t, authErr)
			}
		})
	}
}
