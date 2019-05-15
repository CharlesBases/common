package message

import "testing"

func TestEmail(t *testing.T) {

	dia := GetDialer("userName", "passWord", "host", 465)

	dia.SendEmail(&EmailMessage{
		From:        "form",
		To:          []string{},
		Cc:          []string{},
		Subject:     "subject",
		ContentType: "text/plain",
		Content:     "content",
		Attach:      "attach",
	})
}
