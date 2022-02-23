package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
	"os"
	"strconv"
	"time"
)

func main() {
	// Log settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	pretty, err := strconv.ParseBool(os.Getenv("LOG_PRETTY"))
	if err != nil {
		pretty = false
	}
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	switch os.Getenv("LOG_LEVEL") {
	case "DISABLED":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	b, err := telebot.NewBot(telebot.Settings{
		Token:     os.Getenv("TELEGRAM_TOKEN"),
		Poller:    &telebot.LongPoller{Timeout: time.Second * 10},
		ParseMode: "html",
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	b.Handle("+триггер", func(c telebot.Context) error {
		if !c.Message().IsReply() {
			return nil
		}

		var object interface{}
		var entities telebot.Entities

		if c.Message().ReplyTo.Text != "" {
			object = c.Message().ReplyTo.Text
			entities = c.Message().ReplyTo.Entities
		} else if c.Message().ReplyTo.Photo != nil {
			s := c.Message().ReplyTo.Photo
			s.Caption = c.Message().ReplyTo.Caption
			object = s
			entities = c.Message().ReplyTo.CaptionEntities
		} else if c.Message().ReplyTo.Animation != nil {
			s := c.Message().ReplyTo.Animation
			s.Caption = c.Message().ReplyTo.Caption
			object = s
			entities = c.Message().ReplyTo.CaptionEntities
		} else if c.Message().ReplyTo.Video != nil {
			s := c.Message().ReplyTo.Video
			s.Caption = c.Message().ReplyTo.Caption
			object = s
			entities = c.Message().ReplyTo.CaptionEntities
		} else if c.Message().ReplyTo.Voice != nil {
			s := c.Message().ReplyTo.Voice
			s.Caption = c.Message().ReplyTo.Caption
			object = s
			entities = c.Message().ReplyTo.CaptionEntities
		} else if c.Message().ReplyTo.VideoNote != nil {
			object = c.Message().ReplyTo.VideoNote
			entities = c.Message().ReplyTo.CaptionEntities
		} else if c.Message().ReplyTo.Sticker != nil {
			object = c.Message().ReplyTo.Sticker
		} else if c.Message().ReplyTo.Document != nil {
			s := c.Message().ReplyTo.Document
			s.Caption = c.Message().ReplyTo.Caption
			object = s
			entities = c.Message().ReplyTo.CaptionEntities
		}

		return c.Send(object, entities)
	})

	b.Start()
}
