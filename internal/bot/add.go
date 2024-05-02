package bot

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/russia9/trigger/pkg/domain"
	"gopkg.in/telebot.v3"
	"regexp"
)

var AddCommand = regexp.MustCompile("\\+триггер\\s+(.*)")

func (b *Bot) Add(ctx telebot.Context) error {
	// Check if the message is from a group
	if !ctx.Message().FromGroup() {
		return b.Send(ctx, "🚫 Бот работает только в чатах")
	}

	// Check if message is a reply
	if !ctx.Message().IsReply() {
		return b.Send(ctx, "ℹ Отправьте команду в ответ на сообщение, которое хотите сохранить")
	}

	// Get Chat member
	member, err := b.bot.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
	if err != nil {
		return err
	}

	// Check if the user is an admin
	if member.Role != telebot.Creator && member.Role != telebot.Administrator {
		return b.Send(ctx, "🚫 Добавлять триггеры может только админ")
	}

	// Check if the command is valid
	if !AddCommand.MatchString(ctx.Text()) {
		return b.Send(ctx, "ℹ Не указано имя триггера")
	}

	// Create a new trigger
	trigger := domain.Trigger{
		ID:      uuid.NewString(),
		Trigger: AddCommand.FindStringSubmatch(ctx.Text())[1],
		Chat:    ctx.Message().Chat.ID,
	}

	// Check the reply message type
	switch {
	case ctx.Message().ReplyTo.Text != "":
		// Text
		trigger.Object = []byte(ctx.Message().ReplyTo.Text)
		trigger.Type = "text"
		trigger.Entities = ctx.Message().ReplyTo.Entities
	case ctx.Message().ReplyTo.Photo != nil:
		// Photo
		s := ctx.Message().ReplyTo.Photo
		s.Caption = ctx.Message().ReplyTo.Caption
		trigger.Object, _ = json.Marshal(s)
		trigger.Type = "photo"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	case ctx.Message().ReplyTo.Animation != nil:
		// Animation
		s := ctx.Message().ReplyTo.Animation
		s.Caption = ctx.Message().ReplyTo.Caption
		trigger.Object, _ = json.Marshal(s)
		trigger.Type = "animation"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	case ctx.Message().ReplyTo.Video != nil:
		// Video
		s := ctx.Message().ReplyTo.Video
		s.Caption = ctx.Message().ReplyTo.Caption
		trigger.Object, _ = json.Marshal(s)
		trigger.Type = "video"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	case ctx.Message().ReplyTo.Voice != nil:
		// Voice
		s := ctx.Message().ReplyTo.Voice
		s.Caption = ctx.Message().ReplyTo.Caption
		trigger.Object, _ = json.Marshal(s)
		trigger.Type = "voice"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	case ctx.Message().ReplyTo.VideoNote != nil:
		// VideoNote
		trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.VideoNote)
		trigger.Type = "videonote"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	case ctx.Message().ReplyTo.Sticker != nil:
		// Sticker
		trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.Sticker)
		trigger.Type = "sticker"
	case ctx.Message().ReplyTo.Document != nil:
		// Document
		s := ctx.Message().ReplyTo.Document
		s.Caption = ctx.Message().ReplyTo.Caption
		trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.Caption)
		trigger.Type = "document"
		trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
	default:
		// Unsupported
		return b.Send(ctx, "треш")
	}

	// Save Trigger to Repository
	err = b.repo.Create(context.Background(), &trigger)
	if err != nil {
		return err
	}

	return b.Send(ctx, "✅ Триггер добавлен")
}
