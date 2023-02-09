// Copyright 2020-2023 NGR Softlab
//
package emailer

type MessageParams struct {
	Topic       string
	ContentType string
	Charset     string
	Recipients  []string
	Body        string
	Files       []AttachData
}
