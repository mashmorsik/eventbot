package main

import (
	"eventbot/SendResponse"
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

func main() {
	token, _ := os.LookupEnv("EVENTBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

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

			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
