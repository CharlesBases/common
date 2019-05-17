package message

import "testing"

func TestEmail(t *testing.T) {

	dia := GetDialer("host", 465, "userName", "passWord")

	dia.SendEmail(&EmailMessage{
		To:          []string{},
		Cc:          []string{},
		Subject:     "subject",
		ContentType: "text/plain",
		Content:     "content",
		Attach:      "attach",
	})
}
