package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/russia9/trigger/internal/bot"
	"github.com/russia9/trigger/internal/trigger/repository/mongodb"
	"github.com/russia9/trigger/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Telegram Bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:     os.Getenv("TELEGRAM_TOKEN"),
		Poller:    &telebot.LongPoller{Timeout: time.Second * 10},
		ParseMode: telebot.ModeHTML,
		OnError: func(err error, ctx telebot.Context) {
			fmt.Println(err)
			_ = ctx.Reply(fmt.Sprintf("Произошла ошибка, пизданите русю и он может быть ее починит\n\n<code>%s</code>", err.Error()), telebot.ModeHTML)
		},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Repository
	repo := mongodb.NewTriggerRepository(db)

	// Start bot
	bot.NewBot(b, repo).Start()
}
