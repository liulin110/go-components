package notification

import "github.com/slack-go/slack"

type SlackSender struct {
	WebhookUrl string
}

func NewSlack(webhookUrl string) *SlackSender {
	return &SlackSender{WebhookUrl: webhookUrl}
}
func (s *SlackSender) SendNotification(info *NotificationInfo) error {
	section := slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, info.SlackContent, false, false), nil, nil)
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			section,
		},
	}
	return slack.PostWebhook(s.WebhookUrl, &slack.WebhookMessage{
		Blocks: &blocks,
	})
}
