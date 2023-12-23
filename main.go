package main

import (
	data "eventbot/data"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db := data.NewData(data.MustConnectPostgres())
	db.AddUser(1480532761)

	db.CreateEvent(1480532761, "Christmas", time.Now(), false, false, true)

	bot := data.NewBot(data.BotStart())
	bot.ReadMessage()
}
