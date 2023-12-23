package main

import (
	data "eventbot/data"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//db :=
	//db.AddUser(1480532761)
	//
	//db.CreateEvent(1480532761, "Christmas", time.Now(), false, false, true)
	//
	bot := data.NewBot(data.BotStart(), data.NewData(data.MustConnectPostgres()))
	bot.ReadMessage()
}
