package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID      string `bson:"_id"`
	Trigger string
	Chat    int64

	Object   []byte
	Type     string
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
		Token:     os.Getenv("TELEGRAM_TOKEN"),
		Poller:    &telebot.LongPoller{Timeout: time.Second * 10},
		ParseMode: telebot.ModeHTML,
		OnError: func(err error, ctx telebot.Context) {
			fmt.Println(err)
			ctx.Reply(fmt.Sprintf("Произошла ошибка, пизданите русю и он может быть ее починит\n\n<code>%s</code>", err.Error()), telebot.ModeHTML)
		},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	b.Handle(telebot.OnText, func(ctx telebot.Context) error {
		if strings.HasPrefix(ctx.Text(), "+триггер") {
			if !ctx.Message().FromGroup() {
				return ctx.Reply("Бот работает только в чатах")
			}

			if !ctx.Message().IsReply() {
				return ctx.Reply("Отправьте команду в ответ на сообщение, которое хотите сохранить")
			}

			member, err := b.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
			if err != nil {
				return err
			}

			if member.Role != telebot.Creator && member.Role != telebot.Administrator {
				return ctx.Reply("Добавлять триггеры может только админ")
			}

			if !AddCommand.MatchString(ctx.Text()) {
				return ctx.Reply("Не указано имя триггера")
			}

			trigger := Trigger{
				ID:      uuid.New().String(),
				Trigger: AddCommand.FindAllStringSubmatch(ctx.Text(), -1)[0][1],
				Chat:    ctx.Message().Chat.ID,
			}

			if ctx.Message().ReplyTo.Text != "" {
				trigger.Object = []byte(ctx.Message().ReplyTo.Text)
				trigger.Type = "text"
				trigger.Entities = ctx.Message().ReplyTo.Entities
			} else if ctx.Message().ReplyTo.Photo != nil {
				s := ctx.Message().ReplyTo.Photo
				s.Caption = ctx.Message().ReplyTo.Caption
				trigger.Object, _ = json.Marshal(s)
				trigger.Type = "photo"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Animation != nil {
				s := ctx.Message().ReplyTo.Animation
				s.Caption = ctx.Message().ReplyTo.Caption
				trigger.Object, _ = json.Marshal(s)
				trigger.Type = "animation"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Video != nil {
				s := ctx.Message().ReplyTo.Video
				s.Caption = ctx.Message().ReplyTo.Caption
				trigger.Object, _ = json.Marshal(s)
				trigger.Type = "video"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Voice != nil {
				s := ctx.Message().ReplyTo.Voice
				s.Caption = ctx.Message().ReplyTo.Caption
				trigger.Object, _ = json.Marshal(s)
				trigger.Type = "voice"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.VideoNote != nil {
				trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.VideoNote)
				trigger.Type = "videonote"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else if ctx.Message().ReplyTo.Sticker != nil {
				trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.Sticker)
				trigger.Type = "sticker"
			} else if ctx.Message().ReplyTo.Document != nil {
				s := ctx.Message().ReplyTo.Document
				s.Caption = ctx.Message().ReplyTo.Caption
				trigger.Object, _ = json.Marshal(ctx.Message().ReplyTo.Caption)
				trigger.Type = "document"
				trigger.Entities = ctx.Message().ReplyTo.CaptionEntities
			} else {
				return ctx.Reply("треш")
			}

			_, err = db.Collection("triggers").InsertOne(context.Background(), trigger)
			if err != nil {
				return err
			}

			return ctx.Reply("✅ Триггер добавлен")
		}

		if strings.HasPrefix(ctx.Text(), "-триггер") {
			if !ctx.Message().FromGroup() {
				return ctx.Reply("Бот работает только в чатах")
			}

			member, err := b.ChatMemberOf(ctx.Chat(), ctx.Message().Sender)
			if err != nil {
				return err
			}

			if member.Role != telebot.Creator && member.Role != telebot.Administrator {
				return ctx.Reply("Удалять триггеры может только админ")
			}

			if !DelCommand.MatchString(ctx.Text()) {
				return ctx.Reply("Не указан триггер, который удаляем")
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
				return ctx.Reply("🚫 Триггер не найден")
			}
			return ctx.Reply("✅ Триггер удален")
		}

		if !ctx.Message().FromGroup() {
			return nil
		}

		find, err := db.Collection("triggers").Find(context.Background(),
			bson.M{"chat": ctx.Message().Chat.ID,
				"trigger": bson.M{
					"$regex": primitive.Regex{Pattern: "^" + regexp.QuoteMeta(ctx.Text()) + "$", Options: "i"},
				},
			},
		)
		if err != nil {
			return err
		}

		for find.Next(context.Background()) {
			var trigger Trigger
			err = find.Decode(&trigger)
			if err != nil {
				return err
			}

			switch trigger.Type {
			case "text":
				err = ctx.Reply(string(trigger.Object), trigger.Entities, telebot.NoPreview)
			case "photo":
				var photo telebot.Photo
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "animation":
				var photo telebot.Animation
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "video":
				var photo telebot.Video
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "voice":
				var photo telebot.Voice
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "videonote":
				var photo telebot.VideoNote
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "sticker":
				var photo telebot.Sticker
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			case "document":
				var photo telebot.Document
				_ = json.Unmarshal(trigger.Object, &photo)
				err = ctx.Reply(&photo, trigger.Entities, telebot.NoPreview)
			}

			if err != nil {
				return err
			}
		}

		return nil
	})

	b.Handle("+триггеры", func(ctx telebot.Context) error {
		if !ctx.Message().FromGroup() {
			return nil
		}

		find, err := db.Collection("triggers").Find(context.Background(), bson.M{"chat": ctx.Message().Chat.ID})
		if err != nil {
			return err
		}

		var triggers []Trigger
		err = find.All(context.Background(), &triggers)
		if err != nil {
			return err
		}

		msg := "<b>Список триггеров:</b>\n"

		for i, trigger := range triggers {
			msg += fmt.Sprintf("%d. %s <i>(%s)</i>\n", i+1, trigger.Trigger, trigger.Type)
		}

		return ctx.Reply(msg)
	})

	b.Start()
}
