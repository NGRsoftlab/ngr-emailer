// Copyright 2020-2023 NGR Softlab
//
package emailer

import (
	"errors"
	"net/smtp"
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
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
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
