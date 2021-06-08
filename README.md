# emailer
Tiny lib for email (with attachments) sending (smtp)

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

	err := s.NewMessage("topic", []string{"user@mail.com"}, "test", []AttachData{{
		fileName: "test.txt",
		fileData: []byte("test"),
	}})
	if err != nil {
		log.Fatal(err)
	}
	err = s.Send()
	if err != nil {
		log.Fatal(err)
	}
```
