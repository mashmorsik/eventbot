package main

import (
	"eventbot/SendResponse"
	data "eventbot/data"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func botStart() *tgbotapi.BotAPI {
	token, _ := os.LookupEnv("EVENTBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

func botSend(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			//userId := update.Message.From.ID
			//EditEvents.AddUser(userId)

			command := update.Message.Text
			fmt.Println(command)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, SendResponse.SendResponse(command)) //SendResponse()
			//switch update.Message.Text {
			//case "open":
			//	msg.ReplyMarkup = SendResponse.SendCalendar()
			//
			//}

			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	db := data.NewData(data.MustConnectPostgres())

	//db.AddUser(6537489206)
	testId := 6537489202

	//bot := botStart()
	//botSend(bot)

	//db.CreateEvent(int64(testId), "Holidays", time.Now(), false, false, true)
	//db.GetEventsList(int64(testId))
	//db.DeleteEvent(2)
	//db.UpdateEvent(3, "New Year", time.Now(), false, false, true)
	//db.DeleteAllEvents(int64(testId))
	//db.GetUsersList()
	db.DeleteUser(int64(testId))
}
