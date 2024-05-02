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
		return ctx.Reply("🚫 Бот работает только в чатах")
	}

	// Get Chat member
	member, err := b.bot.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
	if err != nil {
		return err
	}

	// Check if the user is an admin
	if member.Role != telebot.Creator && member.Role != telebot.Administrator {
		return ctx.Reply("🚫 Удалять триггеры может только админ")
	}

	// Check if the command is valid
	if !DeleteCommand.MatchString(ctx.Text()) {
		return ctx.Reply("🚫 Не указан триггер, который нужно удалить")
	}

	// Parse command
	trigger := DeleteCommand.FindStringSubmatch(ctx.Text())[1]

	// Delete trigger
	count, err := b.repo.Delete(context.Background(), trigger, ctx.Chat().ID)
	if err != nil {
		return err
	}

	if count == 0 {
		return ctx.Reply("⚠️ Триггер не найден")
	}

	return ctx.Reply(fmt.Sprintf("✅ Удалено триггеров: %d", count))
}
