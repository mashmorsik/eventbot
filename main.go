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
	defer func() {
		err := s.Shutdown()
		if err != nil {
			Logger.Sugar.Errorln(err)
		}
	}()

	//db.AddUser(1480532761)
	//

	//db := data.NewData(data.MustConnectPostgres())
	//db.CreateEvent(1480532761, 1480532761, "Christmas", time.Now().Local(), false, false, true)
	//
	////_, err := db.FindRemindEvent()
	////if err != nil {
	////	return
	////}
	//fmt.Println(time.Now())

	bot := NewBot(BotStart(), data.MustConnectPostgres())
	bot.ReadMessage()
	//zone := send_response.MakeDateTimeField("2023-12-29", "11:30")
	//fmt.Println(zone)
}

func RunEventExecute() gocron.Scheduler {
	db := data.MustConnectPostgres()
	dat := data.NewData(db)
	ue := command.NewUserEvent(db, nil, BotStart())

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
				// FIXME: if time_date <= current time and event not fired, execute and set fired
			}
		default:
			cronExp = event.Cron
			cronHandler = func() {
				if ue.HandleCronResponse(event.ChatId, event.Name) != nil {
					Logger.Sugar.Errorln(err)
					return
				}
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

	//// Once
	//o, err := s.NewJob(
	//	gocron.OneTimeJob(
	//		// Get time from DB
	//		gocron.OneTimeJobStartDateTime(time.Now())),
	//	gocron.NewTask(
	//		// Send reminder message
	//		func() {},
	//	),
	//)
	//if err != nil {
	//	Logger.Sugar.Panic(err)
	//}
	//Logger.Sugar.Info("run cron job id ", o.ID())
	//
	//// Daily
	//d, err := s.NewJob(
	//	gocron.DailyJob(
	//		1,
	//		gocron.NewAtTimes(
	//			// Get time from DB
	//			gocron.NewAtTime(10, 30, 0),
	//		),
	//	),
	//	gocron.NewTask(
	//		// Send reminder message
	//		func() {},
	//	),
	//)
	//if err != nil {
	//	Logger.Sugar.Panic(err)
	//}
	//
	//Logger.Sugar.Info("run cron job id ", d.ID())

	s.Start()

	return s
}
