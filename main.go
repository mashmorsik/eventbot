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
	//db.AddUser(1480532761)
	//
	//db.CreateEvent(1480532761, "Christmas", time.Now(), false, false, true)
	//
	//db := data.NewData(data.MustConnectPostgres())
	//
	//_, err := db.FindRemindEvent()
	//if err != nil {
	//	return
	//}
	bd := data.MustConnectPostgres()
	bot := NewBot(BotStart(), bd)
	bot.ReadMessage()

}
