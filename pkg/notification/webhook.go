package notification

import (
	"context"
	"encoding/base64"
	"github.com/liulin110/go-components/utils"
	"net/http"
	"net/url"
)

type WebhookSender struct {
	Url string `json:"url"`
	// optional basic auth
	Authorization string `json:"authorization"`
}

type webhookSenderOption func(WebhookSender *WebhookSender)

func NewWebhookSender(url string, fns ...webhookSenderOption) *WebhookSender {
	sender := &WebhookSender{
		Url: url,
	}
	for _, fn := range fns {
		fn(sender)
	}
	return sender
}

func WithAuthorization(authorization string) webhookSenderOption {
	return func(sender *WebhookSender) {
		sender.Authorization = authorization
	}
}

func (w *WebhookSender) SendNotification(info *NotificationInfo) error {
	response := make(map[string]any)
	webUrl, err := url.Parse(w.Url)
	if err != nil {
		return err
	}
	option := utils.HttpOption{
		Method: http.MethodPost,
		Url:    webUrl,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		RequestBody: info.WebhookContent,
		Response:    response,
	}
	if w.Authorization != "" {
		option.Header["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(w.Authorization))
	}
	return option.Send(context.Background())
}
