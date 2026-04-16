package telegram

import (
	"context"
	"fmt"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type client struct {
	bot *bot.Bot
}

func NewClient(bot *bot.Bot) *client {

	if bot == nil {
		return nil
	}

	return &client{
		bot: bot,
	}
}

func (c *client) SendMessage(ctx context.Context, chatID int64, text string) error {
	logger.Info(ctx, "Send telegram message")

	if c.bot == nil {
		logger.Error(ctx, "Telegram bot not initialized (network issue)",
			zap.String("chat_id", fmt.Sprintf("%d", chatID)),
			zap.String("message_text", text))
		return fmt.Errorf("telegram bot not initialized")
	}
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
