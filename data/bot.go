package data

import (
	"eventbot/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type BotUser struct {
	UserId  int64
	Message string
	ChatId  int64
}

type Bot struct {
	bot *tgbotapi.BotAPI
	db  *Data
}

func NewBot(bot *tgbotapi.BotAPI, db *Data) *Bot {
	if bot == nil {
		panic("Bot is nil")
	}

	return &Bot{bot: bot, db: db}
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
			command.NewUserEvent(b.db, update.Message, b.bot).HandleCommand()
		}
	}
}

type UserMessage struct {
	UserId int64
	Text   string
}

func WelcomeMessage() string {
	return "Hello! This is a reminder bot. \nSend \n/newevent to create an event " +
		"\n/myevents to get the list of your events " +
		"\n/edit to update your event, " +
		"\n/delete to delete one of your events, " +
		"\n/deleteall to delete all of your events."
}

func AskForName() string {
	return "Write the name of your event."
}

func AskForDate() string {
	return "Write the date of your event. Use the YYYY-MM-DD format."
}

func AskForTime() string {
	return "Write the time of your event. Use the HH:MM 24h format."
}

func AskHowFrequently() string {
	return "Write if I should remind you about this event weekly, monthly or yearly. If you want to be reminded just once, " +
		"write: Once."
}

//func (r *Bot) SendMessage() {
//
//	msg := tgbotapi.NewMessage(User.ChatId, send-response.SendResponse(User.Message))
//	_, err := r.bot.Send(msg)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
