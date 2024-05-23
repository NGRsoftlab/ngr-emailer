# emailer
Tiny lib for email (with attachments) sending (smtp, only login auth is supported now)

# import
```import "github.com/NGRsoftlab/ngr-emailer"```

# example
```
s := NewSender(
		"user@mail.com",
		"password",
		"user@mail.com",
		fmt.Sprintf("%v:%v", "test.com", "587"),
	)

		err := s.NewMessage(
		&MessageParams{
			Topic:       "topic",
			ContentType: HtmlContentType,
			Charset:     Utf8Charset,
			Recipients:  []string{"user@test.ru"},
			Body:        body,
			Files: []AttachData{{
				FileName: "test.txt",
				FileData: []byte("test"),
			}},
		})
	if err != nil {
		log.Fatal(err)
	}
	err = s.Send()
	if err != nil {
		log.Fatal(err)
	}
```
