package http

import "context"

type TelegramClient interface {
	SendMessage(ctx context.Context, chatId int64, text string) error
}
