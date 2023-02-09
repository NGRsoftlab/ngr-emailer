// Copyright 2020-2023 NGR Softlab
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
	err = s.Send()
	if err != nil {
		t.Fatal(err)
	}
}
