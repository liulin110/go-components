package notification

import (
	"context"
	"errors"
	"fmt"
	"github.com/liulin110/go-components/utils"
	"net/http"
	"net/url"
)

type TelegramSender struct {
	BotId  string
	ChatId int64
}

func NewTelegramSender(botId string, chatId int64) *TelegramSender {
	return &TelegramSender{
		BotId:  botId,
		ChatId: chatId,
	}
}

type TelegramBotResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int64  `json:"error_code"`
	Description string `json:"description"`
}

func (t *TelegramSender) SendNotification(info *NotificationInfo) error {
	parse, err := url.Parse(fmt.Sprintf("https://api.telegram.org/%s/sendMessage", "bot"+t.BotId))
	if err != nil {
		return err
	}
	values := url.Values{}
	values.Add("chat_id", fmt.Sprintf("%d", t.ChatId))
	values.Add("text", info.TelegramContent)
	values.Add("parse_mode", "Markdown")
	parse.RawQuery = values.Encode()
	response := TelegramBotResponse{}
	option := utils.HttpOption{
		Method:      http.MethodGet,
		Url:         parse,
		Header:      nil,
		RequestBody: nil,
		Response:    &response,
	}
	err = option.Send(context.Background())
	if err != nil {
		return errors.New(response.Description)
	}
	return nil
}
