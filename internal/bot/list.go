package bot

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
)

func (b *Bot) List(ctx telebot.Context) error {
	// Check if the message is from a group
	if !ctx.Message().FromGroup() {
		return nil
	}

	// Get triggers from Repository
	triggers, err := b.repo.List(context.Background(), ctx.Chat().ID)
	if err != nil {
		return err
	}

	// Create message
	msg := "<b>Список триггеров:</b>\n"
	for i, trigger := range triggers {
		msg += fmt.Sprintf("%d. %s <i>(%s)</i>\n", i+1, trigger.Trigger, trigger.Type)
	}

	return b.Send(ctx, msg)
}
