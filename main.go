package main

import (
	"eventbot/Logger"
	"eventbot/command"
	"eventbot/data"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	Logger.InitSugarLogger()
	s := RunEventExecute()

	updateEventsChan := make(chan any)
	go EventUpdate(updateEventsChan, s)

	defer func() {
		err := s.Shutdown()
		if err != nil {
			Logger.Sugar.Errorln(err)
		}
	}()

	bot := NewBot(BotStart(), data.MustConnectPostgres(), updateEventsChan)
	bot.ReadMessage()
}

func EventUpdate(ch <-chan any, scheduler gocron.Scheduler) {
	for e := range ch {
		Logger.Sugar.Infof("received a new event: %s, rerun all events\n", e)
		RerunEventExecute(scheduler)
	}
}

func RerunEventExecute(scheduler gocron.Scheduler) {
	for _, j := range scheduler.Jobs() {
		Logger.Sugar.Debugf("remove event: %s", j.Name())

		err := scheduler.RemoveJob(j.ID())
		if err != nil {
			Logger.Sugar.Errorf("remove job fail %v\n", err)
		}
	}

	RunEventExecute()
}

func RunEventExecute() gocron.Scheduler {
	db := data.MustConnectPostgres()
	dat := data.NewData(db)
	ue := command.NewUserEvent(db, nil, BotStart(), nil)

	events, err := dat.GetAllEvents()
	if err != nil {
		Logger.Sugar.Errorln(err)
		return nil
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		Logger.Sugar.Errorln(err)
		return nil
	}

	for _, event := range events {
		Logger.Sugar.Infoln("set new cron job", slog.Any("event", event))

		var (
			cronExp     string
			cronHandler func()
		)

		switch event.Cron {
		case "once":
			cronExp = "* * * * *"
			cronHandler = func() {
				if ue.HandleOnceReminder() != nil {
					Logger.Sugar.Errorln(err)
					return
				}
			}
		default:
			cronExp = event.Cron
			cronHandler = func() {
				if ue.HandleOnceReminder() != nil {
					Logger.Sugar.Errorln(err)
					return
				}
				//if ue.HandleCronResponse(event.ChatId, event.Name, event.EventId) != nil {
				//	Logger.Sugar.Errorln(err)
				//	return
				//}
			}
		}

		e, err := s.NewJob(
			gocron.CronJob(
				cronExp,
				false,
			),
			gocron.NewTask(cronHandler),
		)
		if err != nil {
			Logger.Sugar.Panic(err)
		}
		Logger.Sugar.Info("run cron job id ", e.ID())
	}

	s.Start()

	return s
}
