package main

import (
	"eventbot/EditEvents"
	"eventbot/SendResponse"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("6668421178:AAEIq5xPhDY17AFGZpz9BxCmFrXU8eopmQo")
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

			userId := update.Message.From.ID
			EditEvents.AddUser(userId)

			command := update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, SendResponse.SendResponse(command)) //SendResponse()

			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
