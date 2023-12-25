package main

import (
	"eventbot/data"
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

	db := data.MustConnectPostgres()
	bot := NewBot(BotStart(), db)
	bot.ReadMessage()
}
