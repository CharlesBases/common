package message

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// EmailMessage 内容
type EmailMessage struct {
	From        string   // 发件人
	To          []string // 收件人
	Cc          []string // 抄送人
	Subject     string   // 标题
	ContentType string   // 内容类型 text/plain text/html
	Content     string   // 内容
	Attach      string   // 附件
}

func NewEmailMessage(from, subject, contentType, content, attach string, to, cc []string) *EmailMessage {
	return &EmailMessage{
		From:        from,
		Subject:     subject,
		ContentType: contentType,
		Content:     content,
		To:          to,
		Cc:          cc,
		Attach:      attach,
	}
}

// EmailClient 发送客户端
type EmailClient struct {
	Host     string // smtp 地址
	Port     int    // 用户名
	Username string // 密码
	Password string // 端口
	Message  *EmailMessage
}

func NewEmailClient(host, username, password string, port int, message *EmailMessage) *EmailClient {
	return &EmailClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Message:  message,
	}
}

// SendMessage 发送邮件
func (ec *EmailClient) SendEmail() (bool, error) {
	d := gomail.NewPlainDialer(ec.Host, ec.Port, ec.Username, ec.Password)
	if 587 == ec.Port {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	m := gomail.NewMessage()
	m.SetHeader("From", ec.Message.From)
	m.SetHeader("To", ec.Message.To...)
	if len(ec.Message.Cc) != 0 {
		m.SetHeader("Cc", ec.Message.Cc...)
	}
	m.SetHeader("Subject", ec.Message.Subject)
	m.SetBody(ec.Message.ContentType, ec.Message.Content)
	if ec.Message.Attach != "" {
		m.Attach(ec.Message.Attach)
	}
	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
	return true, nil
}
