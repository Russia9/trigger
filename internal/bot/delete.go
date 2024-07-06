package bot

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"regexp"
)

var DeleteCommand = regexp.MustCompile("-триггер\\s+(.*)")

func (b *Bot) Delete(ctx telebot.Context) error {
	// Check if the message is from a group
	if !ctx.Message().FromGroup() {
		return b.Send(ctx, "🚫 Бот работает только в чатах")
	}

	// Get Chat Admins
	admins, err := b.bot.AdminsOf(ctx.Chat())
	if err != nil {
		return err
	}

	// Check if sender is an admin
	admin := false
	for _, a := range admins {
		if a.User.ID == ctx.Message().Sender.ID {
			admin = true
			break
		}
	}

	if !admin {
		return b.Send(ctx, "🚫 Добавлять триггеры может только админ")
	}

	// Check if the command is valid
	if !DeleteCommand.MatchString(ctx.Text()) {
		return b.Send(ctx, "🚫 Не указан триггер, который нужно удалить")
	}

	// Parse command
	trigger := DeleteCommand.FindStringSubmatch(ctx.Text())[1]

	// Delete trigger
	count, err := b.repo.Delete(context.Background(), trigger, ctx.Chat().ID)
	if err != nil {
		return err
	}

	if count == 0 {
		return b.Send(ctx, "⚠️ Триггер не найден")
	}

	return b.Send(ctx, fmt.Sprintf("✅ Удалено триггеров: %d", count))
}
