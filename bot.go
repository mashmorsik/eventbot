package main

import (
	"database/sql"
	"eventbot/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	db          *sql.DB
	rerunEvents chan any
}

func NewBot(bot *tgbotapi.BotAPI, db *sql.DB, updateEventsChan chan any) *Bot {
	if bot == nil {
		panic("Bot is nil")
	}

	return &Bot{bot: bot, db: db, rerunEvents: updateEventsChan}
}

func BotStart() *tgbotapi.BotAPI {
	token, _ := os.LookupEnv("EVENTBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func (b *Bot) ReadMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			command.NewUserEvent(b.db, update.Message, b.bot, b.rerunEvents).HandleCommand()
		}
	}
}
