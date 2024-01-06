package main

import (
	"eventbot/Logger"
	"eventbot/command"
	"eventbot/cron"
	"eventbot/data"
	"github.com/go-co-op/gocron"
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

	botApi := BotStart()
	bot := NewBot(BotStart(), db, updateEventsChan)

	ue := command.UserEvent{
		Data:        dat,
		Message:     nil,
		BotAPI:      botApi,
		RerunEvents: updateEventsChan,
	}

	s, err := cron.RunScheduler(&ue, dat)
	if err != nil {
		Logger.Sugar.Panic(err)
	}

	//defer func() {
	//	err = s.Stop()
	//	if err != nil {
	//		Logger.Sugar.Errorln(err)
	//	}
	//}()

	go func() {
		err = RerunSchedulerObserver(updateEventsChan, &ue, dat, s)
		if err != nil {
			Logger.Sugar.Panic(err)
		}
	}()

	bot.ReadMessage()
}

func RerunSchedulerObserver(ch <-chan any, userEvent *command.UserEvent, data *data.Data, scheduler *gocron.Scheduler,
) error {
	for e := range ch {
		Logger.Sugar.Infof("received a new event: %s, rerun all events", e)

		err := cron.RerunScheduler(userEvent, data, scheduler)
		if err != nil {
			return err
		}
	}

	Logger.Sugar.Warnf("return from RerunSchedulerObserver")
	return nil
}
