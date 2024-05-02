package bot

import (
	"github.com/russia9/trigger/pkg/domain"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot *telebot.Bot

	repo domain.TriggerRepository
}

func NewBot(bot *telebot.Bot, repo domain.TriggerRepository) *Bot {
	b := &Bot{
		bot:  bot,
		repo: repo,
	}

	b.bot.Handle("+триггер", b.Add)
	b.bot.Handle("-триггер", b.Delete)
	b.bot.Handle("+триггеры", b.List)
	b.bot.Handle(telebot.OnText, b.Trigger)

	return b
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) Send(ctx telebot.Context, what interface{}, opts ...interface{}) error {
	if ctx.Message().TopicMessage {
		_, err := b.bot.Reply(ctx.Message().ReplyTo, what, opts...)
		if err != nil {
			return err
		}
	}
	return ctx.Send(what, opts...)
}
