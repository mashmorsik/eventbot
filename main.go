package main

import (
	"eventbot/Logger"
	"eventbot/command"
	"eventbot/cron"
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
	Logger.InitSugarLogger()

	db := data.MustConnectPostgres()
	dat := data.NewData(db)
	updateEventsChan := make(chan any)
	sc := cron.StartScheduler()
	s := cron.NewScheduler(sc)

	botApi := BotStart()
	bot := NewBot(BotStart(), db, updateEventsChan, sc)

	ue := command.UserEvent{
		Data:        dat,
		Message:     nil,
		BotAPI:      botApi,
		RerunEvents: updateEventsChan,
		Scheduler:   s,
	}

	_, err := ue.RunScheduler()
	if err != nil {
		Logger.Sugar.Panic(err)
	}

	bot.ReadMessage()
}
