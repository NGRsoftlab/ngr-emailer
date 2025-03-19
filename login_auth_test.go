package emailer

import (
	"crypto/tls"
	"net/smtp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoginAuthMethods(t *testing.T) {
	auth := LoginAuth("user", "pass")
	require.NotNil(t, auth)

	proto, resp, err := auth.Start(&smtp.ServerInfo{})
	require.NoError(t, err)
	require.Equal(t, "LOGIN", proto)
	require.Empty(t, resp)

	tests := []struct {
		name         string
		fromServer   []byte
		more         bool
		expectedResp []byte
		expectError  bool
	}{
		{
			name:         "Username prompt",
			fromServer:   []byte("Username:"),
			more:         true,
			expectedResp: []byte("user"),
			expectError:  false,
		},
		{
			name:         "Password prompt",
			fromServer:   []byte("Password:"),
			more:         true,
			expectedResp: []byte("pass"),
			expectError:  false,
		},
		{
			name:         "Unknown prompt",
			fromServer:   []byte("Unknown:"),
			more:         true,
			expectedResp: nil,
			expectError:  true,
		},
		{
			name:         "No more data",
			fromServer:   nil,
			more:         false,
			expectedResp: nil,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				resp, err := auth.Next(tt.fromServer, tt.more)

				if tt.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.expectedResp, resp)
			},
		)
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
		t.Run(
			tt.name, func(t *testing.T) {
				connErr, authErr := TestLoginAuth(
					tt.params.host, tt.params.port,
					tt.params.tlsOn, tt.params.tlsConfig,
					tt.params.user, tt.params.password,
					tt.params.connTimeout,
				)
				if tt.wantConnError {
					require.Error(t, connErr)
				}
				if tt.wantAuthError {
					require.Error(t, authErr)
				}
			},
		)
	}
}
