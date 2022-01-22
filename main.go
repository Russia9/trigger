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

		if c.Message().ReplyTo.Text != "" {
			return c.Send(c.Message().ReplyTo.Text, c.Message().ReplyTo.Entities)
		} else if c.Message().ReplyTo.Photo != nil {
			s := c.Message().ReplyTo.Photo
			s.Caption = c.Message().ReplyTo.Caption
			return c.Send(c.Message().ReplyTo.Photo, c.Message().ReplyTo.CaptionEntities)
		} else if c.Message().ReplyTo.Animation != nil {
			return c.Send(c.Message().ReplyTo.Animation)
		} else if c.Message().ReplyTo.Video != nil {
			return c.Send(c.Message().ReplyTo.Video)
		} else if c.Message().ReplyTo.Voice != nil {
			return c.Send(c.Message().ReplyTo.Voice)
		} else if c.Message().ReplyTo.VideoNote != nil {
			return c.Send(c.Message().ReplyTo.VideoNote)
		} else if c.Message().ReplyTo.Sticker != nil {
			return c.Send(c.Message().ReplyTo.Sticker)
		} else if c.Message().ReplyTo.Document != nil {
			return c.Send(c.Message().ReplyTo.Document)
		}

		//return c.Send(send)
		return nil
	})

	b.Start()
}
