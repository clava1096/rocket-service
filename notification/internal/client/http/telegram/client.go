package telegram

import (
	"context"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/go-telegram/bot"
)

type client struct {
	bot *bot.Bot
}

func NewClient(bot *bot.Bot) *client {
	return &client{
		bot: bot,
	}
}

func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {

	logger.Info(ctx, "Send telegram message:")

	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	})

	if err != nil {
		return err
	}

	return nil
}
