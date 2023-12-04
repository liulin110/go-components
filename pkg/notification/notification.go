package notification

type NotificationInfo struct {
	// email
	EmailSubject string
	EmailContent string
	ToEmail      []string
	// tg
	TelegramContent string
	// Slack
	SlackContent string
	// webhook
	WebhookContent any
}

type Notifier interface {
	SendNotification(info *NotificationInfo) error
}

type Informer struct {
	Notifications []Notifier
}

func NewInformer(notifiers ...Notifier) *Informer {
	return &Informer{
		Notifications: notifiers,
	}
}

func (i *Informer) Add(n Notifier) {
	i.Notifications = append(i.Notifications, n)
}

func (i *Informer) NotifyAll(info *NotificationInfo) error {
	for _, n := range i.Notifications {
		err := n.SendNotification(info)
		if err != nil {
			return err
		}
	}
	return nil
}
