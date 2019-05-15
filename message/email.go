package message

import (
	"gopkg.in/gomail.v2"
)

// 拨号器
type Dialer struct {
	*gomail.Dialer
}

// 内容
type EmailMessage struct {
	From        string   // 发件人
	To          []string // 收件人
	Cc          []string // 抄送人
	Subject     string   // 标题
	ContentType string   // 内容类型 text/plain text/html
	Content     string   // 内容
	Attach      string   // 附件
}

func GetDialer(userName string, passWord string, host string, port int) *Dialer {
	return &Dialer{gomail.NewDialer(host, port, userName, passWord)}
}

// SendMessage 发送邮件
func (dia *Dialer) SendEmail(em *EmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", em.From)
	m.SetHeader("To", em.To...)
	if len(em.Cc) != 0 {
		m.SetHeader("Cc", em.Cc...)
	}
	m.SetHeader("Subject", em.Subject)
	m.SetBody(em.ContentType, em.Content)
	if em.Attach != "" {
		m.Attach(em.Attach)
	}
	return dia.DialAndSend(m)
}
