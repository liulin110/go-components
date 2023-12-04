package notification

import "gopkg.in/gomail.v2"

type EmailSender struct {
	user     string
	alias    string
	password string
	host     string
	port     int
}

func NewEmailSender(user, alias, password, host string, port int) *EmailSender {
	return &EmailSender{
		user:     user,
		alias:    alias,
		password: password,
		host:     host,
		port:     port,
	}
}

func (e *EmailSender) SendNotification(n *NotificationInfo) error {
	msg := gomail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From":    {msg.FormatAddress(e.user, e.alias)},
		"To":      n.ToEmail,
		"Subject": {n.EmailSubject},
	})
	msg.SetBody("text/html", n.EmailContent)
	d := gomail.Dialer{
		Host:     e.host,
		Port:     e.port,
		Username: e.user,
		Password: e.password,
		SSL:      false,
	}

	err := d.DialAndSend(msg)
	return err
}
