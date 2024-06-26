package bot

import (
	"context"
	"encoding/json"
	"gopkg.in/telebot.v3"
)

func (b *Bot) Trigger(ctx telebot.Context) error {
	// Check if the message is from a group
	if !ctx.Message().FromGroup() {
		return nil
	}

	if AddCommand.MatchString(ctx.Text()) {
		return b.Add(ctx)
	} else if DeleteCommand.MatchString(ctx.Text()) {
		return b.Delete(ctx)
	}

	// Get triggers from Repository
	triggers, err := b.repo.Get(context.Background(), ctx.Text(), ctx.Message().Chat.ID)
	if err != nil {
		return err
	}

	// Loop through triggers
	for _, trigger := range triggers {
		switch trigger.Type {
		case "text":
			err = b.Send(ctx, string(trigger.Object), trigger.Entities, telebot.NoPreview)
		case "photo":
			var photo telebot.Photo
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "animation":
			var photo telebot.Animation
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "video":
			var photo telebot.Video
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "voice":
			var photo telebot.Voice
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "videonote":
			var photo telebot.VideoNote
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "sticker":
			var photo telebot.Sticker
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		case "document":
			var photo telebot.Document
			_ = json.Unmarshal(trigger.Object, &photo)
			err = b.Send(ctx, &photo, trigger.Entities, telebot.NoPreview)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
