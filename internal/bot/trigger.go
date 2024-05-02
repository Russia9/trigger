package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v3"
)

func (b *Bot) Trigger(ctx telebot.Context) error {
	// Check if the message is from a group
	if !ctx.Message().FromGroup() {
		return nil
	}

	fmt.Println(ctx.Message().TopicMessage)
	fmt.Println(ctx.Message().TopicCreated)

	// Get triggers from Repository
	triggers, err := b.repo.Get(context.Background(), ctx.Text(), ctx.Message().Chat.ID)
	if err != nil {
		return err
	}

	// Loop through triggers
	for _, trigger := range triggers {
		switch trigger.Type {
		case "text":
			err = ctx.Send(string(trigger.Object), trigger.Entities, telebot.NoPreview)
		case "photo":
			var photo telebot.Photo
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "animation":
			var photo telebot.Animation
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "video":
			var photo telebot.Video
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "voice":
			var photo telebot.Voice
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "videonote":
			var photo telebot.VideoNote
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "sticker":
			var photo telebot.Sticker
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		case "document":
			var photo telebot.Document
			_ = json.Unmarshal(trigger.Object, &photo)
			err = ctx.Send(&photo, trigger.Entities, telebot.NoPreview)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
