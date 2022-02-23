package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/telebot.v3"
	"os"
	"regexp"
	"russia9.dev/trigger/utils"
	"strconv"
	"strings"
	"time"
)

type Trigger struct {
	ID       string `bson:"_id"`
	Trigger  string
	Chat     int64
	Object   interface{}
	Entities telebot.Entities
}

var AddCommand = regexp.MustCompile("\\+триггер\\s+(.*)")
var DelCommand = regexp.MustCompile("-триггер\\s+(.*)")

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

	// Mongo
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal().Str("module", "mongo").Err(err).Send()
	}
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal().Str("module", "mongo").Err(err).Send()
	}
	db := client.Database(utils.GetEnv("MONGO_DB", "bot"))

	b, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: time.Second * 10},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	b.OnError = func(err error, ctx telebot.Context) {
		fmt.Println(err)
		ctx.Send(fmt.Sprintf("Произошла ошибка, пизданите русю и он может быть ее починит\n\n<code>%s</code>", err.Error()), telebot.ModeHTML)
	}

	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		if strings.HasPrefix(ctx.Text(), "+триггер") {
			if !ctx.Message().FromGroup() {
				return ctx.Send("Бот работает только в чатах")
			}

			if !ctx.Message().IsReply() {
				return ctx.Send("Отправьте команду в ответ на сообщение, которое хотите сохранить")
			}

			member, err := b.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
			if err != nil {
				return err
			}

			if member.Role != telebot.Creator && member.Role != telebot.Administrator {
				return ctx.Send("Добавлять триггеры может только админ")
			}

			if !AddCommand.MatchString(ctx.Text()) {
				return ctx.Send("Не указано имя триггера")
			}

			var object interface{}
			var entities telebot.Entities

			if ctx.Message().ReplyTo.Text != "" {
				object = ctx.Message().ReplyTo.Text
				entities = ctx.Message().ReplyTo.Entities
			} else if ctx.Message().ReplyTo.Photo != nil {
				s := ctx.Message().ReplyTo.Photo
				s.Caption = ctx.Message().ReplyTo.Caption
				object = s
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Animation != nil {
				s := ctx.Message().ReplyTo.Animation
				s.Caption = ctx.Message().ReplyTo.Caption
				object = s
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Video != nil {
				s := ctx.Message().ReplyTo.Video
				s.Caption = ctx.Message().ReplyTo.Caption
				object = s
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Voice != nil {
				s := ctx.Message().ReplyTo.Voice
				s.Caption = ctx.Message().ReplyTo.Caption
				object = s
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.VideoNote != nil {
				object = ctx.Message().ReplyTo.VideoNote
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Sticker != nil {
				object = ctx.Message().ReplyTo.Sticker
			} else if ctx.Message().ReplyTo.Document != nil {
				s := ctx.Message().ReplyTo.Document
				s.Caption = ctx.Message().ReplyTo.Caption
				object = s
				entities = ctx.Message().ReplyTo.CaptionEntities
			} else {
				return ctx.Send("треш")
			}

			trigger := Trigger{
				ID:       uuid.New().String(),
				Trigger:  AddCommand.FindAllStringSubmatch(ctx.Text(), -1)[0][1],
				Chat:     ctx.Message().Chat.ID,
				Object:   object,
				Entities: entities,
			}

			_, err = db.Collection("triggers").InsertOne(context.Background(), trigger)
			if err != nil {
				return err
			}

			return ctx.Send("✅ Триггер добавлен")
		}

		if strings.HasPrefix(ctx.Text(), "-триггер") {
			if !ctx.Message().FromGroup() {
				return ctx.Send("Бот работает только в чатах")
			}

			member, err := b.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
			if err != nil {
				return err
			}

			if member.Role != telebot.Creator && member.Role != telebot.Administrator {
				return ctx.Send("Удалять триггеры может только админ")
			}

			if !DelCommand.MatchString(ctx.Text()) {
				return ctx.Send("Не указан триггер, который удаляем")
			}

			count, err := db.Collection("triggers").DeleteMany(context.Background(),
				bson.M{
					"chat":    ctx.Message().Chat.ID,
					"trigger": DelCommand.FindAllStringSubmatch(ctx.Text(), -1)[0][1],
				},
			)
			if err != nil {
				return err
			}

			if count.DeletedCount == 0 {
				return ctx.Send("🚫 Триггер не найден")
			}
			return ctx.Send("✅ Триггер удален")
		}

		if !ctx.Message().FromGroup() {
			return nil
		}

		find, err := db.Collection("triggers").Find(context.Background(), bson.M{"chat": ctx.Message().Chat.ID, "trigger": ctx.Text()})
		if err != nil {
			return err
		}

		for find.Next(context.Background()) {
			var trigger Trigger
			err = find.Decode(&trigger)
			if err != nil {
				return err
			}

			if object, ok := trigger.Object.(string); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Photo); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Animation); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Video); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Voice); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.VideoNote); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Sticker); ok {
				err = ctx.Send(object, trigger.Entities)
			}
			if object, ok := trigger.Object.(*telebot.Document); ok {
				err = ctx.Send(object, trigger.Entities)
			}

			if err != nil {
				return err
			}
		}

		return nil
	})

	b.Start()
}
