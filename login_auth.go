// Copyright 2020-2024 NGR Softlab
package emailer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	logging "github.com/NGRsoftlab/ngr-logging"
	"net/smtp"
	"time"
)

// loginAuth struct for login auth
type loginAuth struct {
	username, password string
}

// LoginAuth auth with login and password
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start login auth start
func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next loginAuth
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown fromServer")
		}
	}
	return nil, nil
}

////////////////////////////////////////////////////

// TestLoginAuth - test smtp server login auth
func TestLoginAuth(
	host string, port int,
	tlsOn bool, tlsConfig *tls.Config,
	user, password string,
	connTimeout time.Duration) (connError error, authError error) {

	c, connError := smtp.Dial(fmt.Sprintf("%s:%v", host, port))
	if connError != nil {
		return connError, nil
	}
	defer func() {
		if errClose := c.Close(); errClose != nil {
			logging.Logger.Errorf("smtp conn close error: %s", errClose.Error())
		}
	}()

	if tlsOn {
		if connError = c.StartTLS(tlsConfig); connError != nil {
			logging.Logger.Errorf("start tsl error: %s", connError.Error())
			return connError, nil
		}
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		connTimeout)

	go func(ctx context.Context) {
		defer cancel()

		if authError = c.Auth(LoginAuth(user, password)); authError != nil {
			logging.Logger.Errorf("auth error: %s", authError.Error())
		}
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			logging.Logger.Error("smtp auth conn deadline")
			authError = errors.New("auth conn deadline is over")
		case context.Canceled:
			logging.Logger.Info("smtp auth conn cancel by timeout")
		}
	}

	if authError != nil {
		logging.Logger.Errorf("smtp auth conn test error: %s", authError.Error())
	}

	return connError, authError
}
