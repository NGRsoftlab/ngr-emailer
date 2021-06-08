// Copyright 2020 NGR Softlab
//
package emailer

import (
	"fmt"
	"testing"
)

/////////////////////////////////////////////////
// Put correct here before testing.
func TestSimpleMail(t *testing.T) {
	s := NewSender(
		"user@mail.com",
		"password",
		"user@mail.com",
		fmt.Sprintf("%v:%v", "test.com", "587"),
	)

	err := s.NewMessage("topic", []string{"user@mail.com"}, "test", []AttachData{{
		FileName: "test.txt",
		FileData: []byte("test"),
	}})
	if err != nil {
		t.Fatal(err)
	}
	err = s.Send()
	if err != nil {
		t.Fatal(err)
	}
}
